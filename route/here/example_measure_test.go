package here_test

import (
	"context"
	"fmt"
	"time"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/route/here"
)

// How to create a `Client`, assuming the points are in the form `{longitude,
// latitude}`.
// Please note that this example does not define an output as it requires API
// authentication and is here for illustrative purposes only.
func Example() {
	points := []route.Point{
		{-74.028297, 4.875835},
		{-74.046965, 4.872842},
		{-74.041763, 4.885648},
	}
	client := here.NewClient("<your-api-key>")

	// All of the matrix-constructing functions take a context as their first
	// parameter, which can be used to cancel a request cycle while it is in
	// progress. For example, it is a good general practice to use this context
	// to impose a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Distance and duration matrices can be constructed with the functions
	// provided in the package. These functions will use a synchronous request
	// flow if the number of points requested is below HERE's size limit for
	// synchronous API calls - otherwise, HERE's asynchronous flow will
	// automatically be used.
	dist, dur, err := client.DistanceDurationMatrices(ctx, points)
	if err != nil {
		panic(err)
	}
	_ = dist // We don't use the distance matrix in this example.

	// The matrix functions also take a variadic list of `MatrixOption`s that
	// configure how HERE will calculate the routes. For example, this code
	// constructs a DistanceMeasure with a specific departure time for a
	// bicycle:
	dist, err = client.DistanceMatrix(
		ctx,
		points,
		here.WithTransportMode(here.TransportModeBicycle),
		here.WithDepartureTime(time.Date(2021, 12, 10, 8, 30, 0, 0, time.Local)),
	)
	if err != nil {
		panic(err)
	}
	_ = dist // We don't use the distance matrix in this example.

	// Or, you can configure a truck profile:
	dist, err = client.DistanceMatrix(
		ctx,
		points,
		here.WithTransportMode(here.TransportModeTruck),
		here.WithTruckProfile(here.Truck{
			Type:                  here.TruckTypeTractor,
			TrailerCount:          2,
			ShippedHazardousGoods: []here.HazardousGood{here.Poison},
		}),
	)
	if err != nil {
		panic(err)
	}

	// Once the measures have been created, you may estimate the distances and
	// durations by calling the Cost function.
	for p1 := range points {
		for p2 := range points {
			fmt.Printf(
				"(%d, %d) = [%.2f, %.2f]\n",
				p1, p2, dist.Cost(p1, p2), dur.Cost(p1, p2),
			)
		}
	}
	// This is the expected output.
	// (0, 0) = [0.00, 0.00]
	// (0, 1) = [6519.00, 791.00]
	// (0, 2) = [4884.00, 632.00]
	// (1, 0) = [5242.00, 580.00]
	// (1, 1) = [0.00, 0.00]
	// (1, 2) = [2270.00, 248.00]
	// (2, 0) = [3800.00, 493.00]
	// (2, 1) = [2270.00, 250.00]
	// (2, 2) = [0.00, 0.00]
}
