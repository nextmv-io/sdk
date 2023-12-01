package routingkit

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/dgraph-io/ristretto"
	rk "github.com/nextmv-io/go-routingkit/routingkit"
	"github.com/nextmv-io/sdk/route"
	"github.com/twpayne/go-polyline"
)

// >>> DistanceClient implementation

// DistanceClient represents a RoutingKit distance client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
type DistanceClient interface {
	// Measure returns a route.ByPoint that can calculate the road network
	// distance between any two points found within the provided mapFile.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
	Measure(radius float64, cacheSize int64, fallback route.ByPoint) (route.ByPoint, error)
	// Matrix returns a route.ByIndex that represents the road network distance
	// matrix as a route.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
	Matrix(srcs []route.Point, dests []route.Point) (route.ByIndex, error)
	// Polyline requests polylines for the given points. The first parameter
	// returns a polyline from start to end and the second parameter returns a
	// list of polylines, one per leg.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
	Polyline(points []route.Point) (string, []string, error)
}

// NewDistanceClient returns a new RoutingKit client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func NewDistanceClient(
	mapFile string,
	profile rk.Profile,
) (DistanceClient, error) {
	client, err := rk.NewDistanceClient(mapFile, profile)
	if err != nil {
		return nil, err
	}
	return distanceClient{
		mapFile: mapFile,
		profile: profile,
		client:  client,
	}, nil
}

type distanceClient struct {
	mapFile string
	profile rk.Profile
	client  rk.DistanceClient
}

// Measure returns a route.ByPoint that can calculate the road network distance
// between any two points found within the provided mapFile.
func (c distanceClient) Measure(radius float64, cacheSize int64, fallback route.ByPoint) (route.ByPoint, error) {
	return newByPoint(c.client, c.mapFile, radius, cacheSize, c.profile, fallback)
}

// Matrix returns a route.ByIndex that represents the road network distance
// matrix as a route.
func (c distanceClient) Matrix(srcs []route.Point, dests []route.Point) (route.ByIndex, error) {
	return newMatrix(c.client, c.mapFile, 0, srcs, dests, c.profile, nil)
}

// Polyline requests polylines for the given points. The first parameter
// returns a polyline from start to end and the second parameter returns a list
// of polylines, one per leg.
func (c distanceClient) Polyline(points []route.Point) (string, []string, error) {
	return getPolyLine(points, c.client.Route)
}

// >>> DurationClient implementation

// DurationClient represents a RoutingKit duration client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
type DurationClient interface {
	// Measure returns a route.ByPoint that can calculate the road network
	// travel time between any two points found within the provided mapFile.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
	Measure(radius float64, cacheSize int64, fallback route.ByPoint) (route.ByPoint, error)
	// Matrix returns a route.ByIndex that represents the road network travel
	// time matrix as a route.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
	Matrix(srcs []route.Point, dests []route.Point) (route.ByIndex, error)
	// Polyline requests polylines for the given points. The first parameter
	// returns a polyline from start to end and the second parameter returns a
	// list of polylines, one per leg.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
	Polyline(points []route.Point) (string, []string, error)
}

// NewDurationClient returns a new RoutingKit client.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
func NewDurationClient(
	mapFile string,
	profile rk.Profile,
) (DurationClient, error) {
	client, err := rk.NewTravelTimeClient(mapFile, profile)
	if err != nil {
		return nil, err
	}
	return durationClient{
		mapFile: mapFile,
		profile: profile,
		client:  client,
	}, nil
}

type durationClient struct {
	mapFile string
	profile rk.Profile
	client  rk.TravelTimeClient
}

// Measure returns a route.ByPoint that can calculate the road network travel
// time between any two points found within the provided mapFile.
func (c durationClient) Measure(radius float64, cacheSize int64, fallback route.ByPoint) (route.ByPoint, error) {
	return newDurationByPoint(c.client, c.mapFile, radius, cacheSize, c.profile, fallback)
}

// Matrix returns a route.ByIndex that represents the road network travel time
// matrix as a route.
func (c durationClient) Matrix(srcs []route.Point, dests []route.Point) (route.ByIndex, error) {
	return newDurationMatrix(c.client, c.mapFile, 0, srcs, dests, c.profile, nil)
}

// Polyline requests polylines for the given points. The first parameter
// returns a polyline from start to end and the second parameter returns a list
// of polylines, one per leg.
func (c durationClient) Polyline(points []route.Point) (string, []string, error) {
	return getPolyLine(points, c.client.Route)
}

const cacheItemCost int64 = 80

// These profile constructors are exported for convenience -
// to avoid having to import two packages both called routingkit

// Car constructs a car routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Car = rk.Car

// Bike constructs a bike routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Bike = rk.Bike

// Truck constructs a truck routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Truck = rk.Truck

// Pedestrian constructs a pedestrian routingkit profile.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
var Pedestrian = rk.Pedestrian

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
	client, err := rk.NewTravelTimeClient(mapFile, profile)
	if err != nil {
		return nil, err
	}
	return newDurationByPoint(client, mapFile, radius, cacheSize, profile, m)
}

func newDurationByPoint(
	client rk.TravelTimeClient,
	mapFile string,
	radius float64,
	cacheSize int64,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByPoint, error) {
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
	client    rk.TravelTimeClient
	m         route.ByPoint
	cache     *ristretto.Cache
	mapFile   string
	radius    float64
	cacheSize int64
	profile   rk.Profile
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
	if b.m != nil && d == rk.MaxDistance {
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
	client, err := rk.NewDistanceClient(mapFile, profile)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return newByPoint(client, mapFile, radius, cacheSize, profile, m)
}

func newByPoint(
	client rk.DistanceClient,
	mapFile string,
	radius float64,
	cacheSize int64,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByPoint, error) {
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
	client    rk.DistanceClient
	m         route.ByPoint
	cache     *ristretto.Cache
	mapFile   string
	radius    float64
	cacheSize int64
	profile   rk.Profile
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
	if b.m != nil && d == rk.MaxDistance {
		c := b.m.Cost(p1, p2)
		b.cache.Set(key, c, cacheItemCost)
		return c
	}
	// the cost of an entry is cacheItemCost
	b.cache.Set(key, float64(d), cacheItemCost)
	return float64(d)
}

// Creates polylines for the given points. First return parameter is a polyline
// from start to end, second parameter is a list of polylines per leg in the
// route.
func (b byPoint) Polyline(points []route.Point) (string, []string, error) {
	return getPolyLine(points, b.client.Route)
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
	client, err := rk.NewDistanceClient(mapFile, profile)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return newMatrix(client, mapFile, radius, srcs, dests, profile, m)
}

func newMatrix(
	client rk.DistanceClient,
	mapFile string,
	radius float64,
	srcs []route.Point,
	dests []route.Point,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByIndex, error) {
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
	client, err := rk.NewTravelTimeClient(mapFile, profile)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return newDurationMatrix(client, mapFile, radius, srcs, dests, profile, m)
}

func newDurationMatrix(
	client rk.TravelTimeClient,
	mapFile string,
	radius float64,
	srcs []route.Point,
	dests []route.Point,
	profile rk.Profile,
	m route.ByPoint,
) (route.ByIndex, error) {
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
			if fallback != nil && c == rk.MaxDistance {
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

// getPolyLine requests the polylines for the given route from the routingkit
// client. It returns the complete polyline and a list of polylines per leg.
func getPolyLine(
	points []route.Point,
	router func(from []float32, to []float32) (uint32, [][]float32),
) (string, []string, error) {
	encodedPolylines := make([]string, len(points)-1)
	completePolyline := [][]float64{}
	for i := 0; i < len(points)-1; i++ {
		p1 := points[i]
		p2 := points[i+1]
		dist, poly32 := router(coords(p1), coords(p2))
		poly64 := make([][]float64, len(poly32))
		for i, p := range poly32 {
			poly64[i] = []float64{float64(p[0]), float64(p[1])}
		}
		encodedPolylines[i] = string(polyline.EncodeCoords(poly64))
		if dist == rk.MaxDistance {
			return "", []string{}, fmt.Errorf("no route found between %v and %v", p1, p2)
		}
		completePolyline = append(completePolyline, poly64...)
	}
	return string(polyline.EncodeCoords(completePolyline)), encodedPolylines, nil
}
