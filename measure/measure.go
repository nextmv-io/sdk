package measure

import "errors"

// BuildMatrixRequestPoints builds a slice of points in the correct format to
// request a matrix from any of the supported platforms (e.g. OSRM, Routingkit,
// Google, HERE). It takes the stops to be routed, start and end stops of
// vehicles (optional) and the number of to be used.
func BuildMatrixRequestPoints(
	stops, starts,
	ends []Point,
	vehiclesCount int,
) ([]Point, error) {
	if len(starts) > 0 && len(starts) != vehiclesCount {
		return nil, errors.New(
			"if starts are given, they must match the number of vehicles",
		)
	}
	if len(ends) > 0 && len(ends) != vehiclesCount {
		return nil, errors.New(
			"if ends are given, they must match the number of vehicles",
		)
	}
	count := len(stops)
	// Create points array of the expected size
	points := make([]Point, count+2*vehiclesCount)
	for i := range points {
		// Set default values
		points[i] = Point{0, 0}
	}
	copy(points, stops)

	if len(starts) > 0 {
		for v, start := range starts {
			points[count+v*2] = start
		}
	}

	if len(ends) > 0 {
		for v, end := range ends {
			points[count+v*2+1] = end
		}
	}
	return points, nil
}

// OverrideZeroPoints overrides points that have been passed as placeholders
// [0,0] to build the matrix with zero values.
func OverrideZeroPoints(
	points []Point,
	m ByIndex,
) ByIndex {
	m = Override(
		m,
		Constant(0),
		func(from, to int) bool {
			// If any point is nil or empty (len==0) or [0,0] return true.
			if points[from] == nil || points[to] == nil ||
				len(points[from]) == 0 || len(points[to]) == 0 ||
				(points[from][0] == 0 && points[from][1] == 0) ||
				(points[to][0] == 0 && points[to][1] == 0) {
				return true
			}
			return false
		},
	)

	return m
}
