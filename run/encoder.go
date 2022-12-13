package run

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/run/encode"
)

type version struct {
	Sdk string `json:"sdk"`
}
type meta[Options any] struct {
	Version version `json:"version"`
	Options Options `json:"options"`
	Store   string  `json:"store"`
}

// GenericEncoder returns a new Encoder that encodes the solution.
func GenericEncoder[Solution, Options any](
	encoder encode.Encoder,
) Encoder[Solution, Options] {
	enc := genericEncoder[Solution, Options]{encoder}
	return enc.Encode
}

type genericEncoder[Solution, Options any] struct {
	encoder encode.Encoder
}

// GenericEncoder is an Encoder that encodes a struct.
func (g *genericEncoder[Solution, Options]) Encode(
	_ context.Context,
	solutions <-chan Solution,
	writer any,
	runnerCfg any,
	options Options,
) error {
	ioWriter, ok := writer.(io.Writer)
	if !ok {
		return errors.New("JsonEncoder is not compatible with configured IOProducer")
	}

	if outputPather, ok := runnerCfg.(OutputPather); ok {
		if strings.HasSuffix(outputPather.OutputPath(), ".gz") {
			ioWriter = gzip.NewWriter(ioWriter)
		}
	}

	if quieter, ok := runnerCfg.(Quieter); ok && !quieter.Quiet() {
		meta := meta[Options]{
			Version: version{
				Sdk: sdk.VERSION,
			},
			Options: options,
		}
		// Write version
		buf := new(bytes.Buffer)
		if err := g.encoder.Encode(buf, meta); err != nil {
			return err
		}
		_, err := ioWriter.Write(bytes.TrimRight(buf.Bytes(), "\"}\n"))
		if err != nil {
			return err
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

	if err := jsonEncodeChan(g.encoder, ioWriter, solutions); err != nil {
		return err
	}

	if quieter, ok := runnerCfg.(Quieter); ok && !quieter.Quiet() {
		if _, err := ioWriter.Write([]byte{'}'}); err != nil {
			return err
		}
	}

	return nil
}

func jsonEncodeChan[Encoder encode.Encoder](
	encoder Encoder, w io.Writer, vc any,
) (err error) {
	cval := reflect.ValueOf(vc)
	if _, err = w.Write([]byte{'['}); err != nil {
		return
	}
	v, ok := cval.Recv()
	if !ok {
		_, err = w.Write([]byte{']'})
		return err
	}
	// create buffer & encoder only if we have a value
	buf := new(bytes.Buffer)
	goto Encode
Loop:
	if v, ok = cval.Recv(); !ok {
		_, err = w.Write([]byte{']'})
		return err
	}
	if _, err = w.Write([]byte{','}); err != nil {
		return err
	}
Encode:
	err = encoder.Encode(buf, v.Interface())
	if err == nil {
		_, err = w.Write(bytes.TrimRight(buf.Bytes(), "\n"))
		if err != nil {
			return err
		}
		buf.Reset()
	}
	if err != nil {
		return err
	}
	goto Loop
}
