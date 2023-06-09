// package main holds the implementation of the measure generation.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/route/osrm"
)

type input struct {
	Stops    []route.Point
	Starts   []route.Point
	Ends     []route.Point
	Vehicles []string
}

type output struct {
	Stops          []route.Stop       `json:"stops"`
	Vehicles       []string           `json:"vehicles"`
	Starts         []route.Position   `json:"starts"`
	Ends           []route.Position   `json:"ends"`
	Shifts         []route.TimeWindow `json:"shifts"`
	DurationMatrix [][]float64        `json:"duration_matrix"`
	DistanceMatrix [][]float64        `json:"distance_matrix"`
}

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
	points, err := route.BuildMatrixRequestPoints(
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
	err = client.SnapRadius(0)
	if err != nil {
		panic(err)
	}

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
	distanceMatrix = route.OverrideZeroPoints(
		points,
		distanceMatrix,
	)
	durationMatrix = route.OverrideZeroPoints(
		points,
		durationMatrix,
	)

	// Create an output for the routing app in the expected format.
	now := time.Date(2022, 10, 17, 9, 0, 0, 0, time.UTC)
	out := output{
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

	/* #nosec */
	err = os.WriteFile("routing-input.json", file, 0o644)
	if err != nil {
		panic(err)
	}
}

// makeFloatMatrix is a helper function that takes a route.ByIndex and returns a
// [][]float64.
func makeFloatMatrix(matrix route.ByIndex, length int) [][]float64 {
	out := make([][]float64, length)
	for i := 0; i < length; i++ {
		out[i] = make([]float64, length)
		for j := 0; j < length; j++ {
			out[i][j] = matrix.Cost(i, j)
		}
	}
	return out
}

// convertToStop is a helper function that converts a []route.Point into a
// []route.Stop.
func convertToStop(points []route.Point) []route.Stop {
	stops := make([]route.Stop, len(points))
	for i, p := range points {
		stops[i] = route.Stop{
			ID: fmt.Sprintf("stop-%s", strconv.Itoa(i)),
			Position: route.Position{
				Lon: p[1],
				Lat: p[0],
			},
		}
	}
	return stops
}

// convertToPosition is a helper function that converts a []route.Point into a
// []route.Position.
func convertToPosition(points []route.Point) []route.Position {
	position := make([]route.Position, len(points))
	for i, p := range points {
		if len(p) == 0 {
			return []route.Position{}
		}
		position[i] = route.Position{
			Lon: p[1],
			Lat: p[0],
		}
	}
	return position
}
