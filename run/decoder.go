package run

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"io"

	"github.com/nextmv-io/sdk/run/decode"
)

// NewGenericDecoder returns a new generic decoder.
func NewGenericDecoder[Input any](
	decoder decode.Decoder,
) Decoder[Input] {
	dec := genericDecoder[Input]{decoder}
	return dec.Decoder
}

type genericDecoder[Input any] struct {
	decoder decode.Decoder
}

// GenericDecoder is an InputDecoder that decodes a json into a struct.
func (g *genericDecoder[Input]) Decoder(
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

	err = g.decoder.Decode(ioReader, &input)
	return input, err
}

// NoopOptionsDecoder is a Decoder that returns the option as is.
func NoopOptionsDecoder[Option any](
	_ context.Context, _ any,
) (option Option, err error) {
	return option, nil
}
