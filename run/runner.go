package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// Runner defines the interface of the runner.
type Runner interface {
	// Run invokes a solver by invoking the associated handler.
	Run()

	// SetDecoder sets decoder to be used to decode input.
	SetDecoder(decoder decode.Decoder)
	// SetEncoder sets encoder to be used to encode output.
	SetEncoder(encoder encode.Encoder)
	// SetHandler sets the handler to be used by the run invocation.
	SetHandler(any)
}

var newFunc func() Runner
