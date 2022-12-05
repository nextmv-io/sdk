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

// GenericEncoder is an Encoder that encodes a struct.
func GenericEncoder[Solution, Options any, Encoder encode.Encoder](
	_ context.Context,
	solutions <-chan Solution,
	writer any,
	runnerCfg any,
	options Options,
) error {
	encoder := *new(Encoder)
	ioWriter, ok := writer.(io.Writer)
	if !ok {
		return errors.New("JsonEncoder is not compatible with configured IOProducer")
	}
	runnerConfig, ok := runnerCfg.(CliRunnerConfig)
	if !ok {
		return errors.New("JsonEncoder is not compatible with configured IOProducer")
	}
	if strings.HasSuffix(runnerConfig.Runner.Output.Path, ".gz") {
		ioWriter = gzip.NewWriter(ioWriter)
	}

	if !runnerConfig.Runner.Output.Quiet {
		meta := meta[Options]{
			Version: version{
				Sdk: sdk.VERSION,
			},
			Options: options,
		}
		// Write version
		buf := new(bytes.Buffer)
		if err := encoder.Encode(buf, meta); err != nil {
			return err
		}
		_, err := ioWriter.Write(bytes.TrimRight(buf.Bytes(), "\"}\n"))
		if err != nil {
			return err
		}
	}

	if runnerConfig.Runner.Output.Solutions == Last {
		var last Solution
		for solution := range solutions {
			last = solution
		}
		tempSolutions := make(chan Solution, 1)
		tempSolutions <- last
		close(tempSolutions)
		solutions = tempSolutions
	}

	if err := jsonEncodeChan(encoder, ioWriter, solutions); err != nil {
		return err
	}

	if !runnerConfig.Runner.Output.Quiet {
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
