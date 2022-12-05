package run

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"io"

	"github.com/nextmv-io/sdk/run/decode"
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
