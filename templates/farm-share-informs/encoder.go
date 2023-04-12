package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/nextmv-io/sdk"
	schema "github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/encode"
)

type output struct {
	Unassigned []route.Stop           `json:"unassigned"`
	Vehicles   []schema.VehicleOutput `json:"vehicles"`
	Objective  schema.JSONObjective   `json:"objective"`
}

// Statistics of the search.
type StatisticsIn struct {
	Time Time `json:"time"`
	// Value of the store. Nil when using a Satisfier.
	Value *int `json:"value,omitempty"`
}

type Time struct {
	Start   time.Time `json:"start"`
	Elapsed Duration  `json:"elapsed"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type StatisticsOut struct {
	Schema string `json:"schema"`
	Run    struct {
		Time float64 `json:"time"`
	} `json:"run"`
	Result result `json:"result"`
}

type result struct {
	Value   float64 `json:"value"`
	Elapsed float64 `json:"elapsed"`
	Custom  custom  `json:"custom"`
}

type version struct {
	Sdk string `json:"sdk"`
}
type meta[Options, Solution any] struct {
	Version    version       `json:"version"`
	Options    Options       `json:"options"`
	Solutions  []Solution    `json:"solutions"`
	Statistics StatisticsOut `json:"statistics"`
}

type custom struct {
	Routing      routing `json:"routing"`
	UsedVehicles int     `json:"used_vehicles"`
}

type routing struct {
	Stops stops `json:"stops"`
}

type stops struct {
	Unassigned int `json:"unassigned"`
	Assigned   int `json:"assigned"`
}

// GenericEncoder returns a new Encoder that encodes the solution using the
// given encoder.
func GenericEncoder[Solution, Options any](
	encoder encode.Encoder,
) run.Encoder[Solution, Options] {
	enc := genericEncoder[Solution, Options]{encoder}
	return &enc
}

type genericEncoder[Solution, Options any] struct {
	encoder encode.Encoder
}

// Encode encodes the solution using the given encoder. If a given output path
// ends in .gz, it will be gzipped after encoding. The writer needs to be an
// io.Writer.
func (g *genericEncoder[Solution, Options]) Encode(
	_ context.Context,
	solutions <-chan Solution,
	writer any,
	runnerCfg any,
	options Options,
) (err error) {
	closer, ok := writer.(io.Closer)
	if ok {
		defer func() {
			tempErr := closer.Close()
			// the first error is the most important
			if err == nil {
				err = tempErr
			}
		}()
	}

	ioWriter, ok := writer.(io.Writer)
	if !ok {
		err = errors.New("encoder is not compatible with configured IOProducer")
		return err
	}

	if outputPather, ok := runnerCfg.(run.OutputPather); ok {
		if strings.HasSuffix(outputPather.OutputPath(), ".gz") {
			ioWriter = gzip.NewWriter(ioWriter)
		}
	}

	if limiter, ok := runnerCfg.(run.SolutionLimiter); ok {
		solutionFlag, retErr := limiter.Solutions()
		if retErr != nil {
			return retErr
		}

		if solutionFlag == run.Last {
			var last Solution
			for solution := range solutions {
				last = solution
			}
			tempSolutions := make(chan Solution, 1)
			tempSolutions <- last
			close(tempSolutions)
			solutions = tempSolutions
		}
	}

	if quieter, ok := runnerCfg.(run.Quieter); ok && !quieter.Quiet() {
		m := meta[Options, Solution]{}
		m.Version = version{
			Sdk: sdk.VERSION,
		}
		m.Options = options
		for solution := range solutions {
			m.Solutions = append(m.Solutions, solution)
		}
		if len(m.Solutions) > 0 {
			s := output{}
			b, err := json.Marshal(m.Solutions[0])
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, &s)
			if err != nil {
				return err
			}

			assigned := 0
			usedVehicles := 0
			for _, v := range s.Vehicles {
				if len(v.Route) > 2 {
					assigned += len(v.Route) - 2
					usedVehicles++
				}
			}

			unassigned := 0
			if len(s.Unassigned) > 0 {
				unassigned = len(s.Unassigned)
			}
			if s.Vehicles != nil {
				m.Statistics = StatisticsOut{
					Schema: "v0",
					Result: result{
						Value:   float64(*&s.Objective.Value),
						Elapsed: 0,
						Custom: custom{
							Routing: routing{
								Stops: stops{
									Unassigned: unassigned,
									Assigned:   assigned,
								},
							},
							UsedVehicles: usedVehicles,
						},
					},
				}
			}
		}
		if err = g.encoder.Encode(ioWriter, m); err != nil {
			return err
		}

		return nil
	}

	m := []Solution{}
	for solution := range solutions {
		m = append(m, solution)
	}
	if err = g.encoder.Encode(ioWriter, m); err != nil {
		return err
	}

	return nil
}

func (g *genericEncoder[Solution, Options]) ContentType() string {
	contentTyper, ok := g.encoder.(run.ContentTyper)
	if !ok {
		return "text/plain"
	}
	return contentTyper.ContentType()
}
