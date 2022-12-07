package run

import (
	"bytes"
	"context"
	"errors"
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
	// SetMaxParallel sets the maximum number of parallel requests.
	SetMaxParallel(int)
	// HandlerToIOProducer configures the IOProducer based on the http request.
	HandlerToIOProducer(
		func(w http.ResponseWriter, req *http.Request,
		) (IOProducer, error))
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

// SetMaxParallel sets the maximum number of parallel requests.
func SetMaxParallel[Input, Option, Solution any](maxParallel int) func(
	HTTPRunner[Input, Option, Solution],
) {
	return func(r HTTPRunner[Input, Option, Solution]) {
		r.SetMaxParallel(maxParallel)
	}
}

// HandlerToIOProducer configures the IOProducer based on the http request.
func HandlerToIOProducer[Input, Option, Solution any](
	f func(w http.ResponseWriter, req *http.Request) (IOProducer, error),
) func(
	HTTPRunner[Input, Option, Solution],
) {
	return func(r HTTPRunner[Input, Option, Solution]) {
		r.HandlerToIOProducer(f)
	}
}

// NewHTTPRunner creates a new HTTPRunner.
func NewHTTPRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...HTTPRunnerOption[Input, Option, Solution],
) HTTPRunner[Input, Option, Solution] {
	runner := &httpRunner[Input, Option, Solution]{
		genericRunner: &genericRunner[Input, Option, Solution]{
			InputDecoder:  NewGenericDecoder[Input](decode.JSON()),
			OptionDecoder: NewGenericDecoder[Option](decode.JSON()),
			Algorithm:     algorithm,
			Encoder:       NewGenericEncoder[Solution, Option](encode.JSON()),
		},
	}

	runnerConfig, decodedOption, err := FlagParser[
		Option, HTTPRunnerConfig,
	]()
	if err != nil {
		panic(err)
	}
	runner.genericRunner.runnerConfig = runnerConfig
	runner.genericRunner.decodedOption = decodedOption

	runner.maxParallel = make(chan struct{}, runnerConfig.Runner.HTTP.MaxParallel)

	// default http server
	runner.httpServer = &http.Server{
		Addr:     runnerConfig.Runner.HTTP.Address,
		ErrorLog: log.New(os.Stderr, "HTTPRunner", log.LstdFlags),
		Handler:  runner,
	}

	// default handler to IOProducer
	runner.handlerToIOProducer = MultiPartHandlerToIOProducer

	for _, option := range options {
		option(runner)
	}

	return runner
}

type httpRunner[Input, Option, Solution any] struct {
	*genericRunner[Input, Option, Solution]
	httpServer          *http.Server
	maxParallel         chan struct{}
	handlerToIOProducer func(
		w http.ResponseWriter, req *http.Request,
	) (IOProducer, error)
}

// MultiPartHandlerToIOProducer allows the input and option to be sent as a
// multipart request.
func MultiPartHandlerToIOProducer(
	w http.ResponseWriter, req *http.Request,
) (IOProducer, error) {
	read_form, err := req.MultipartReader()
	if err != nil {
		return nil, err
	}
	var inputReader, optionReader io.Reader
	for {
		part, err_part := read_form.NextPart()
		if err_part == io.EOF {
			break
		}
		switch part.FormName() {
		case "input":
			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, part); err != nil {
				return nil, err
			}
			inputReader = buf
		case "option":
			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, part); err != nil {
				return nil, err
			}
			optionReader = buf
		}
	}
	if inputReader == nil {
		return nil, errors.New("input not found")
	}
	var writer io.Writer = w
	return func(ctx context.Context, config any) IOData {
		return NewIOData(
			inputReader,
			optionReader,
			writer,
		)
	}, nil
}

func (h *httpRunner[Input, Option, Solution]) SetHTTPAddr(addr string) {
	h.httpServer.Addr = addr
}

func (h *httpRunner[Input, Option, Solution]) SetLogger(l *log.Logger) {
	h.httpServer.ErrorLog = l
}

func (h *httpRunner[Input, Option, Solution]) SetMaxParallel(maxParallel int) {
	h.maxParallel = make(chan struct{}, maxParallel)
}

func (h *httpRunner[Input, Option, Solution]) HandlerToIOProducer(
	f func(w http.ResponseWriter, req *http.Request) (IOProducer, error),
) {
	h.handlerToIOProducer = f
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
	for {
		select {
		case h.maxParallel <- struct{}{}:
			// We have a free slot, so we can start a new run.
			defer func() { <-h.maxParallel }()
			break
		default:
			// No free slot, so we immediately return an error.
			http.Error(w, "max number of parallel requests exceeded",
				http.StatusTooManyRequests)
		}
		break
	}

	// configure how to turn the request and response into an IOProducer.
	producer, err := h.handlerToIOProducer(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.SetIOProducer(producer)

	err = h.genericRunner.Run(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
			MaxParallel int    `default:"1" usage:"The maximum number of parallel requests"`
		}
		Output struct {
			Solutions string `default:"all" usage:"Return all or last solution"`
			Quiet     bool   `default:"false" usage:"Do not return statistics"`
		}
	}
}

// Quiet returns the quiet flag.
func (c HTTPRunnerConfig) Quiet() bool {
	return c.Runner.Output.Quiet
}

// Solutions returns the configured solutions.
func (c HTTPRunnerConfig) Solutions() (Solutions, error) {
	return ParseSolutions(c.Runner.Output.Solutions)
}
