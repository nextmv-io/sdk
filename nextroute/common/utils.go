package common

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Filter filters a slice using a predicate function.
func Filter[T any](v []T, f func(T) bool) []T {
	r := make([]T, 0, len(v))
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

// UniqueDefined is a universal duplicate removal function for type instances in
// a slice that implement the comparable interface. The function f is used to
// extract the comparable value from the type instance.
func UniqueDefined[T any, I comparable](items []T, f func(T) I) []T {
	inResult := make(map[I]bool)
	var result []T
	for _, item := range items {
		if _, ok := inResult[f(item)]; !ok {
			inResult[f(item)] = true
			result = append(result, item)
		}
	}
	return result
}

// GroupBy groups the elements of a slice by a key function.
func GroupBy[T any, K comparable](s []T, f func(T) K) map[K][]T {
	inResult := make(map[K]bool)
	result := make(map[K][]T)
	for _, instance := range s {
		key := f(instance)
		if _, ok := inResult[key]; !ok {
			result[key] = make([]T, 0)
		}
		result[key] = append(result[key], instance)
	}
	return result
}

// Map maps a slice of type T to a slice of type R using the function f.
func Map[T any, R any](v []T, f func(T) R) []R {
	r := make([]R, len(v))
	for idx, x := range v {
		r[idx] = f(x)
	}
	return r
}

// MapSlice maps a slice of type T to a slice of type R using the function f
// returning a slice of R.
func MapSlice[T any, R any](v []T, f func(T) []R) []R {
	var r []R
	for _, x := range v {
		r = append(r, f(x)...)
	}
	return r
}

// FindIndex returns the first index i satisfying predicate(s[i]).
func FindIndex[E any](s []E, predicate func(E) bool) int {
	for i, v := range s {
		if predicate(v) {
			return i
		}
	}
	return -1
}

// AllTrue returns true if all the given predicate evaluations are true.
func AllTrue[E any](s []E, predicate func(E) bool) bool {
	return All(s, true, predicate)
}

// AllFalse returns true if all the given predicate evaluations is false.
func AllFalse[E any](s []E, predicate func(E) bool) bool {
	return All(s, false, predicate)
}

// All returns true if all the given predicate evaluations evaluate to
// condition.
func All[E any](s []E, condition bool, predicate func(E) bool) bool {
	for _, v := range s {
		if predicate(v) != condition {
			return false
		}
	}
	return true
}

// HasTrue returns true if any of the given predicate evaluations evaluate to
// true.
func HasTrue[E any](s []E, predicate func(E) bool) bool {
	return Has(s, true, predicate)
}

// HasFalse returns true if any of the given predicate evaluations evaluate to
// false.
func HasFalse[E any](s []E, predicate func(E) bool) bool {
	return Has(s, false, predicate)
}

// Has returns true if any of the given predicate evaluations is condition.
func Has[E any](s []E, condition bool, predicate func(E) bool) bool {
	for _, v := range s {
		if predicate(v) == condition {
			return true
		}
	}
	return false
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
	seconds := distance.Value(Meters) / speed.Value(MetersPerSecond)

	return seconds / timeUnit.Seconds()
}

// RandomElement returns a random element from the given slice. If the slice is
// empty, panic is raised. If source is nil, a new source is created using the
// current time.
func RandomElement[T any](
	source *rand.Rand,
	elements []T,
) T {
	if len(elements) == 0 {
		panic(fmt.Errorf("cannot select random element from empty slice"))
	}
	if source == nil {
		source = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return elements[source.Intn(len(elements))]
}

// RandomElements returns a slice of n random elements from the
// given slice. If n is greater than the length of the slice, all elements are
// returned. If n is less than or equal to zero, an empty slice is returned.
// If source is nil, a new source is created using the current time.
func RandomElements[T any](
	source *rand.Rand,
	elements []T,
	n int,
) []T {
	if source == nil {
		source = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if n <= 0 {
		return []T{}
	}
	if n >= len(elements) {
		return elements
	}
	result := make([]T, 0, n)
	indicesUsed := make(map[int]bool, 0)
	for i := 0; i < n; i++ {
		index := RandomIndex(
			source,
			len(elements),
			indicesUsed,
		)
		result = append(result, elements[index])
	}

	return result
}

// RandomElementIndices returns a slice of n random element indices from the
// given slice. If n is greater than the length of the slice, all indices are
// returned. If n is less than or equal to zero, an empty slice is returned.
// If source is nil, a new source is created using the current time.
func RandomElementIndices[T any](
	source *rand.Rand,
	elements []T,
	n int,
) []int {
	if source == nil {
		source = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if n <= 0 {
		return []int{}
	}
	if n >= len(elements) {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[n] = i
		}
		return result
	}
	result := make([]int, 0, n)
	indicesUsed := make(map[int]bool, 0)
	for i := 0; i < n; i++ {
		index := RandomIndex(
			source,
			len(elements),
			indicesUsed,
		)
		result = append(result, index)
	}

	return result
}

// RandomIndex returns a random index from the given size. If the index has
// already been used, a new index is generated. If source is nil, a new source
// is created using the current time.
func RandomIndex(source *rand.Rand, size int, indicesUsed map[int]bool) int {
	if source == nil {
		source = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	for {
		index := source.Intn(size)
		if used, ok := indicesUsed[index]; !ok || !used {
			indicesUsed[index] = true
			return index
		}
	}
}

// Lazy is a function that returns a value of type T. The value is only
// calculated once, and the result is cached for subsequent calls.
type Lazy[T any] func() T

// DefineLazy returns a Lazy[T] that will call the given function to
// calculate the value. The value is only calculated once, and the result
// is cached for subsequent calls.
func DefineLazy[T any](f func() T) Lazy[T] {
	var value T
	var once sync.Once
	return func() T {
		once.Do(func() {
			value = f()
			f = nil
		})
		return value
	}
}
