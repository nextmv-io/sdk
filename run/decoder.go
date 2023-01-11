package run

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/nextmv-io/sdk/run/decode"
)

// GenericDecoder returns a new generic decoder.
func GenericDecoder[Input any](
	decoder decode.Decoder,
) Decoder[Input] {
	dec := genericDecoder[Input]{decoder}
	return dec.Decoder
}

type genericDecoder[Input any] struct {
	decoder decode.Decoder
}

// Decoder is a function that decodes the input from the reader. It uses the
// given decoder to decode the input. If the input is gzipped, it will be
// decoded using the gzip.reader.
func (g *genericDecoder[Input]) Decoder(
	_ context.Context, reader any) (input Input, err error,
) {
	closer, ok := reader.(io.Closer)
	if ok {
		defer func() {
			tempErr := closer.Close()
			// the first error is the most important
			if err == nil {
				err = tempErr
			}
		}()
	}

	ioReader, ok := reader.(io.Reader)
	if !ok {
		err = errors.New(
			"Decoder is not compatible with configured IOProducer",
		)
		return input, err
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

// QueryParamDecoder is a Decoder that returns option from query params.
func QueryParamDecoder[Option any](
	_ context.Context, reader any,
) (option Option, err error) {
	urlValues, ok := reader.(url.Values)
	if !ok {
		return option, errors.New(
			"QueryParamDecoder is not compatible with configured IOProducer",
		)
	}

	if len(urlValues) == 0 {
		return option, nil
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&option, urlValues)
	return option, err
}
