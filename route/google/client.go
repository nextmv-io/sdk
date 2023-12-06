package google

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/nextmv-io/sdk/route"
	"googlemaps.github.io/maps"
)

const (
	maxAddresses         int = 25
	maxElements          int = 100
	maxElementsPerSecond int = 1000
)

// reference references a slice of the complete origins or destinations.
type reference struct {
	start int
	num   int
}

// matrixRequest represents a reference to a Google Maps request with additional
// information for allocating resulting information in the correct matrix
// indices.
type matrixRequest struct {
	r            *maps.DistanceMatrixRequest
	origins      reference
	destinations reference
}

// matrixResponse represents a reference to a Google Maps response with
// additional information for allocating resulting information in the correct
// matrix indices.
type matrixResponse struct {
	r            *maps.DistanceMatrixResponse
	origins      reference
	destinations reference
}

// DistanceDurationMatrices makes requests to the Google Distance Matrix API
// and returns route.ByIndex types to estimate distances (in meters) and
// durations (in seconds). It receives a Google Maps Client and Request. The
// coordinates passed to the request must be in the form latitude, longitude.
// The resulting distance and duration matrices are saved in memory. To find
// out how to create a client and request, please refer to the go package docs.
// This function takes into consideration the usage limits of the Distance
// Matrix API and thus may transform the request into multiple ones and handle
// them accordingly. You can find more about usage limits here in the official
// google maps documentation for the distance matrix, usage and billing.
//
// Deprecated: This package is deprecated and will be removed in the future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/google].
func DistanceDurationMatrices(c *maps.Client, r *maps.DistanceMatrixRequest) (
	route.ByIndex,
	route.ByIndex,
	error,
) {
	// Split request into multiple ones respecting the maxAddresses and
	// maxElements usage limits.
	requests := split(r)

	// Group requests that can go together respecting the maxElementsPerSecond
	// limit. Sequential calls are performed, where each call is composed of
	// concurrent requests to the API. Sleep between concurrent calls to avoid
	// exceeding the usage limit.
	groups := group(requests)
	var responses []matrixResponse
	for _, group := range groups {
		resp, err := makeMatrixRequest(c, group)
		if err != nil {
			return nil, nil, err
		}
		responses = append(responses, resp...)
		time.Sleep(1 * time.Second)
	}

	// Instantiate an empty numOrigins x numDestinations matrix.
	numOrigins, numDestinations := len(r.Origins), len(r.Destinations)
	distances := make([][]float64, numOrigins)
	durations := make([][]float64, numOrigins)
	for i := 0; i < numOrigins; i++ {
		distances[i] = make([]float64, numDestinations)
		durations[i] = make([]float64, numDestinations)
	}

	// The responses are processed to fill the distance and duration matrices.
	for _, response := range responses {
		for i, row := range response.r.Rows {
			from := i + response.origins.start
			for j, element := range row.Elements {
				to := j + response.destinations.start
				distances[from][to] = float64(element.Distance.Meters)
				durations[from][to] = element.Duration.Seconds()
			}
		}
	}

	return route.Matrix(distances), route.Matrix(durations), nil
}

/*
split splits a single request into multiple ones following these rules:

  - there can only be a maximum of maxAddresses origins or destinations,
  - there can only be a maximum of maxElements elements.

For each new request, a reference of the indices of the original origins and
destinations is saved to allocate the correct positions of the matrix.
*/
func split(r *maps.DistanceMatrixRequest) []matrixRequest {
	var requests []matrixRequest
	numOrigins, numDestinations := len(r.Origins), len(r.Destinations)

	// Origins are looped over first.
	originsStart := 0
	originsRemaining := numOrigins
	for originsRemaining > 0 {
		// At a maximum we take maxAddresses origins at a time and create a
		// temporary slice.
		originsIncrease := int(math.Min(
			float64(maxAddresses),
			float64(originsRemaining)),
		)
		originsEnd := originsStart + originsIncrease
		originsSliced := r.Origins[originsStart:originsEnd]

		// For each batch of origins, we adjust the number of destinations that
		// we will be slicing.
		destinationsStart := 0
		destinationsRemaining := numDestinations
		for destinationsRemaining > 0 {
			// At a maximum we take the number of destinations that result in
			// maxElements elements and create a temporary slice.
			destinationsIncrease := int(math.Min(
				float64(maxElements)/float64(originsIncrease),
				float64(destinationsRemaining)),
			)
			destinationsEnd := destinationsStart + destinationsIncrease
			destinationsSliced := r.Destinations[destinationsStart:destinationsEnd]

			// The sliced origins and destinations are used to create a new
			// request that has all the other attributes of the original one.
			req := *r
			req.Origins = originsSliced
			req.Destinations = destinationsSliced
			request := matrixRequest{
				r: &req,
				origins: reference{
					start: originsStart,
					num:   originsIncrease,
				},
				destinations: reference{
					start: destinationsStart,
					num:   destinationsIncrease,
				},
			}
			requests = append(requests, request)

			destinationsStart = destinationsEnd
			destinationsRemaining -= destinationsIncrease
		}

		originsStart = originsEnd
		originsRemaining -= originsIncrease
	}

	return requests
}

// group groups requests that can be called concurrently respecting the
// maxElementsPerSecond limit.
func group(requests []matrixRequest) [][]matrixRequest {
	// Create groups as long as there are requests left to process.
	var groups [][]matrixRequest
	remaining := len(requests)
	for remaining > 0 {
		// Create a group as long as the total elements do not exceed the limit.
		var group []matrixRequest
		remainingElements := maxElementsPerSecond
		for i := len(requests) - remaining; i < len(requests); i++ {
			elements := requests[i].origins.num * requests[i].destinations.num
			if remainingElements-elements >= 0 {
				// The element count is verified to check if the request can be
				// added to the current group. If the request is added, the
				// overall element and request count is updated.
				group = append(group, requests[i])
				remainingElements -= elements
				remaining--
			} else {
				// If the request cannot be added, the current group is complete
				// and a new group is created all over again.
				break
			}
		}

		groups = append(groups, group)
	}

	return groups
}

// matrixResult gathers a response and possible error from concurrent requests.
type matrixResult struct {
	res *matrixResponse
	err error
}

// makeMatrixRequest performs concurrent requests to the provided client.
func makeMatrixRequest(
	c *maps.Client,
	group []matrixRequest,
) ([]matrixResponse, error) {
	// Define channels to gather results and quit other goroutines in case of an
	// error.
	out := make(chan matrixResult, len(group))
	// Perform concurrent requests.
	for _, req := range group {
		go func(req matrixRequest) {
			// Actually make the request.
			r, err := c.DistanceMatrix(context.Background(), req.r)
			if err != nil {
				// If the request errors, push the error to the chan and signal
				// closure.
				out <- matrixResult{res: nil, err: err}
			} else {
				// If the request does not error, push the request to the chan.
				response := matrixResponse{
					r:            r,
					origins:      req.origins,
					destinations: req.destinations,
				}
				out <- matrixResult{res: &response, err: nil}
			}
		}(req)
	}

	// Empty out the chan to gather responses. If there is an error found,
	// immediately return it. Loop over the requests in the group to have
	// control over the number of times we are getting an element from the chan
	// but the resulting responses in the out chan are not expected to have the
	// same order as the requests in group.
	var responses []matrixResponse
	for i := 0; i < len(group); i++ {
		result := <-out
		if result.err != nil {
			return nil, result.err
		}
		responses = append(responses, *result.res)
	}

	return responses, nil
}

type directionRequest struct {
	r     *maps.DirectionsRequest
	index int
}

type directionResponse struct {
	r     []maps.Route
	index int
}

// directionsResult gathers a response and possible error from concurrent
// requests.
type directionsResult struct {
	res *directionResponse
	err error
}

// makeDistanceRequest performs concurrent requests to the provided client.
func makeDirectionsRequest(
	c *maps.Client,
	points []string,
	orgRequest *maps.DirectionsRequest,
) ([]directionResponse, error) {
	// Loop over points to make directions requests taking into account Google's
	// waypoint limitation per request.
	remaining := len(points) - 1
	blockSize := 25
	start := 0
	// Determine the number of requests to be made.
	numberRequests := int(math.Ceil(float64(remaining) / float64(blockSize)))
	directionRequests := make([]directionRequest, numberRequests)
	for i := range directionRequests {
		count := int(math.Min(float64(blockSize), float64(remaining)))
		end := start + count

		if count <= 0 {
			break
		}

		newReq := *orgRequest
		newReq.Origin = points[start]
		newReq.Destination = points[end]
		// To make the waypoints we need the points to end-1, because the end
		// point is the destination for the request.
		waypoints := make([]string, end-1-start)
		for j := start; j < end-1; j++ {
			waypoints[j-start] = points[j+1]
		}
		newReq.Waypoints = waypoints

		start = end
		remaining -= count
		directionRequests[i] = directionRequest{
			r:     &newReq,
			index: i,
		}
	}

	// Define channels to gather results and quit other goroutines in case of an
	// error.
	out := make(chan directionsResult, len(directionRequests))
	// Perform concurrent requests.
	for _, req := range directionRequests {
		go func(req directionRequest) {
			// Actually make the request.
			r, _, err := c.Directions(context.Background(), req.r)
			if err != nil {
				// If the request errors, push the error to the chan and signal
				// closure.
				out <- directionsResult{res: nil, err: err}
			} else {
				response := directionResponse{
					r:     r,
					index: req.index,
				}
				// If the request does not error, push the request to the chan.
				out <- directionsResult{res: &response, err: nil}
			}
		}(req)
	}

	// Empty out the chan to gather responses. If there is an error found,
	// immediately return it.
	var responses []directionResponse
	for i := 0; i < len(directionRequests); i++ {
		result := <-out
		if result.err != nil {
			return nil, result.err
		}
		responses = append(responses, *result.res)
	}

	// Sort the responses by index, that orders the row packs correctly.
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].index < responses[j].index
	})

	return responses, nil
}

// Polylines requests polylines for the given points. The first parameter
// returns a polyline from start to end and the second parameter returns a list
// of polylines, one per leg.
//
// Deprecated: This package is deprecated and will be removed in the future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/google].
func Polylines(
	c *maps.Client,
	orgRequest *maps.DirectionsRequest,
) (string, []string, error) {
	// Extract all points from the given request
	points := make([]string, len(orgRequest.Waypoints)+2)
	points[0] = orgRequest.Origin
	points[len(points)-1] = orgRequest.Destination
	for i, w := range orgRequest.Waypoints {
		points[i+1] = w
	}

	// Make requests to Google and retrieve results
	responses, err := makeDirectionsRequest(c, points, orgRequest)
	if err != nil {
		return "", []string{}, err
	}

	// the number of total legs for the original request.
	decodedLegs := make([][]maps.LatLng, len(points)-1)

	// Stich results together.
	index := -1
	for _, resp := range responses {
		for _, route := range resp.r {
			for i, leg := range route.Legs {
				index++
				for _, steps := range leg.Steps {
					dec, err := steps.Polyline.Decode()
					if err != nil {
						return "", []string{}, err
					}
					decodedLegs[index] = append(decodedLegs[i], dec...) //nolint:gocritic
				}
			}
		}
	}

	// Make a list of encoded legs to return
	legLines := make([]string, len(points)-1)
	for i, leg := range decodedLegs {
		legLines[i] = maps.Encode(leg)
	}

	// Finally, make a single request to get the polyline from start to end.
	completeReq := *orgRequest
	completeReq.Waypoints = nil
	completeResp, _, err := c.Directions(context.Background(), &completeReq)
	if err != nil {
		return "", []string{}, err
	}

	return completeResp[0].OverviewPolyline.Points, legLines, nil
}
