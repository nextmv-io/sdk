// package main holds the implementation of a simple runner example.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/schema"
)

func main() {
	err := run.HTTP(algorithm,
		// listen on port 9003
		run.SetAddr[input, option, schema.Output](":9003"),
		// set the maximum number of parallel requests to 2
		run.SetMaxParallel[input, option, schema.Output](2),
		// override the default logger
		run.SetLogger[input, option, schema.Output](
			log.New(os.Stdout, "[demo] - ", log.LstdFlags),
		),
	).Run(context.Background())
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
