package osrm

import (
	"github.com/nextmv-io/sdk/measure"
)

func overrideZeroes(m measure.ByIndex, points []measure.Point) measure.ByIndex {
	return measure.Override(m, measure.Constant(0.0), func(i, j int) bool {
		return len(points[i]) == 0 || len(points[j]) == 0
	})
}

// DistanceMatrix makes a request for a distance table from an OSRM server and
// returns a Matrix. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
func DistanceMatrix(
	c Client, points []measure.Point,
	parallelQueries int,
) (measure.ByIndex, error) {
	p1, _, err := c.Table(points, WithDistance(), ParallelRuns(parallelQueries))
	if err != nil {
		// preserve the error type for callers
		return nil, err
	}

	return overrideZeroes(measure.Matrix(p1), points), nil
}

// DurationMatrix makes a request for a duration table from an OSRM server and
// returns a Matrix. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
func DurationMatrix(
	c Client, points []measure.Point,
	parallelQueries int,
) (measure.ByIndex, error) {
	_, p2, err := c.Table(points, WithDuration(), ParallelRuns(parallelQueries))
	if err != nil {
		// preserve the error type for callers
		return nil, err
	}

	return overrideZeroes(measure.Matrix(p2), points), nil
}

// DistanceDurationMatrices fetches a distance and duration table from an OSRM
// server and returns a Matrix of each. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
func DistanceDurationMatrices(
	c Client,
	points []measure.Point,
	parallelQueries int,
) (
	distance, duration measure.ByIndex,
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
			measure.Matrix(p1),
			points),
		overrideZeroes(measure.Matrix(p2),
			points,
		), nil
}
