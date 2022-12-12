package run

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

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
	// SetHTTPRequestHandler sets the function that handles the http request.
	SetHTTPRequestHandler(HTTPRequestHandler)
}

// Callback is a function that is called after the request is processed.
type Callback func() error

// HTTPRequestHandler is a function that handles an http request.
type HTTPRequestHandler func(
	w http.ResponseWriter, req *http.Request,
) (Callback, IOProducer, error)

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

// SetHTTPRequestHandler sets the function that handles the http request.
func SetHTTPRequestHandler[Input, Option, Solution any](
	f HTTPRequestHandler) func(HTTPRunner[Input, Option, Solution],
) {
	return func(r HTTPRunner[Input, Option, Solution]) {
		r.SetHTTPRequestHandler(f)
	}
}

// NewHTTPRunner creates a new HTTPRunner.
func NewHTTPRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...HTTPRunnerOption[Input, Option, Solution],
) HTTPRunner[Input, Option, Solution] {
	runner := &httpRunner[Input, Option, Solution]{
		genericRunner: genericRunner[Input, Option, Solution]{
			InputDecoder:  NewGenericDecoder[Input](decode.JSON()),
			OptionDecoder: QueryParamDecoder[Option],
			Algorithm:     algorithm,
			Encoder:       NewGenericEncoder[Solution, Option](encode.JSON()),
		},
	}

	runnerConfig, decodedOption, err := FlagParser[
		Option, HTTPRunnerConfig,
	]()
	if err != nil {
		log.Fatal(err)
	}
	runner.genericRunner.runnerConfig = runnerConfig
	runner.genericRunner.decodedOption = decodedOption

	runner.maxParallel = make(chan struct{}, runnerConfig.Runner.HTTP.MaxParallel)

	// default http server
	runner.httpServer = &http.Server{
		Addr:     runnerConfig.Runner.HTTP.Address,
		ErrorLog: log.New(os.Stderr, "[Nextmv HTTPRunner] ", log.LstdFlags),
		Handler:  runner,
	}

	// default handler to IOProducer
	runner.httpRequestHandler = NewHTTPRequestHandler

	for _, option := range options {
		option(runner)
	}

	return runner
}

type httpRunner[Input, Option, Solution any] struct {
	genericRunner[Input, Option, Solution]
	httpServer         *http.Server
	maxParallel        chan struct{}
	httpRequestHandler HTTPRequestHandler
}

// NewHTTPRequestHandler allows the input and option to be sent as body and
// query parameters.
func NewHTTPRequestHandler(
	w http.ResponseWriter, req *http.Request,
) (Callback, IOProducer, error) {
	return nil, func(ctx context.Context, config any) IOData {
		return NewIOData(
			req.Body,
			req.URL.Query(),
			w,
		)
	}, nil
}

// AsyncHTTPRequestHandlerOption configures an AsyncHTTPRequestHandler.
type AsyncHTTPRequestHandlerOption func(*asyncHTTPHandler)

// CallbackURL sets a default callback url.
func CallbackURL(url string) AsyncHTTPRequestHandlerOption {
	return func(h *asyncHTTPHandler) { h.callbackURL = url }
}

// RequestOverride sets whether to allow the callback url to be overridden by
// the request header (callback_url).
func RequestOverride(allow bool) AsyncHTTPRequestHandlerOption {
	return func(h *asyncHTTPHandler) { h.requestOverride = allow }
}

// NewAsyncHTTPRequestHandler creates a new HTTPRequestHandler.
func NewAsyncHTTPRequestHandler(
	options ...AsyncHTTPRequestHandlerOption,
) HTTPRequestHandler {
	handler := &asyncHTTPHandler{
		httpClient:      http.DefaultClient,
		requestOverride: true,
	}
	for _, option := range options {
		option(handler)
	}
	return handler.Handler
}

type asyncHTTPHandler struct {
	httpClient      *http.Client
	callbackURL     string
	requestOverride bool
}

func (a asyncHTTPHandler) Handler(
	_ http.ResponseWriter, req *http.Request,
) (Callback, IOProducer, error) {
	callbackURL := a.callbackURL
	if a.requestOverride {
		headerCBURL := req.Header.Get("callback_url")
		if headerCBURL != "" {
			callbackURL = headerCBURL
		}
		if callbackURL == "" {
			return nil, nil, errors.New(
				"callback_url not configured and not found in header",
			)
		}
	} else if callbackURL == "" {
		return nil, nil, errors.New("callback_url not configured")
	}

	buf := new(bytes.Buffer)
	callbackFunc := func() (err error) {
		resp, err := a.httpClient.Post(callbackURL, "application/json", buf)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				err = cerr
			}
		}()
		return err
	}
	return callbackFunc, func(ctx context.Context, config any) IOData {
		return NewIOData(
			req.Body,
			req.URL.Query(),
			buf,
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

func (h *httpRunner[Input, Option, Solution]) SetHTTPRequestHandler(
	f HTTPRequestHandler,
) {
	h.httpRequestHandler = f
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
	select {
	case h.maxParallel <- struct{}{}:
	default:
		// No free slot, so we immediately return an error.
		http.Error(w, "max number of parallel requests exceeded",
			http.StatusTooManyRequests)
		return
	}

	// control mechanism to let the request by run async or not.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() { <-h.maxParallel }()
		// configure how to turn the request and response into an IOProducer.
		callbackFunc, producer, err := h.httpRequestHandler(w, req)
		async := callbackFunc != nil
		if async {
			wg.Done()
		} else {
			defer wg.Done()
		}
		if err != nil {
			handleError(async, err, w)
			return
		}
		// get a copy of the genericRunner set the IOProducer and run it.
		genericRunner := h.genericRunner
		genericRunner.SetIOProducer(producer)
		err = genericRunner.Run(context.Background())
		if err != nil {
			handleError(async, err, w)
			return
		}

		// if the request is async, call the callbackFunc.
		if async {
			err = callbackFunc()
			if err != nil {
				handleError(async, err, w)
				return
			}
		}
	}()
	wg.Wait()
}

func handleError(async bool, err error, w http.ResponseWriter) {
	log.Println(err)
	if !async {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HTTPRunnerConfig is the configuration of the HTTPRunner.
type HTTPRunnerConfig struct {
	Runner struct {
		Log  *log.Logger
		HTTP struct {
			Address     string `default:":9000" usage:"The host address"`
			Certificate string `usage:"The certificate file path"`
			Key         string `usage:"The key file path"`
			MaxParallel int    `default:"1" usage:"The max number of requests"`
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
