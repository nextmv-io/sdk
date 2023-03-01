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

// Unique is a universal duplicate removal function for type instances in
// a slice that implement the comparable interface.
func Unique[T comparable](s []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

// FindIndex returns the first index i satisfying predicate(s[i]),
func FindIndex[E any](s []E, predicate func(E) bool) int {
	for i, v := range s {
		if predicate(v) {
			return i
		}
	}
	return -1
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
