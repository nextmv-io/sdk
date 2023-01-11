package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/nextmv-io/sdk/route"
	osrm "github.com/nextmv-io/sdk/route/osrm"
)

func main() {
	// Some input data to create matrices for. This input data routes points in
	// Berlin, Germany.
	input := input{
		Stops: []route.Point{
			{53.21077, 13.57843},
			{52.51017, 13.80583},
		},
		Starts: []route.Point{},
		Ends: []route.Point{
			{52.81254, 13.35778},
			{},
		},
		Vehicles: []string{"vehicle-1", "vehicle-2"},
	}

	// Build a slice of points in the correct format to use by a matrix provider.
	points, err := route.BuildRequestSlice(
		input.Stops,
		input.Starts,
		input.Ends,
		len(input.Vehicles),
	)
	if err != nil {
		panic(err)
	}

	// Create an OSRM client
	client := osrm.DefaultClient(
		"<YOUR-OSRM-SERVER-URL>",
		true,
	)
	client.SnapRadius(0)

	// Get distance matrices.
	distanceMatrix, durationMatrix, err := osrm.DistanceDurationMatrices(
		client,
		points,
		0,
	)
	if err != nil {
		panic(err)
	}

	// Override values for empty locations in the input.
	distanceMatrix = overrideMeasureValues(
		input.Stops,
		input.Starts,
		input.Ends,
		len(input.Vehicles),
		distanceMatrix,
	)
	durationMatrix = overrideMeasureValues(
		input.Stops,
		input.Starts,
		input.Ends,
		len(input.Vehicles),
		durationMatrix,
	)

	// Create an output for the routing app in the expected format.
	now := time.Now()
	out := Output{
		Stops:    convertToStop(input.Stops),
		Starts:   convertToPosition(input.Starts),
		Ends:     convertToPosition(input.Ends),
		Vehicles: input.Vehicles,
		Shifts: []route.TimeWindow{
			{
				Start: now,
				End:   now.Add(24 * time.Hour),
			},
			{
				Start: now,
				End:   now.Add(24 * time.Hour),
			},
		},
		DistanceMatrix: makeFloatMatrix(distanceMatrix, len(points)),
		DurationMatrix: makeFloatMatrix(durationMatrix, len(points)),
	}

	// Write the output.
	file, _ := json.MarshalIndent(out, "", " ")
	_ = ioutil.WriteFile("routing-input.json", file, 0644)
}
