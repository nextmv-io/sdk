package run

import (
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

type Runner interface {
	Run()

	SetDecoder(decode.Decoder)
	SetEncoder(encode.Encoder)
	SetHandler(any)
}

var newFunc func() Runner
