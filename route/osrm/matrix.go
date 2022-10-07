// © 2019-2022 nextmv.io inc. All rights reserved.
// nextmv.io, inc. CONFIDENTIAL
//
// This file includes unpublished proprietary source code of nextmv.io, inc.
// The copyright notice above does not evidence any actual or intended
// publication of such source code. Disclosure of this source code or any
// related proprietary information is strictly prohibited without the express
// written permission of nextmv.io, inc.

package osrm

import (
	"fmt"

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
func DistanceMatrix(
	c Client, points []route.Point,
	parallelQueries int,
) (route.ByIndex, error) {
	p1, _, err := c.Table(points, WithDistance(), ParallelRuns(parallelQueries))
	if err != nil {
		return nil, fmt.Errorf("fetching matrix: %v", err)
	}

	return overrideZeroes(route.Matrix(p1), points), nil
}

// DurationMatrix makes a request for a duration table from an OSRM server and
// returns a Matrix. ParallelQueries specifies the number of
// parallel queries to be made, pass 0 to calculate a default, based on the
// number of points given.
func DurationMatrix(
	c Client, points []route.Point,
	parallelQueries int,
) (route.ByIndex, error) {
	_, p2, err := c.Table(points, WithDuration(), ParallelRuns(parallelQueries))
	if err != nil {
		return nil, fmt.Errorf("fetching matrix: %v", err)
	}

	return overrideZeroes(route.Matrix(p2), points), nil
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
	p1, p2, err := c.Table(
		points,
		WithDistance(),
		WithDuration(),
		ParallelRuns(parallelQueries),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching matrices: %v", err)
	}

	return overrideZeroes(
			route.Matrix(p1),
			points),
		overrideZeroes(route.Matrix(p2),
			points,
		), nil
}
