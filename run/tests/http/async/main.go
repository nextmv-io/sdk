// package main holds the implementation of a simple runner example.
package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/schema"
)

func main() {
	// start a callback server listening on port 8080
	go func() {
		handler := http.HandlerFunc(callback)
		http.Handle("/callback", handler)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := run.HTTP(algorithm,
		// listen on port 9001
		run.SetAddr[input, option, schema.Output](":9001"),
		// set the maximum number of parallel requests to 2
		run.SetMaxParallel[input, option, schema.Output](2),
		// override the default logger
		run.SetLogger[input, option, schema.Output](
			log.New(os.Stdout, "[demo] - ", log.LstdFlags),
		),
		// send solutions to the callback URL instead of returning them directly
		run.SetHTTPRequestHandler[input, option, schema.Output](
			run.AsyncHTTPRequestHandler(
				// configure the callback URL to send solutions to the callback
				// server started above
				run.CallbackURL("http://localhost:8080/callback"),
				// a caller of the HTTP endpoint is not allowed to override the
				// default callback URL defined above
				run.RequestOverride(false),
			),
		),
	).Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func callback(_ http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	_, err := b.ReadFrom(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("callback.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(b.String())
	if err != nil {
		log.Fatal(err)
	}
	// you can access the request ID via the request header
	// r.Header.Get("request_id")
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
