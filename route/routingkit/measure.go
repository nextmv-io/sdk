package routingkit

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/dgraph-io/ristretto"
	"github.com/nextmv-io/go-routingkit/routingkit"
	"github.com/nextmv-io/sdk/route"
)

const cacheItemCost int64 = 80

// These profile constructors are exported for convenience -
// to avoid having to import two packages both called routingkit

// Car constructs a car routingkit profile.
var Car = routingkit.Car

// Bike constructs a bike routingkit profile.
var Bike = routingkit.Bike

// Truck constructs a truck routingkit profile.
var Truck = routingkit.Truck

// Pedestrian constructs a pedestrian routingkit profile.
var Pedestrian = routingkit.Pedestrian

// DurationByPoint is a measure that uses routingkit to calculate car travel
// times between given points. It needs a .osm.pbf map file, a radius in which
// points can be snapped to the road network, a cache size in bytes (1 << 30 = 1
// GB), a profile and a measure that is used in case no travel time can be
// computed.
func DurationByPoint(
	mapFile string,
	radius float64,
	cacheSize int64,
	profile routingkit.Profile,
	m route.ByPoint,
) (route.ByPoint, error) {
	client, err := routingkit.NewTravelTimeClient(mapFile, profile)
	if err != nil {
		return nil, err
	}
	cache, err := ristretto.NewCache(&ristretto.Config{
		// NumCounters should be 10 times the max number of cached items. Since
		// the cost of each item is cacheItemCost , 10 * cacheSize /
		// cacheItemCost gives the correct value
		NumCounters: 10 * cacheSize / cacheItemCost,
		MaxCost:     cacheSize,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("creating cache: %v", err)
	}
	return durationByPoint{
		client:    client,
		mapFile:   mapFile,
		radius:    radius,
		m:         m,
		cache:     cache,
		cacheSize: cacheSize,
		profile:   profile,
	}, nil
}

type durationByPoint struct {
	client    routingkit.TravelTimeClient
	m         route.ByPoint
	cache     *ristretto.Cache
	mapFile   string
	radius    float64
	cacheSize int64
	profile   routingkit.Profile
}

// Cost calculates the road network travel time between the points.
func (b durationByPoint) Cost(p1, p2 route.Point) float64 {
	key := make([]byte, 32)
	binary.LittleEndian.PutUint64(key[0:], math.Float64bits(p1[0]))
	binary.LittleEndian.PutUint64(key[8:], math.Float64bits(p1[1]))
	binary.LittleEndian.PutUint64(key[16:], math.Float64bits(p2[0]))
	binary.LittleEndian.PutUint64(key[24:], math.Float64bits(p2[1]))
	val, found := b.cache.Get(key)
	if found {
		return val.(float64)
	}

	d := b.client.TravelTime(coords(p1), coords(p2))
	if b.m != nil && d == routingkit.MaxDistance {
		c := b.m.Cost(p1, p2)
		b.cache.Set(key, c, cacheItemCost)
		return c
	}
	dInSeconds := float64(d) / 1000.0
	// the cost of an entry is cacheItemCost
	b.cache.Set(key, dInSeconds, cacheItemCost)
	return dInSeconds
}

// Triangular indicates that the measure does have the triangularity property.
func (b durationByPoint) Triangular() bool {
	return true
}

// MarshalJSON serializes the route.
func (b durationByPoint) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["type"] = "routingkitDuration"
	m["osm"] = b.mapFile
	m["radius"] = b.radius
	m["cache_size"] = b.cacheSize
	m["profile"] = ProfileLoader{&b.profile}
	if b.m != nil {
		m["measure"] = b.m
	}
	return json.Marshal(m)
}

// ByPoint constructs a route.ByPoint that computes the road network distance
// connecting any two points found within the provided mapFile. It needs a
// radius in which points can be snapped to the road network, a cache size in
// bytes (1 << 30 = 1 GB), a profile and a measure that is used in case no
// travel time can be computed.
func ByPoint(
	mapFile string,
	radius float64,
	cacheSize int64,
	profile routingkit.Profile,
	m route.ByPoint,
) (route.ByPoint, error) {
	client, err := routingkit.NewDistanceClient(mapFile, profile)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	cache, err := ristretto.NewCache(&ristretto.Config{
		// NumCounters should be 10 times the max number of cached items. Since
		// the cost of each item is cacheItemCost , 10 * cacheSize /
		// cacheItemCost gives the correct value
		NumCounters: 10 * cacheSize / cacheItemCost,
		MaxCost:     cacheSize,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("creating cache: %v", err)
	}
	return byPoint{
		client:    client,
		mapFile:   mapFile,
		radius:    radius,
		m:         m,
		cache:     cache,
		cacheSize: cacheSize,
		profile:   profile,
	}, nil
}

type byPoint struct {
	client    routingkit.DistanceClient
	m         route.ByPoint
	cache     *ristretto.Cache
	mapFile   string
	radius    float64
	cacheSize int64
	profile   routingkit.Profile
}

// Cost calculates the road network distance between the points.
func (b byPoint) Cost(p1, p2 route.Point) float64 {
	key := make([]byte, 32)
	binary.LittleEndian.PutUint64(key[0:], math.Float64bits(p1[0]))
	binary.LittleEndian.PutUint64(key[8:], math.Float64bits(p1[1]))
	binary.LittleEndian.PutUint64(key[16:], math.Float64bits(p2[0]))
	binary.LittleEndian.PutUint64(key[24:], math.Float64bits(p2[1]))
	val, found := b.cache.Get(key)
	if found {
		return val.(float64)
	}

	d := b.client.Distance(coords(p1), coords(p2))
	if b.m != nil && d == routingkit.MaxDistance {
		c := b.m.Cost(p1, p2)
		b.cache.Set(key, c, cacheItemCost)
		return c
	}
	// the cost of an entry is cacheItemCost
	b.cache.Set(key, float64(d), cacheItemCost)
	return float64(d)
}

// Triangular indicates that the measure does have the triangularity property.
func (b byPoint) Triangular() bool {
	return true
}

// MarshalJSON serializes the route.
func (b byPoint) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["type"] = "routingkit"
	m["osm"] = b.mapFile
	m["radius"] = b.radius
	m["profile"] = ProfileLoader{&b.profile}
	m["cache_size"] = b.cacheSize
	if b.m != nil {
		m["measure"] = b.m
	}
	return json.Marshal(m)
}

// Matrix uses the provided mapFile to construct a route.ByIndex that can find
// the road network distance between any point in srcs and any point in dests.
// In addition to the mapFile, srcs and dests it needs a radius in which points
// can be snapped to the road network, a profile and a measure that is used in
// case no distances can be computed.
func Matrix(
	mapFile string,
	radius float64,
	srcs []route.Point,
	dests []route.Point,
	profile routingkit.Profile,
	m route.ByPoint,
) (route.ByIndex, error) {
	client, err := routingkit.NewDistanceClient(mapFile, profile)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	mx := client.Matrix(coordsSlice(srcs), coordsSlice(dests))

	return matrix{
		ByIndex:    route.Matrix(float64Matrix(mx, srcs, dests, m, false)),
		mapFile:    mapFile,
		radius:     radius,
		srcs:       srcs,
		dests:      dests,
		profile:    &ProfileLoader{&profile},
		ByPoint:    m,
		clientType: "routingkitMatrix",
	}, nil
}

type matrix struct {
	route.ByPoint
	route.ByIndex
	mapFile    string
	clientType string
	srcs       []route.Point
	dests      []route.Point
	radius     float64
	profile    *ProfileLoader
}

// Cost returns the road network distance between the points.
func (m matrix) Cost(i, j int) float64 {
	return m.ByIndex.Cost(i, j)
}

// MarshalJSON serializes the route.
func (m matrix) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"type":         m.clientType,
		"osm":          m.mapFile,
		"radius":       m.radius,
		"sources":      m.srcs,
		"destinations": m.dests,
		"profile":      m.profile,
	}
	if m.ByPoint != nil {
		data["measure"] = m.ByPoint
	}

	return json.Marshal(data)
}

// DurationMatrix uses the provided mapFile to construct a route.ByIndex that
// can find the road network durations between any point in srcs and any point
// in dests.
// In addition to the mapFile, srcs and dests it needs a radius in which points
// can be snapped to the road network, a profile and a measure that is used in
// case no travel durations can be computed.
func DurationMatrix(
	mapFile string,
	radius float64,
	srcs []route.Point,
	dests []route.Point,
	profile routingkit.Profile,
	m route.ByPoint,
) (route.ByIndex, error) {
	client, err := routingkit.NewTravelTimeClient(mapFile, profile)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	mx := client.Matrix(coordsSlice(srcs), coordsSlice(dests))

	return matrix{
		ByIndex:    route.Matrix(float64Matrix(mx, srcs, dests, m, true)),
		mapFile:    mapFile,
		radius:     radius,
		srcs:       srcs,
		dests:      dests,
		ByPoint:    m,
		profile:    &ProfileLoader{&profile},
		clientType: "routingkitDurationMatrix",
	}, nil
}

func coords(p route.Point) []float32 {
	return []float32{float32(p[0]), float32(p[1])}
}

func coordsSlice(ps []route.Point) [][]float32 {
	cs := make([][]float32, len(ps))
	for i, p := range ps {
		cs[i] = coords(p)
	}
	return cs
}

func float64Matrix(m [][]uint32,
	srcs []route.Point,
	dests []route.Point,
	fallback route.ByPoint,
	duration bool,
) [][]float64 {
	fM := make([][]float64, len(m))
	for i, r := range m {
		fM[i] = make([]float64, len(r))
		for j, c := range r {
			if fallback != nil && c == routingkit.MaxDistance {
				fM[i][j] = fallback.Cost(srcs[i], dests[j])
			} else {
				if duration {
					fM[i][j] = float64(c) / 1000.0 // convert to seconds
				} else {
					fM[i][j] = float64(c)
				}
			}
		}
	}
	return fM
}
