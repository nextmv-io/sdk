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
	// SetHTTPAddr sets the address the http server listens on.
	SetHTTPAddr(string)
	// SetLogger sets the logger of the http server.
	SetLogger(*log.Logger)
}

// HTTPRunnerOption configures a HTTPRunner.
type HTTPRunnerOption[Input, Option, Solution any] func(
	HTTPRunner[Input, Option, Solution],
)

// SetAddr sets the address the http server listens on.
func SetAddr[Input, Option, Solution any](addr string) func(
	HTTPRunner[Input, Option, Solution],
) {
	return func(r HTTPRunner[Input, Option, Solution]) { r.SetHTTPAddr(addr) }
}

// SetLogger sets the logger of the http server.
func SetLogger[Input, Option, Solution any](l *log.Logger) func(
	HTTPRunner[Input, Option, Solution],
) {
	return func(r HTTPRunner[Input, Option, Solution]) { r.SetLogger(l) }
}

// NewHTTPRunner creates a new HTTPRunner.
func NewHTTPRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...HTTPRunnerOption[Input, Option, Solution],
) HTTPRunner[Input, Option, Solution] {
	runner := &httpRunner[Input, Option, Solution]{
		genericRunner: &genericRunner[Input, Option, Solution]{
			InputDecoder:  NewGenericDecoder[Input](decode.JSON()),
			OptionDecoder: HeaderDecoder[Option],
			Algorithm:     algorithm,
			Encoder:       NewGenericEncoder[Solution, Option](encode.JSON()),
		},
	}

	runnerConfig, decodedOption, err := FlagParser[
		Option, HTTPRunnerConfig,
	]()
	runner.genericRunner.runnerConfig = runnerConfig
	runner.genericRunner.decodedOption = decodedOption
	if err != nil {
		panic(err)
	}

	// default http server
	runner.httpServer = &http.Server{
		Addr:     runnerConfig.Runner.HTTP.Address,
		ErrorLog: log.New(os.Stderr, "HTTPRunner", log.LstdFlags),
		Handler:  runner,
	}

	for _, option := range options {
		option(runner)
	}

	return runner
}

type httpRunner[Input, Option, Solution any] struct {
	*genericRunner[Input, Option, Solution]
	httpServer *http.Server
}

func (h *httpRunner[Input, Option, Solution]) SetHTTPAddr(addr string) {
	h.httpServer.Addr = addr
}

func (h *httpRunner[Input, Option, Solution]) SetLogger(l *log.Logger) {
	h.httpServer.ErrorLog = l
}

func (h *httpRunner[Input, Option, Solution]) Run(
	context context.Context,
) error {
	httpRunnerConfig := h.genericRunner.runnerConfig.(HTTPRunnerConfig)
	if httpRunnerConfig.Runner.HTTP.Certificate != "" ||
		httpRunnerConfig.Runner.HTTP.Key != "" {
		return h.httpServer.ListenAndServeTLS(
			httpRunnerConfig.Runner.HTTP.Certificate,
			httpRunnerConfig.Runner.HTTP.Key,
		)
	}
	return h.httpServer.ListenAndServe()
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
	err := h.genericRunner.Run(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HeaderDecoder is a Decoder that decodes a header into a struct.
func HeaderDecoder[Option any](
	context context.Context, header any, option Option,
) (Option, error) {
	// TODO: transform headers to output
	return option, nil
}

// HTTPRunnerConfig is the configuration of the HTTPRunner.
type HTTPRunnerConfig struct {
	Runner struct {
		Log  *log.Logger
		HTTP struct {
			Address     string `default:":9000" usage:"The host address"`
			Certificate string `usage:"The certificate file path"`
			Key         string `usage:"The key file path"`
		}
	}
}
