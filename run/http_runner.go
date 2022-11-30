package run

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
)

// HTTPRunner is a Runner that uses HTTP as its IO.
type HTTPRunner[Input, Option, Solution any] interface {
	Runner[Input, Option, Solution]
	// SetRun allows to fully configure and control the run process, e.g.
	// setting up and listening to a http server. Since a handler is needed, we
	// also provide GetHttpHandler to get the handler.
	SetRun(func(context.Context) error)
	GetHTTPHandler() http.Handler
	// GetOneOffRunner returns a Runner that is used internally. This method is
	// useful when users want to implement their own http.Handler.
	GetOneOffRunner() Runner[Input, Option, Solution]
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

// DefaultHTTPOneOffRunner is the default http one-off runner.
func DefaultHTTPOneOffRunner[Input, Option, Solution any](
	handler Algorithm[Input, Option, Solution],
) Runner[Input, Option, Solution] {
	return NewOneOffRunner(
		nil,
		JSONDecoder[Input],
		HeaderDecoder[Option],
		handler,
		JSONEncoder[Solution],
	)
}

// NewHTTPRunner creates a new HTTPRunner.
func NewHTTPRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...HTTPRunnerOption[Input, Option, Solution],
) HTTPRunner[Input, Option, Solution] {
	runner := &httpRunner[Input, Option, Solution]{
		oneOffRunner: DefaultHTTPOneOffRunner(algorithm),
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
	oneOffRunner Runner[Input, Option, Solution]
	run          func(context.Context) error
	httpServer   *http.Server
}

func (h *httpRunner[Input, Option, Solution]) setHTTPAddr(addr string) {
	h.httpServer.Addr = addr
}

func (h *httpRunner[Input, Option, Solution]) Run(
	context context.Context,
) error {
	return h.run(context)
}

// SetIOHandler sets the ioHandler of a runner using f.
func (h *httpRunner[Input, Option, Solution]) SetIOProducer(
	ioHandler IOProducer,
) {
	h.oneOffRunner.SetIOProducer(ioHandler)
}

func (h *httpRunner[Input, Option, Solution]) SetInputDecoder(
	decoder InputDecoder[Input],
) {
	h.oneOffRunner.SetInputDecoder(decoder)
}

func (h *httpRunner[Input, Option, Solution]) SetOptionDecoder(
	decoder OptionDecoder[Option],
) {
	h.oneOffRunner.SetOptionDecoder(decoder)
}

// SetHandler sets the handler of a runner using f.
func (h *httpRunner[Input, Option, Solution]) SetAlgorithm(
	algorithm Algorithm[Input, Option, Solution],
) {
	h.oneOffRunner.SetAlgorithm(algorithm)
}

// SetEncoder sets the encoder of a runner using f.
func (h *httpRunner[Input, Option, Solution]) SetEncoder(
	encoder Encoder[Solution],
) {
	h.oneOffRunner.SetEncoder(encoder)
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

// GetOneOffRunner returns the one-off runner of the runner.
func (h *httpRunner[Input, Option, Solution]) GetOneOffRunner() Runner[
	Input, Option, Solution,
] {
	return h.oneOffRunner
}

// ServeHTTP implements the http.Handler interface.
func (h *httpRunner[Input, Option, Solution]) ServeHTTP(
	w http.ResponseWriter, req *http.Request,
) {
	runner := h.oneOffRunner
	var reader io.Reader = req.Body
	var writer io.Writer = w
	runner.SetIOProducer(
		func(ctx context.Context, config any) IOData {
			return NewIOData(
				reader,
				req.Header,
				writer,
			)
		},
	)
	err := runner.Run(context.Background())
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
