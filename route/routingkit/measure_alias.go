package routingkit

import (
	rk "github.com/nextmv-io/go-routingkit/routingkit"
	r "github.com/nextmv-io/sdk/measure/routingkit"
	"github.com/nextmv-io/sdk/route"
)

// >>> DistanceClient implementation

// DistanceClient represents a RoutingKit distance client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
type DistanceClient = r.DistanceClient

// NewDistanceClient returns a new RoutingKit client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func NewDistanceClient(
	mapFile string,
	profile rk.Profile,
) (DistanceClient, error) {
	return r.NewDistanceClient(mapFile, profile)
}

// >>> DurationClient implementation

// DurationClient represents a RoutingKit duration client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
type DurationClient = r.DurationClient

// NewDurationClient returns a new RoutingKit client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func NewDurationClient(
	mapFile string,
	profile rk.Profile,
) (DurationClient, error) {
	return r.NewDurationClient(mapFile, profile)
}

// These profile constructors are exported for convenience -
// to avoid having to import two packages both called routingkit

// Car constructs a car routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Car = r.Car

// Bike constructs a bike routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Bike = r.Bike

// Truck constructs a truck routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Truck = r.Truck

// Pedestrian constructs a pedestrian routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Pedestrian = r.Pedestrian

// DurationByPoint is a measure that uses routingkit to calculate car travel
// times between given points. It needs a .osm.pbf map file, a radius in which
// points can be snapped to the road network, a cache size in bytes (1 << 30 = 1
// GB), a profile and a measure that is used in case no travel time can be
// computed.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func DurationByPoint(
	mapFile string,
	radius float64,
	cacheSize int64,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByPoint, error) {
	return r.DurationByPoint(mapFile, radius, cacheSize, profile, m)
}

// ByPoint constructs a route.ByPoint that computes the road network distance
// connecting any two points found within the provided mapFile. It needs a
// radius in which points can be snapped to the road network, a cache size in
// bytes (1 << 30 = 1 GB), a profile and a measure that is used in case no
// travel time can be computed.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func ByPoint(
	mapFile string,
	radius float64,
	cacheSize int64,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByPoint, error) {
	return r.ByPoint(mapFile, radius, cacheSize, profile, m)
}

// Matrix uses the provided mapFile to construct a route.ByIndex that can find
// the road network distance between any point in srcs and any point in dests.
// In addition to the mapFile, srcs and dests it needs a radius in which points
// can be snapped to the road network, a profile and a measure that is used in
// case no distances can be computed.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func Matrix(
	mapFile string,
	radius float64,
	srcs []route.Point,
	dests []route.Point,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByIndex, error) {
	return r.Matrix(mapFile, radius, srcs, dests, profile, m)
}

// DurationMatrix uses the provided mapFile to construct a route.ByIndex that
// can find the road network durations between any point in srcs and any point
// in dests.
// In addition to the mapFile, srcs and dests it needs a radius in which points
// can be snapped to the road network, a profile and a measure that is used in
// case no travel durations can be computed.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func DurationMatrix(
	mapFile string,
	radius float64,
	srcs []route.Point,
	dests []route.Point,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByIndex, error) {
	return r.DurationMatrix(mapFile, radius, srcs, dests, profile, m)
}
