package main

import (
	"fmt"
	"strconv"

	"github.com/nextmv-io/sdk/route"
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

// overrideMeasureValues overrides the given measure if start/ends are not given
// with a constant 0 value. Otherwise it returns the same measure again.
func overrideMeasureValues(
	stops,
	starts,
	ends []route.Point,
	vehicleCount int,
	m route.ByIndex,
) route.ByIndex {
	if len(starts) > 0 && len(ends) > 0 {
		return m
	}
	overrideIndices := map[int]bool{}
	if len(starts) == 0 {
		for idx := len(stops); idx < len(stops)+2*vehicleCount; idx += 2 {
			overrideIndices[idx] = true
		}
	}
	if len(ends) == 0 {
		for idx := len(stops) + 1; idx < len(stops)+2*vehicleCount; idx += 2 {
			overrideIndices[idx] = true
		}
	}

	m = route.Override(
		m,
		route.Constant(0),
		func(from, to int) bool {
			_, fromOk := overrideIndices[from]
			_, toOk := overrideIndices[to]
			if fromOk || toOk {
				return true
			}
			return false
		},
	)

	return m
}

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
