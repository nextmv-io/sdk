package run

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// HTTPRunner is a Runner that uses HTTP as its IO.
type HTTPRunner[Input, Option, Solution any] interface {
	Runner[Input, Option, Solution]
	// SetRun allows to fully configure and control the run process, e.g.
	// setting up and listening to a http server. Since a handler is needed, we
	// also provide GetHttpHandler to get the handler.
	SetRun(func(context.Context) error)
	GetHTTPHandler() http.Handler
	// GetGenericRunner returns a Runner that is used internally. This method is
	// useful when users want to implement their own http.Handler.
	GetGenericRunner() Runner[Input, Option, Solution]
}

// HTTPRunnerOption configures a HttpRunner.
type HTTPRunnerOption[Input, Option, Solution any] func(
	*httpRunner[Input, Option, Solution],
)

// SetAddr sets the address the http server listens on.
func SetAddr[Input, Option, Solution any](addr string) func(
	*httpRunner[Input, Option, Solution],
) {
	return func(r *httpRunner[Input, Option, Solution]) { r.setHTTPAddr(addr) }
}

// NewHTTPRunner creates a new HTTPRunner.
func NewHTTPRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...HTTPRunnerOption[Input, Option, Solution],
) HTTPRunner[Input, Option, Solution] {
	runner := &httpRunner[Input, Option, Solution]{
		Runner: NewGenericRunner(
			nil,
			NewGenericDecoder[Input](decode.JSON()),
			HeaderDecoder[Option],
			algorithm,
			GenericEncoder[Solution, Option, encode.JSONEncoder],
		),
	}

	// default http server
	runner.httpServer = &http.Server{
		Addr:     ":9000",
		ErrorLog: log.New(os.Stderr, "", log.LstdFlags),
		Handler:  runner,
	}
	// default value for run
	runner.run = func(context.Context) error {
		return runner.httpServer.ListenAndServe()
	}

	for _, option := range options {
		option(runner)
	}

	return runner
}

type httpRunner[Input, Option, Solution any] struct {
	Runner[Input, Option, Solution]
	run        func(context.Context) error
	httpServer *http.Server
}

func (h *httpRunner[Input, Option, Solution]) setHTTPAddr(addr string) {
	h.httpServer.Addr = addr
}

func (h *httpRunner[Input, Option, Solution]) Run(
	context context.Context,
) error {
	return h.run(context)
}

// SetRun sets the run function of a runner using f.
func (h *httpRunner[Input, Option, Solution]) SetRun(
	f func(context.Context) error,
) {
	h.run = f
}

// GetHTTPHandler returns the http.Handler of the runner.
func (h *httpRunner[Input, Option, Solution]) GetHTTPHandler() http.Handler {
	return h
}

// GetGenericRunner returns the one-off runner of the runner.
func (h *httpRunner[Input, Option, Solution]) GetGenericRunner() Runner[
	Input, Option, Solution,
] {
	return h
}

// ServeHTTP implements the http.Handler interface.
func (h *httpRunner[Input, Option, Solution]) ServeHTTP(
	w http.ResponseWriter, req *http.Request,
) {
	var reader io.Reader = req.Body
	var writer io.Writer = w
	h.SetIOProducer(
		func(ctx context.Context, config any) IOData {
			return NewIOData(
				reader,
				req.Header,
				writer,
			)
		},
	)
	err := h.Run(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HeaderDecoder is a Decoder that decodes a header into a struct.
func HeaderDecoder[Input any](
	context context.Context, header any, input Input,
) (Input, error) {
	// TODO: transform headers to output
	return input, nil
}
