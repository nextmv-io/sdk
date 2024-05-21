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
	file, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = run.HTTP(algorithm,
		// listen on port 9002
		run.SetAddr[input, option, schema.Output](":9002"),
		// set the maximum number of parallel requests to 2
		run.SetMaxParallel[input, option, schema.Output](2),
		// override the default logger
		run.SetLogger[input, option, schema.Output](
			log.New(file, "[demo] - ", log.Lshortfile),
		),
	).Run(context.Background())
	if err != nil {
		log.Println(err)
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
