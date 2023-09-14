package main

import (
	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/common"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/route"
	"github.com/twpayne/go-polyline"
)

func newDistanceExpression(input schema.FleetInput) nextroute.DistanceExpression {
	points := make([]measure.Point, len(input.Stops)+2*len(input.Vehicles))
	for i, s := range input.Stops {
		points[i] = measure.Point{
			s.Position.Lon, s.Position.Lat,
		}
	}

	counter := 0
	for i := len(input.Stops); i < len(input.Stops); i += 2 {
		if input.Vehicles[counter].Start != nil {
			points[i] = measure.Point{
				input.Vehicles[counter].Start.Lon, input.Vehicles[counter].Start.Lat,
			}
		}

		if input.Vehicles[counter].End != nil {
			points[i+1] = measure.Point{
				input.Vehicles[counter].End.Lon, input.Vehicles[counter].End.Lat,
			}
		}
		counter++
	}

	m := measure.HaversineByPoint()
	mIndexed := measure.Indexed(m, points)

	distanceMatrix := newFloatMatrix(mIndexed, len(points))
	distanceExpression := distanceExpression(&distanceMatrix)
	return distanceExpression
}

// distanceExpression creates a distance expression for later use.
func distanceExpression(distanceMatrix *[][]float64) nextroute.DistanceExpression {
	distanceExpression := nextroute.NewHaversineExpression()
	if distanceMatrix != nil {
		distanceExpression = nextroute.NewDistanceExpression(
			"travelDistance",
			nextroute.NewMeasureByIndexExpression(measure.Matrix(*distanceMatrix)),
			common.Meters,
		)
	}
	return distanceExpression
}

type distanceData struct {
	distance nextroute.DistanceExpression
}

// newFloatMatrix is a helper function that takes a route.ByIndex and returns a
// [][]float64.
func newFloatMatrix(matrix route.ByIndex, length int) [][]float64 {
	out := make([][]float64, length)
	for i := 0; i < length; i++ {
		out[i] = make([]float64, length)
		for j := 0; j < length; j++ {
			out[i][j] = matrix.Cost(i, j)
		}
	}
	return out
}

func haversinePolyline(p []measure.Point) (string, []string) {
	if len(p) < 2 {
		return "", []string{}
	}
	legs := make([]string, len(p))
	allCoords := make([][]float64, len(p))
	for i := 0; i < len(p); i++ {
		current := p[i]
		coord := []float64{
			current[1], current[0],
		}
		allCoords[i] = coord
		// the last stop does not have an individual leg
		if i == len(p)-1 {
			break
		}
		next := p[i+1]
		coords := [][]float64{
			{current[1], current[0]},
			{next[1], next[0]},
		}
		leg := string(polyline.EncodeCoords(coords))
		legs[i] = leg
	}
	routePolyline := string(polyline.EncodeCoords(allCoords))
	return routePolyline, legs
}

// mapWithError maps a slice of type T to a slice of type R using the function f.
func mapWithError[T any, R any](v []T, f func(T) (R, error)) ([]R, error) {
	r := make([]R, len(v))
	for idx, x := range v {
		o, err := f(x)
		if err != nil {
			return r, err
		}
		r[idx] = o
	}
	return r, nil
}
