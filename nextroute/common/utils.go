package common

import (
	"fmt"
	"math"
	"time"
)

// Filter filters a slice using a predicate function.
func Filter[T any](v []T, f func(T) bool) []T {
	var r []T
	for _, x := range v {
		if f(x) {
			r = append(r, x)
		}
	}
	return r
}

// DefensiveCopy returns a defensive copy of a slice.
func DefensiveCopy[T any](v []T) []T {
	c := make([]T, len(v))
	copy(c, v)
	return c
}

// WithinTolerance returns true if a and b are within the given tolerance.
func WithinTolerance(a, b, tolerance float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	if b == 0 {
		return d < tolerance
	}
	return (d / math.Abs(b)) < tolerance
}

// DurationValue returns the value of a duration in the given time unit.
// Will panic if the time unit is zero.
func DurationValue(
	distance Distance,
	speed Speed,
	timeUnit time.Duration,
) float64 {
	if timeUnit.Seconds() == 0 {
		panic(
			fmt.Errorf(
				"time unit is zero in duration calculation",
			),
		)
	}
	seconds := distance.Value(Meters) / speed.Value(NewMetersPerSecond())

	return seconds / timeUnit.Seconds()
}
