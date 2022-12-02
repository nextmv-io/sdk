package run

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/itzg/go-flagsfiller"
	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// GenericDecoder is a Decoder that decodes a json into a struct.
func GenericDecoder[Input any, Decoder decode.Decoder](
	_ context.Context, reader any) (input Input, err error,
) {
	ioReader, ok := reader.(io.Reader)
	if !ok {
		return input, errors.New(
			"JsonDecoder is not compatible with configured IOProducer",
		)
	}

	// Convert to buffered reader and read magic bytes
	bufferedReader := bufio.NewReader(ioReader)
	testBytes, err := bufferedReader.Peek(2)

	// Test for gzip magic bytes and use corresponding reader, if given
	if err == nil && testBytes[0] == 31 && testBytes[1] == 139 {
		var gzipReader *gzip.Reader
		if gzipReader, err = gzip.NewReader(bufferedReader); err != nil {
			return input, err
		}
		ioReader = gzipReader
	} else {
		// Default case: assume text input
		ioReader = bufferedReader
	}

	decoder := *new(Decoder)
	err = decoder.Decode(ioReader, &input)
	return input, err
}

// NoopOptionsDecoder is a Decoder that returns the option as is.
func NoopOptionsDecoder[Input any](
	_ context.Context, _ any, input Input,
) (Input, error) {
	return input, nil
}

// FlagParser parses flags and env vars.
func FlagParser[Option, RunnerCfg any]() (
	runnerConfig RunnerCfg, option Option, err error,
) {
	// create a FlagSetFiller
	filler := flagsfiller.New(
		flagsfiller.WithEnv(""),
		flagsfiller.WithFieldRenamer(
			func(name string) string {
				repl := strings.ReplaceAll(name, "-", ".")
				return strings.ToLower(repl)
			},
		),
	)
	err = filler.Fill(flag.CommandLine, &option)
	if err != nil {
		return runnerConfig, option, err
	}

	err = filler.Fill(flag.CommandLine, &runnerConfig)
	if err != nil {
		return runnerConfig, option, err
	}

	flag.Parse()

	return runnerConfig, option, nil
}

// CliIOProducer is a test IOProducer.
func CliIOProducer(_ context.Context, config any) IOData {
	cfg, ok := config.(CliRunnerConfig)
	if !ok {
		log.Fatal("DefaultIOProducer is not compatible with the runner")
	}
	reader := os.Stdin
	if cfg.Runner.Input.Path != "" {
		r, err := os.Open(cfg.Runner.Input.Path)
		if err != nil {
			log.Fatal(err)
		}
		reader = r
	}
	var writer io.Writer = os.Stdout
	if cfg.Runner.Output.Path != "" {
		w, err := os.Create(cfg.Runner.Output.Path)
		if err != nil {
			log.Fatal(err)
		}
		writer = w
	}
	return NewIOData(
		reader,
		nil,
		writer,
	)
}

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
