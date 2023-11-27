package common

import (
	"math/rand"
	"time"

	c "github.com/nextmv-io/sdk/common"
)

// Alias is an interface that allows for sampling from a discrete
// distribution in O(1) time.
type Alias = c.Alias

// NewAlias creates a new Alias from the given weights.
// The weights must be positive and at least one weight must be given.
// The weights are normalized to sum to 1.
// NewAlias([]float64{1, 2, 3}) will return an Alias that will
// return 0 with probability 1/6, 1 with probability 1/3 and 2 with
// probability 1/2.
func NewAlias(weights []float64) (Alias, error) {
	return c.NewAlias(weights)
}

// NewLocation creates a new Location. An error is returned if the longitude is
// not between (-180, 180) or the latitude is not between (-90, 90).
func NewLocation(longitude float64, latitude float64) (Location, error) {
	return c.NewLocation(longitude, latitude)
}

// NewInvalidLocation creates a new invalid Location. Longitude and latitude
// are not important.
func NewInvalidLocation() Location {
	return c.NewInvalidLocation()
}

// Locations is a slice of Location.
type Locations = c.Locations

// Location represents a physical location on the earth.
type Location = c.Location

// BoundingBox contains information about a box.
type BoundingBox = c.BoundingBox

// NewInvalidBoundingBox returns an invalid bounding box.
func NewInvalidBoundingBox() BoundingBox {
	return c.NewInvalidBoundingBox()
}

// NewBoundingBox returns a bounding box for the given locations.
func NewBoundingBox(locations Locations) BoundingBox {
	return c.NewBoundingBox(locations)
}

// DistanceUnit is the unit of distance.
type DistanceUnit = c.DistanceUnit

// Distance is a distance in a given unit.
type Distance = c.Distance

// NewDistance returns a new distance.
func NewDistance(
	value float64,
	unit DistanceUnit,
) Distance {
	return c.NewDistance(value, unit)
}

// DurationUnit is the unit of duration.
type DurationUnit = c.DurationUnit

// FastHaversine is a fast approximation of the haversine distance.
type FastHaversine = c.FastHaversine

// NewFastHaversine returns a new FastHaversine.
func NewFastHaversine(lat float64) FastHaversine {
	return c.NewFastHaversine(lat)
}

// Haversine calculates the distance between two locations using the
// Haversine formula. Haversine is a good approximation for short
// distances (up to a few hundred kilometers).
func Haversine(from, to Location) (Distance, error) {
	return c.Haversine(from, to)
}

// Intersect returns the intersection of two slices.
func Intersect[T comparable](a []T, b []T) []T {
	return c.Intersect[T](a, b)
}

// NSmallest returns the n-smallest items in the slice items using the
// function f to determine the value of each item. If n is greater than
// the length of items, all items are returned.
func NSmallest[T any](items []T, f func(T) float64, n int) []T {
	return c.NSmallest[T](items, f, n)
}

// Meters is the distance travelled by light in a vacuum in
// 1/299,792,458 seconds.
const Meters = c.Meters

// Kilometers is 1000 meters.
const Kilometers = c.Kilometers

// Miles is 1609.34 meters.
const Miles = c.Miles

// KilometersPerHour is a speed unit of kilometers per hour.
var KilometersPerHour = c.KilometersPerHour

// MetersPerSecond is a speed unit of meters per second.
var MetersPerSecond = c.MetersPerSecond

// MilesPerHour is a speed unit of miles per hour.
var MilesPerHour = c.MilesPerHour

const (
	// NanoSecond is 1/1,000,000,000 of a second.
	NanoSecond = c.NanoSecond
	// MicroSecond is 1/1,000,000 of a second.
	MicroSecond = c.MicroSecond
	// MilliSecond is 1/1,000 of a second.
	MilliSecond = c.MilliSecond
	// Second is the SI unit of time.
	Second = c.Second
	// Minute is 60 seconds.
	Minute = c.Minute
	// Hour is 60 minutes.
	Hour = c.Hour
	// Day is 24 hours.
	Day = c.Day
)

// Speed is the interface for a speed.
type Speed = c.Speed

// SpeedUnit represents a unit of speed.
type SpeedUnit = c.SpeedUnit

// NewSpeed creates a new speed.
func NewSpeed(
	distance float64,
	unit SpeedUnit,
) Speed {
	return c.NewSpeed(distance, unit)
}

// NewSpeedUnit returns a new speed unit.
func NewSpeedUnit(
	distanceUnit DistanceUnit,
	duration time.Duration,
) SpeedUnit {
	return c.NewSpeedUnit(distanceUnit, duration)
}

// Statistics describes statistics.
type Statistics = c.Statistics

// NewStatistics creates a new statistics object.
func NewStatistics[T any](v []T, f func(T) float64) Statistics {
	return c.NewStatistics[T](v, f)
}

// Filter filters a slice using a predicate function.
func Filter[T any](v []T, f func(T) bool) []T {
	return c.Filter[T](v, f)
}

// Unique is a universal duplicate removal function for type instances in
// a slice that implement the comparable interface.
func Unique[T comparable](s []T) []T {
	return c.Unique[T](s)
}

// UniqueDefined is a universal duplicate removal function for type instances in
// a slice that implement the comparable interface. The function f is used to
// extract the comparable value from the type instance.
func UniqueDefined[T any, I comparable](items []T, f func(T) I) []T {
	return c.UniqueDefined[T, I](items, f)
}

// NotUnique returns the duplicate instances.
func NotUnique[T comparable](s []T) []T {
	return c.NotUnique[T](s)
}

// NotUniqueDefined returns the instances for which f returns identical
// values.
func NotUniqueDefined[T any, I comparable](items []T, f func(T) I) []T {
	return c.NotUniqueDefined[T, I](items, f)
}

// GroupBy groups the elements of a slice by a key function.
func GroupBy[T any, K comparable](s []T, f func(T) K) map[K][]T {
	return c.GroupBy[T, K](s, f)
}

// Map maps a slice of type T to a slice of type R using the function f.
func Map[T any, R any](v []T, f func(T) R) []R {
	return c.Map[T, R](v, f)
}

// MapSlice maps a slice of type T to a slice of type R using the function f
// returning a slice of R.
func MapSlice[T any, R any](v []T, f func(T) []R) []R {
	return c.MapSlice[T, R](v, f)
}

// FindIndex returns the first index i satisfying predicate(s[i]).
func FindIndex[E any](s []E, predicate func(E) bool) int {
	return c.FindIndex[E](s, predicate)
}

// AllTrue returns true if all the given predicate evaluations are true.
func AllTrue[E any](s []E, predicate func(E) bool) bool {
	return c.AllTrue[E](s, predicate)
}

// AllFalse returns true if all the given predicate evaluations is false.
func AllFalse[E any](s []E, predicate func(E) bool) bool {
	return c.AllFalse[E](s, predicate)
}

// All returns true if all the given predicate evaluations evaluate to
// condition.
func All[E any](s []E, condition bool, predicate func(E) bool) bool {
	return c.All[E](s, condition, predicate)
}

// HasTrue returns true if any of the given predicate evaluations evaluate to
// true.
func HasTrue[E any](s []E, predicate func(E) bool) bool {
	return c.HasTrue[E](s, predicate)
}

// HasFalse returns true if any of the given predicate evaluations evaluate to
// false.
func HasFalse[E any](s []E, predicate func(E) bool) bool {
	return c.HasFalse[E](s, predicate)
}

// Has returns true if any of the given predicate evaluations is condition.
func Has[E any](s []E, condition bool, predicate func(E) bool) bool {
	return c.Has[E](s, condition, predicate)
}

// CopyMap copies all key/value pairs in source adding them to destination.
// If a key exists in both maps, the value in destination is overwritten.
func CopyMap[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](
	destination M1,
	source M2) {
	c.CopyMap[M1, M2, K, V](destination, source)
}

// DefensiveCopy returns a defensive copy of a slice.
func DefensiveCopy[T any](v []T) []T {
	return c.DefensiveCopy[T](v)
}

// WithinTolerance returns true if a and b are within the given tolerance.
func WithinTolerance(a, b, tolerance float64) bool {
	return c.WithinTolerance(a, b, tolerance)
}

// Truncate truncates a float64 to the given unit.
func Truncate(f float64, unit float64) float64 {
	return c.Truncate(f, unit)
}

// DurationValue returns the value of a duration in the given time unit.
// Will panic if the time unit is zero.
func DurationValue(
	distance Distance,
	speed Speed,
	timeUnit time.Duration,
) float64 {
	return c.DurationValue(distance, speed, timeUnit)
}

// RandomElement returns a random element from the given slice. If the slice is
// empty, panic is raised. If source is nil, a new source is created using the
// current time.
func RandomElement[T any](
	source *rand.Rand,
	elements []T,
) T {
	return c.RandomElement[T](source, elements)
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
	return c.RandomElements[T](source, elements, n)
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
	return c.RandomElementIndices[T](source, elements, n)
}

// RandomIndex returns a random index from the given size. If the index has
// already been used, a new index is generated. If source is nil, a new source
// is created using the current time.
func RandomIndex(source *rand.Rand, size int, indicesUsed map[int]bool) int {
	return c.RandomIndex(source, size, indicesUsed)
}

// Shuffle returns a shuffled copy of the given slice. If source is nil, a new
// source is created using the current time.
func Shuffle[T any](source *rand.Rand, slice []T) []T {
	return c.Shuffle[T](source, slice)
}

// DefineLazy returns a Lazy[T] that will call the given function to
// calculate the value. The value is only calculated once, and the result
// is cached for subsequent calls.
func DefineLazy[T any](f func() T) c.Lazy[T] {
	return c.DefineLazy(f)
}

// Comparable is a type constraint for three types: int, int64, and string. By
// using this type constraint for a generic parameter, the parameter can be used
// as a map key and it can be sorted.
type Comparable = c.Comparable

// Keys returns a slice of all values in the given map.
func Keys[M ~map[K]V, K Comparable, V any](m M) []K {
	return c.Keys(m)
}

// Values returns a slice of all values in the given map.
func Values[M ~map[K]V, K Comparable, V any](m M) []V {
	return c.Values(m)
}

// RangeMap ranges over a map in deterministic order by first sorting the
// keys. It provides a function which will be called for each key/value pair.
// The function should return true to stop iteration.
func RangeMap[M ~map[K]V, K Comparable, V any](
	m M,
	f func(key K, value V) bool,
) {
	c.RangeMap(m, f)
}

// Reverse reverses the given slice in place and returns it.
func Reverse[T any](slice []T) []T {
	return c.Reverse(slice)
}
