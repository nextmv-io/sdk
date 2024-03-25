// package main holds the implementation of a simple runner example.
package main

import (
	"context"
	"log"
	"time"

	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/schema"
)

func main() {
	err := run.CLI(algorithm).Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

type input struct {
	Message string `json:"message" usage:"Message to print."`
}

type option struct {
	Duration time.Duration `json:"duration" default:"1s" usage:"Sleep duration."`
}

type output struct {
	Message string `json:"message"`
}

func algorithm(_ context.Context, input input, opts option) (schema.Output, error) {
	// sleep for the specified duration, 1s by default as defined via go tags
	time.Sleep(opts.Duration)
	return schema.NewOutput(opts, output{Message: input.Message + " World!"}), nil
}
