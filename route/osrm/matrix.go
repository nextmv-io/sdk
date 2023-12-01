package osrm

import (
	"github.com/nextmv-io/sdk/route"
)

func overrideZeroes(m route.ByIndex, points []route.Point) route.ByIndex {
	return route.Override(m, route.Constant(0.0), func(i, j int) bool {
		return len(points[i]) == 0 || len(points[j]) == 0
	})
}

// DistanceMatrix makes a request for a distance table from an OSRM server and
// returns a Matrix. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/osrm].
func DistanceMatrix(
	c Client, points []route.Point,
	parallelQueries int,
) (route.ByIndex, error) {
	p1, _, err := c.Table(points, WithDistance(), ParallelRuns(parallelQueries))
	if err != nil {
		// preserve the error type for callers
		return nil, err
	}

	return overrideZeroes(route.Matrix(p1), points), nil
}

// DurationMatrix makes a request for a duration table from an OSRM server and
// returns a Matrix. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/osrm].
func DurationMatrix(
	c Client, points []route.Point,
	parallelQueries int,
) (route.ByIndex, error) {
	_, p2, err := c.Table(points, WithDuration(), ParallelRuns(parallelQueries))
	if err != nil {
		// preserve the error type for callers
		return nil, err
	}

	return overrideZeroes(route.Matrix(p2), points), nil
}

// DistanceDurationMatrices fetches a distance and duration table from an OSRM
// server and returns a Matrix of each. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/osrm].
func DistanceDurationMatrices(
	c Client,
	points []route.Point,
	parallelQueries int,
) (
	distance, duration route.ByIndex,
	err error,
) {
	p1, p2, err := c.Table(
		points,
		WithDistance(),
		WithDuration(),
		ParallelRuns(parallelQueries),
	)
	if err != nil {
		// preserve the error type for callers
		return nil, nil, err
	}

	return overrideZeroes(
			route.Matrix(p1),
			points),
		overrideZeroes(route.Matrix(p2),
			points,
		), nil
}
