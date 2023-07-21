package run

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/nextmv-io/sdk/run/encode"
)

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
// ends in .gz, it will be gzipped after encoding. The writer needs to be an
// io.Writer.
func (g *genericEncoder[Solution, Options]) Encode(
	_ context.Context,
	solutions <-chan Solution,
	writer any,
	runnerCfg any,
	_ Options,
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
			lastIsSet := false
			for solution := range solutions {
				last = solution
				lastIsSet = true
			}
			if !lastIsSet {
				return nil
			}
			tempSolutions := make(chan Solution, 1)
			tempSolutions <- last
			close(tempSolutions)
			solutions = tempSolutions
		}
	}

	for solution := range solutions {
		err := g.encoder.Encode(ioWriter, solution)
		if err != nil {
			return err
		}
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
