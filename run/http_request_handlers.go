package run

import (
	"bytes"
	"context"
	"errors"
	"net/http"
)

// SyncHTTPRequestHandler allows the input and option to be sent as body and
// query parameters.
func SyncHTTPRequestHandler(
	w http.ResponseWriter, req *http.Request,
) (Callback, IOProducer[HTTPRunnerConfig], error) {
	return nil, func(ctx context.Context, config HTTPRunnerConfig) IOData {
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

// AsyncHTTPRequestHandler creates a new HTTPRequestHandler.
func AsyncHTTPRequestHandler(
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
) (Callback, IOProducer[HTTPRunnerConfig], error) {
	callbackURL := a.callbackURL
	if a.requestOverride {
		headerCallbackURL := req.Header.Get("callback_url")
		if headerCallbackURL != "" {
			callbackURL = headerCallbackURL
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
	callbackFunc := func(requestID string) (err error) {
		// Create a new request
		callbackReq, err := http.NewRequestWithContext(
			context.Background(), "POST", callbackURL, buf,
		)
		if err != nil {
			return err
		}
		// Set the GUID header
		callbackReq.Header.Set("request_id", requestID)
		// Send the request
		resp, err := a.httpClient.Do(callbackReq)
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
	return callbackFunc, func(
		ctx context.Context, config HTTPRunnerConfig,
	) IOData {
		return NewIOData(
			req.Body,
			req.URL.Query(),
			buf,
		)
	}, nil
}
