package run

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/nextmv-io/sdk/run/decode"
	"github.com/nextmv-io/sdk/run/encode"
)

// Callback is a function that is called after the request is processed. It is
// used to send the result asynchronously to some other service. The argument is
// the request id.
type Callback func(string) error

// HTTPRequestHandler is a function that handles an http request.
type HTTPRequestHandler func(
	w http.ResponseWriter, req *http.Request,
) (Callback, IOProducer[HTTPRunnerConfig], error)

// HTTPRunnerOption configures a HTTPRunner.
type HTTPRunnerOption[Input, Option, Solution any] func(
	*httpRunner[Input, Option, Solution],
)

// SetAddr sets the address the http server listens on.
func SetAddr[Input, Option, Solution any](addr string) func(
	*httpRunner[Input, Option, Solution],
) {
	return func(r *httpRunner[Input, Option, Solution]) { r.setHTTPAddr(addr) }
}

// SetLogger sets the logger of the http server.
func SetLogger[Input, Option, Solution any](l *log.Logger) func(
	*httpRunner[Input, Option, Solution],
) {
	return func(r *httpRunner[Input, Option, Solution]) { r.setLogger(l) }
}

// SetMaxParallel sets the maximum number of parallel requests.
func SetMaxParallel[Input, Option, Solution any](maxParallel int) func(
	*httpRunner[Input, Option, Solution],
) {
	return func(r *httpRunner[Input, Option, Solution]) {
		r.setMaxParallel(maxParallel)
	}
}

// SetHTTPRequestHandler sets the function that handles the http request.
func SetHTTPRequestHandler[Input, Option, Solution any](
	f HTTPRequestHandler) func(*httpRunner[Input, Option, Solution],
) {
	return func(r *httpRunner[Input, Option, Solution]) {
		r.setHTTPRequestHandler(f)
	}
}

// HTTPRunner creates a new HTTPRunner.
func HTTPRunner[Input, Option, Solution any](
	algorithm Algorithm[Input, Option, Solution],
	options ...HTTPRunnerOption[Input, Option, Solution],
) Runner[HTTPRunnerConfig, Input, Option, Solution] {
	runner := &httpRunner[Input, Option, Solution]{
		// the IOProducer will be dynamically set by the http request handler.
		Runner: GenericRunner[HTTPRunnerConfig](
			nil,
			GenericDecoder[Input](decode.JSON()),
			QueryParamDecoder[Option],
			algorithm,
			GenericEncoder[Solution, Option](encode.JSON()),
		),
	}

	runnerConfig := runner.Runner.RunnerConfig()
	runner.maxParallel = make(chan struct{}, runnerConfig.Runner.HTTP.MaxParallel)

	// default http server
	runner.httpServer = &http.Server{
		Addr:     runnerConfig.Runner.HTTP.Address,
		ErrorLog: log.New(os.Stderr, "[Nextmv HTTPRunner] ", log.LstdFlags),
		Handler:  runner,
	}

	// default handler to IOProducer
	runner.httpRequestHandler = SyncHTTPRequestHandler

	for _, option := range options {
		option(runner)
	}

	return runner
}

type httpRunner[Input, Option, Solution any] struct {
	Runner[HTTPRunnerConfig, Input, Option, Solution]
	httpServer         *http.Server
	maxParallel        chan struct{}
	httpRequestHandler HTTPRequestHandler
}

func (h *httpRunner[Input, Option, Solution]) setHTTPAddr(addr string) {
	h.httpServer.Addr = addr
}

func (h *httpRunner[Input, Option, Solution]) setLogger(l *log.Logger) {
	h.httpServer.ErrorLog = l
}

func (h *httpRunner[Input, Option, Solution]) setMaxParallel(maxParallel int) {
	h.maxParallel = make(chan struct{}, maxParallel)
}

func (h *httpRunner[Input, Option, Solution]) setHTTPRequestHandler(
	f HTTPRequestHandler,
) {
	h.httpRequestHandler = f
}

func (h *httpRunner[Input, Option, Solution]) Run(
	context context.Context,
) error {
	httpRunnerConfig := h.Runner.RunnerConfig()
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
		// generate a new requestID
		requestID := uuid.New().String()
		if async {
			// write the guid to the response.
			_, err = w.Write([]byte(requestID))
			if err != nil {
				handleError(async, err, w)
				return
			}
			wg.Done()
		} else {
			defer wg.Done()
		}
		if err != nil {
			handleError(async, err, w)
			return
		}
		// get a copy of the genericRunner set the IOProducer and run it.
		genericRunner := h.Runner
		genericRunner.SetIOProducer(producer)
		err = genericRunner.Run(context.Background())
		if err != nil {
			handleError(async, err, w)
			return
		}

		// if the request is async, call the callbackFunc.
		if async {
			err = callbackFunc(requestID)
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
