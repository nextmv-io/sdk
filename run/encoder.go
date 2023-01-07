package run

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/run/encode"
)

type version struct {
	Sdk string `json:"sdk"`
}
type meta[Options, Solution any] struct {
	Version   version    `json:"version"`
	Options   Options    `json:"options"`
	Solutions []Solution `json:"solutions"`
}

// GenericEncoder returns a new Encoder that encodes the solution using the
// given encoder.
func GenericEncoder[Solution, Options any](
	encoder encode.Encoder,
) Encoder[Solution, Options] {
	enc := genericEncoder[Solution, Options]{encoder}
	return &enc
}

type genericEncoder[Solution, Options any] struct {
	encoder encode.Encoder
}

// Encode encodes the solution using the given encoder. If a given output path
// ends in .gz, it will be gzipped after encoding.
func (g *genericEncoder[Solution, Options]) Encode(
	_ context.Context,
	solutions <-chan Solution,
	writer any,
	runnerCfg any,
	options Options,
) error {
	ioWriter, ok := writer.(io.Writer)
	if !ok {
		return errors.New("Encoder is not compatible with configured IOProducer")
	}

	if outputPather, ok := runnerCfg.(OutputPather); ok {
		if strings.HasSuffix(outputPather.OutputPath(), ".gz") {
			ioWriter = gzip.NewWriter(ioWriter)
		}
	}

	if limiter, ok := runnerCfg.(SolutionLimiter); ok {
		solutionFlag, err := limiter.Solutions()
		if err != nil {
			return err
		}

		if solutionFlag == Last {
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
	if quieter, ok := runnerCfg.(Quieter); ok && !quieter.Quiet() {
		m := meta[Options, Solution]{
			Version: version{
				Sdk: sdk.VERSION,
			},
			Options: options,
		}
		for solution := range solutions {
			m.Solutions = append(m.Solutions, solution)
		}
		if err := g.encoder.Encode(ioWriter, m); err != nil {
			return err
		}

		return nil
	}

	m := []Solution{}
	for solution := range solutions {
		m = append(m, solution)
	}
	if err := g.encoder.Encode(ioWriter, m); err != nil {
		return err
	}

	return nil
}

func (g *genericEncoder[Solution, Options]) ContentType() string {
	contentTyper, ok := g.encoder.(ContentTyper)
	if !ok {
		return "text/plain"
	}
	return contentTyper.ContentType()
}
