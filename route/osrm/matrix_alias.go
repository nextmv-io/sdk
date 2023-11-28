package osrm

import (
	o "github.com/nextmv-io/sdk/measure/osrm"
	"github.com/nextmv-io/sdk/route"
)

// DistanceMatrix makes a request for a distance table from an OSRM server and
// returns a Matrix. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
func DistanceMatrix(
	c Client, points []route.Point,
	parallelQueries int,
) (route.ByIndex, error) {
	return o.DistanceMatrix(c, points, parallelQueries)
}

// DurationMatrix makes a request for a duration table from an OSRM server and
// returns a Matrix. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
func DurationMatrix(
	c Client, points []route.Point,
	parallelQueries int,
) (route.ByIndex, error) {
	return o.DurationMatrix(c, points, parallelQueries)
}

// DistanceDurationMatrices fetches a distance and duration table from an OSRM
// server and returns a Matrix of each. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
func DistanceDurationMatrices(
	c Client,
	points []route.Point,
	parallelQueries int,
) (
	distance, duration route.ByIndex,
	err error,
) {
	return o.DistanceDurationMatrices(c, points, parallelQueries)
}
