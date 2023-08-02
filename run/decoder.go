package run

import (
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

// Decoder is a function that decodes the input from the reader, which needs to
// be an io.Reader. It uses the given decoder to decode the input. If the input
// is gzipped, it will be decoded using the gzip.Reader.
func (g *genericDecoder[Input]) Decoder(
	_ context.Context, reader any) (input Input, err error,
) {
	ioReader, ok := reader.(io.Reader)
	if !ok {
		err = errors.New(
			"decoder is not compatible with configured IOProducer",
		)
		return input, err
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
