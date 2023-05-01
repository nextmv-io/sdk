package here

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nextmv-io/sdk/route"
)

// Client represents a HERE maps client. See official documentation for HERE
// topics, getting started.
type Client interface {
	// DistanceMatrix retrieves a HERE distance matrix. It uses the async HERE API
	// if there are more than 500 points given.
	DistanceMatrix(
		ctx context.Context,
		points []route.Point,
		opts ...MatrixOption,
	) (route.ByIndex, error)
	// DurationMatrix retrieves a HERE duration matrix. It uses the async HERE API
	// if there are more than 500 points given.
	DurationMatrix(
		ctx context.Context,
		points []route.Point,
		opts ...MatrixOption,
	) (route.ByIndex, error)

	// DistanceDurationMatrices retrieves a HERE distance and duration matrix. It
	// uses the async HERE API if there are more than 500 points given.
	DistanceDurationMatrices(
		ctx context.Context,
		points []route.Point,
		opts ...MatrixOption,
	) (distances, durations route.ByIndex, err error)
}

type client struct {
	// scheme and host of the HERE API - currently configurable for testing
	// although this may eventually be useful for failover to alternative
	// regions
	schemeHost              string
	maxAsyncPollingInterval time.Duration
	minAsyncPollingInterval time.Duration
	retries                 int
	denyRedirectedRequests  []string
	APIKey                  string
	httpClient              *http.Client
	// maxSyncPoints controls the maximum number of points that should
	// be requested from the sync endpoint - there is a maximum set by
	// HERE (500 points) and is configurable below that number for testing
	maxSyncPoints int
}

const (
	defaultHereAPIHost = "matrix.router.hereapi.com"
)

// NewClient returns a new OSRM Client.
func NewClient(apiKey string, opts ...ClientOption) Client {
	c := &client{
		schemeHost:              fmt.Sprintf("https://%s", defaultHereAPIHost),
		maxAsyncPollingInterval: time.Second * 5,
		minAsyncPollingInterval: time.Millisecond * 200,
		retries:                 10,
		APIKey:                  apiKey,
		httpClient:              http.DefaultClient,
		maxSyncPoints:           500,
		denyRedirectedRequests:  []string{},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.denyRedirectedRequests = append(c.denyRedirectedRequests, defaultHereAPIHost)
	c.httpClient.CheckRedirect = func(r *http.Request, via []*http.Request) error {
		for _, host := range c.denyRedirectedRequests {
			if strings.HasSuffix(r.URL.Hostname(), host) {
				return http.ErrUseLastResponse
			}
		}

		return nil
	}

	return c
}

func cleanPoints(points []route.Point) []route.Point {
	cleanPoints := make([]route.Point, len(points))
	copy(cleanPoints, points)
	for i, p := range cleanPoints {
		if len(p) == 2 && p[0] == 0 && p[1] == 0 {
			cleanPoints[i] = route.Point{}
		}
	}
	return cleanPoints
}

// DistanceMatrix retrieves a HERE distance matrix. It uses the async HERE API
// if there are more than 500 points given.
func (c *client) DistanceMatrix(
	ctx context.Context,
	points []route.Point,
	opts ...MatrixOption,
) (route.ByIndex, error) {
	points = cleanPoints(points)
	if len(points) > c.maxSyncPoints {
		distances, _, err := c.fetchMatricesAsync(ctx, points, true, false, opts)
		return distances, err
	}
	distances, _, err := c.fetchMatricesSync(ctx, points, true, false, opts)
	return distances, err
}

// DurationMatrix retrieves a HERE duration matrix. It uses the async HERE API
// if there are more than 500 points given.
func (c *client) DurationMatrix(
	ctx context.Context,
	points []route.Point,
	opts ...MatrixOption,
) (route.ByIndex, error) {
	points = cleanPoints(points)
	if len(points) > c.maxSyncPoints {
		_, durations, err := c.fetchMatricesAsync(
			ctx, points, false, true, opts)
		return durations, err
	}
	_, durations, err := c.fetchMatricesSync(ctx, points, false, true, opts)
	return durations, err
}

// DistanceDurationMatrices retrieves a HERE distance and duration matrix. It
// uses the async HERE API if there are more than 500 points given.
func (c *client) DistanceDurationMatrices(
	ctx context.Context,
	points []route.Point,
	opts ...MatrixOption,
) (distances, durations route.ByIndex, err error) {
	points = cleanPoints(points)
	if len(points) > c.maxSyncPoints {
		return c.fetchMatricesAsync(ctx, points, true, true, opts)
	}
	return c.fetchMatricesSync(ctx, points, true, true, opts)
}

// fetchMatricesSync makes a call to the sync HERE API endpoint.
func (c *client) fetchMatricesSync(
	ctx context.Context,
	points []route.Point,
	includeDistance,
	includeDuration bool,
	opts []MatrixOption,
) (distances, durations route.ByIndex, err error) {
	resp, err := c.calculate(
		ctx, points, false, includeDistance, includeDuration, opts...)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if tempErr := resp.Body.Close(); tempErr != nil {
			err = tempErr
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, badStatusError(resp)
	}

	var hereResponse matrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&hereResponse); err != nil {
		return nil, nil, fmt.Errorf("decoding response: %v", err)
	}

	if includeDistance {
		distances = route.Matrix(reshape(
			hereResponse.Matrix.Distances,
			points,
		))
	}
	if includeDuration {
		durations = route.Matrix(reshape(
			hereResponse.Matrix.TravelTimes,
			points,
		))
	}

	return distances, durations, nil
}

// fetchMatricesAsync makes a call to the async HERE API endpoint.
func (c *client) fetchMatricesAsync( //nolint:gocyclo
	ctx context.Context,
	points []route.Point,
	includeDistance,
	includeDuration bool,
	opts []MatrixOption,
) (distances, durations route.ByIndex, err error) {
	statusURL, err := c.startAsyncCalculation(
		ctx, points, includeDistance, includeDuration, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("starting async calculation: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, statusURL, nil)
	if err != nil {
		return nil, nil, err
	}

	pollInterval := c.minAsyncPollingInterval
	retries := 0
	// returns true if there are more retries allowed, and increments retry
	// counters
	shouldRetry := func() bool {
		if retries >= c.retries {
			return false
		}
		retries++
		pollInterval *= 2
		if pollInterval > c.maxAsyncPollingInterval {
			pollInterval = c.maxAsyncPollingInterval
		}
		return true
	}

	var resultURL string
	for {
		var ready, retry bool
		resultURL, ready, retry, err = c.poll(req)
		if err != nil {
			if !retry {
				return nil, nil, fmt.Errorf(
					"an error occurred while polling status: %v", err)
			}
			if !shouldRetry() {
				return nil, nil, fmt.Errorf(
					"maximum number of retries (%d) exceeded", c.retries)
			}
		}
		if ready {
			break
		}

		select {
		case <-ctx.Done():
			return nil, nil, context.Canceled
		case <-time.After(pollInterval):
		}
	}

	retries = 0
	pollInterval = c.minAsyncPollingInterval
	var resp *http.Response
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, resultURL, nil)
	if err != nil {
		return nil, nil, err
	}
	for {
		resp, err = c.httpClient.Do(req)
		if err != nil && !shouldRetry() || errors.Is(err, context.Canceled) {
			return nil, nil, fmt.Errorf("getting result: %v", err)
		}
		defer func() {
			if tempErr := resp.Body.Close(); tempErr != nil {
				err = tempErr
			}
		}()
		if err != nil {
			return nil, nil, fmt.Errorf("closing response body: %w", err)
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		if resp.StatusCode > http.StatusBadRequest &&
			resp.StatusCode < http.StatusInternalServerError {
			return nil, nil, validationError{
				message: fmt.Sprintf("aborting request due to received status code: %d",
					resp.StatusCode),
			}
		}
		if resp.StatusCode >= http.StatusInternalServerError && !shouldRetry() {
			return nil, nil, fmt.Errorf("received status code: %d",
				resp.StatusCode)
		}

		select {
		case <-ctx.Done():
			return nil, nil, context.Canceled
		case <-time.After(pollInterval):
		}
	}

	var hereResponse matrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&hereResponse); err != nil {
		return nil, nil, fmt.Errorf("decoding result: %v", err)
	}

	if includeDistance {
		distances = route.Matrix(reshape(
			hereResponse.Matrix.Distances,
			points,
		))
	}
	if includeDuration {
		durations = route.Matrix(reshape(
			hereResponse.Matrix.TravelTimes,
			points,
		))
	}

	return distances, durations, nil
}

func badStatusError(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(
			"bad status code: %d, error reading response body: %v",
			resp.StatusCode,
			err,
		)
	}
	return fmt.Errorf(
		"bad status code: %d, response body:%s", resp.StatusCode, string(b))
}

func urlWithAPIKey(u string, apiKey string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("apiKey", apiKey)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func (c *client) startAsyncCalculation(
	ctx context.Context,
	points []route.Point,
	includeDistance, includeDuration bool,
	opts ...MatrixOption,
) (string, error) {
	resp, err := c.calculate(
		ctx, points, true, includeDistance, includeDuration, opts...)
	if err != nil {
		return "", err
	}
	defer func() {
		if tempErr := resp.Body.Close(); tempErr != nil {
			err = tempErr
		}
	}()
	if resp.StatusCode != http.StatusAccepted {
		return "", badStatusError(resp)
	}

	var statusResponse statusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		return "", fmt.Errorf(
			"decoding status response from starting calculation: %v", err)
	}

	statusURL, err := urlWithAPIKey(statusResponse.StatusURL, c.APIKey)
	if err != nil {
		return "", fmt.Errorf("parsing status URL: %v", err)
	}

	return statusURL, nil
}

func (c *client) poll(
	req *http.Request,
) (resultURL string, ready bool, retry bool, err error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", false, true, fmt.Errorf("getting status: %w", err)
	}
	defer func() {
		if tempErr := resp.Body.Close(); tempErr != nil {
			err = tempErr
		}
	}()
	if err != nil {
		return "", false, true, fmt.Errorf("closing response body: %w", err)
	}

	if resp.StatusCode > http.StatusBadRequest &&
		resp.StatusCode < http.StatusInternalServerError {
		return "", false, false, validationError{
			message: fmt.Sprintf("aborting request due to received status code: %d",
				resp.StatusCode),
		}
	}
	if resp.StatusCode >= http.StatusInternalServerError {
		return "", false, true, fmt.Errorf(
			"retry request due to received status code: %d",
			resp.StatusCode)
	}
	var statusResponse statusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		return "", false, true, fmt.Errorf("decoding status response: %v", err)
	}

	resultURL, err = urlWithAPIKey(statusResponse.ResultURL, c.APIKey)
	if err != nil {
		return "", false, false, fmt.Errorf("parsing result URL: %v", err)
	}

	if !isKnownStatusResponse(statusResponse.Status) {
		return "", false, false, fmt.Errorf(
			"unknown status: %s", statusResponse.Status)
	}

	return resultURL, statusResponse.Status == responseStatusComplete, true, nil
}

func (c *client) calculate(
	ctx context.Context,
	points []route.Point,
	async, includeDistance, includeDuration bool,
	opts ...MatrixOption,
) (*http.Response, error) {
	url := fmt.Sprintf(
		"%s/v8/matrix?apiKey=%s&async=%t", c.schemeHost, c.APIKey, async)

	var herePoints []point
	for _, p := range points {
		if len(p) == 2 {
			herePoints = append(herePoints, point{
				Lon: p[0],
				Lat: p[1],
			})
		}
	}

	hereReq := &matrixRequest{
		Origins: herePoints,
		RegionDefinition: regionDefinition{
			Type: "autoCircle",
		},
	}

	if includeDistance {
		hereReq.MatrixAttributes = append(hereReq.MatrixAttributes, "distances")
	}
	if includeDuration {
		hereReq.MatrixAttributes = append(
			hereReq.MatrixAttributes, "travelTimes")
	}
	for _, opt := range opts {
		opt(hereReq)
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(hereReq); err != nil {
		return nil, fmt.Errorf("encoding request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("constructing request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	return resp, nil
}

func reshape(m []int, points []route.Point) [][]float64 {
	// TODO: this can happen when we construct the here points
	// so we don't need to iterate over the matrix again
	widthWithoutZeroes := 0
	for _, p := range points {
		if len(p) == 2 {
			widthWithoutZeroes++
		}
	}

	width := len(points)
	reshaped := make([][]float64, width)
	iZeroes := 0
	for i := 0; i < width; i++ {
		reshaped[i] = make([]float64, width)
		if len(points[i]) == 0 {
			iZeroes++
			continue
		}
		jZeroes := 0
		for j := 0; j < width; j++ {
			if len(points[j]) == 0 {
				jZeroes++
				reshaped[i][j] = 0
			} else {
				reshaped[i][j] = float64(
					m[(i-iZeroes)*widthWithoutZeroes+j-jZeroes])
			}
		}
	}

	return reshaped
}
