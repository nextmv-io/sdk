package run

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/run/encode"
	"github.com/nextmv-io/sdk/run/schema"
)

// LegacyEncoder returns a new Encoder that encodes the solution using the
// given encoder.
func LegacyEncoder[Solution, Options any](
	encoder encode.Encoder,
) Encoder[Solution, Options] {
	enc := legacyEncoder[Solution, Options]{encoder}
	return &enc
}

type legacyEncoder[Solution, Options any] struct {
	encoder encode.Encoder
}

// Encode encodes the solution using the given encoder. If a given output path
// ends in .gz, it will be gzipped after encoding. The writer needs to be an
// io.Writer.
func (g *legacyEncoder[Solution, Options]) Encode(
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

	if outputPather, ok := runnerCfg.(OutputPather); ok {
		if strings.HasSuffix(outputPather.OutputPath(), ".gz") {
			ioWriter = gzip.NewWriter(ioWriter)
		}
	}

	if limiter, ok := runnerCfg.(SolutionLimiter); ok {
		solutionFlag, retErr := limiter.Solutions()
		if retErr != nil {
			return retErr
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
	m := schema.Output{
		Version: schema.Version{
			Sdk: sdk.VERSION,
		},
		Options: options,
	}
	for solution := range solutions {
		m.Solutions = append(m.Solutions, solution)
	}
	return g.encoder.Encode(ioWriter, m)
}

func (g *legacyEncoder[Solution, Options]) ContentType() string {
	contentTyper, ok := g.encoder.(ContentTyper)
	if !ok {
		return "text/plain"
	}
	return contentTyper.ContentType()
}
