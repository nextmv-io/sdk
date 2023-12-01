package osrm_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/measure/osrm"
)

// Please note that this example does not define an output as it requires a
// server and is here for illustrative purposes only.
func ExampleDefaultClient() {
	client := osrm.DefaultClient(
		"YOUR-OSRM-SERVER e.g.:http://localhost:5000",
		true,
	)
	_ = client // The client is instantiated as an example, but it is not used.
}

// Please note that this example does not define an output as it requires a
// server and is here for illustrative purposes only.
func ExampleDistanceMatrix() {
	client := osrm.DefaultClient(
		"YOUR-OSRM-SERVER e.g.:http://localhost:5000",
		true,
	)
	points := []measure.Point{
		{-123.1041788, 43.9965908},
		{-123.1117056, 44.0568198},
	}
	dist, err := osrm.DistanceMatrix(client, points, 0)
	if err != nil {
		panic(err)
	}
	cost := dist.Cost(0, 1)
	fmt.Println(cost)
}

// Please note that this example does not define an output as it requires a
// server and is here for illustrative purposes only.
func ExampleDurationMatrix() {
	client := osrm.DefaultClient(
		"YOUR-OSRM-SERVER e.g.:http://localhost:5000",
		true,
	)
	points := []measure.Point{
		{-123.1041788, 43.9965908},
		{-123.1117056, 44.0568198},
	}
	dur, err := osrm.DurationMatrix(client, points, 0)
	if err != nil {
		panic(err)
	}
	cost := dur.Cost(0, 1)
	fmt.Println(cost)
}

// Please note that this example does not define an output as it requires a
// server and is here for illustrative purposes only.
func ExampleDistanceDurationMatrices() {
	client := osrm.DefaultClient(
		"YOUR-OSRM-SERVER e.g.:http://localhost:5000",
		true,
	)
	points := []measure.Point{
		{-123.1041788, 43.9965908},
		{-123.1117056, 44.0568198},
	}
	dist, dur, err := osrm.DistanceDurationMatrices(client, points, 0)
	if err != nil {
		panic(err)
	}
	costDist := dist.Cost(0, 1)
	costDur := dur.Cost(0, 1)
	fmt.Println(costDist)
	fmt.Println(costDur)
}
