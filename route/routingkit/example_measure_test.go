package routingkit_test

import (
	"fmt"

	goroutingkit "github.com/nextmv-io/go-routingkit/routingkit"
	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/route/routingkit"
)

func ExampleByPoint() {
	carProfile := routingkit.Car()
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.3)
	byPointDistance, err := routingkit.ByPoint(
		"testdata/rk_test.osm.pbf",
		1000,
		1<<30,
		carProfile,
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := byPointDistance.Cost(
		route.Point{7.33745, 52.14758},
		route.Point{7.34979, 52.15149},
	)
	fmt.Println(int(cost))
	// Output:
	// 1225
}

func ExampleDurationByPoint() {
	carProfile := routingkit.Car()
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.2)
	byPointDuration, err := routingkit.DurationByPoint(
		"testdata/rk_test.osm.pbf",
		1000,
		1<<30,
		carProfile,
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := byPointDuration.Cost(
		route.Point{7.33745, 52.14758},
		route.Point{7.34979, 52.15149},
	)
	fmt.Println(int(cost))
	// Output:
	// 187
}

func ExampleMatrix() {
	srcs := []route.Point{
		{7.33745, 52.14758},
		{7.32486, 52.14280},
	}
	dests := []route.Point{
		{7.34979, 52.15149},
		{7.33293, 52.13893},
	}
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.3)
	byIndexDistance, err := routingkit.Matrix(
		"testdata/rk_test.osm.pbf",
		1000,
		srcs,
		dests,
		routingkit.Car(),
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := byIndexDistance.Cost(0, 1)
	fmt.Println(int(cost))
	// Output:
	// 1219
}

func ExampleDurationMatrix() {
	srcs := []route.Point{
		{7.33745, 52.14758},
		{7.32486, 52.14280},
	}
	dests := []route.Point{
		{7.34979, 52.15149},
		{7.33293, 52.13893},
	}
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.2)
	byIndexDistance, err := routingkit.DurationMatrix(
		"testdata/rk_test.osm.pbf",
		1000,
		srcs,
		dests,
		routingkit.Car(),
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := byIndexDistance.Cost(0, 1)
	fmt.Println(int(cost))
	// Output:
	// 215
}

// The following code example shows how to create your own vehicle profile and
// use it with routingkit. It customizes the vehicle speed and makes it depend
// on the tags present. In this example the speed is fixed to a single value
// (defined in `customVehicleSpeedMapper`). Furthermore, only ways are allowed
// to be used by the `customVehicle` which have the highway tag and its value
// is not motorway (defined in `customVehicleTagMapFilter`). Please refer to the
// OpenStreetMaps [documentation on ways][osm-docs] to learn more about [tags
// and their values][osm-ways].
// [osm-docs]: https://wiki.openstreetmap.org/wiki/Way
// [osm-ways]: https://wiki.openstreetmap.org/wiki/Key:highway
func Example_customProfile() {
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.3)

	// Restricts ways to defined OSM way tags.
	filter := func(id int, tags map[string]string) bool {
		// Uses the default filter to filter out ways which are not routable by
		// a car.
		if !routingkit.Car().Filter(id, tags) {
			return false
		}
		// Additionally filter out motorway and trunk (only use small
		// streets/roads).
		// Returning true here, should give different costs for the points used
		// below. The shortest path uses a 'trunk' way.
		return tags["highway"] != "motorway" && tags["highway"] != "trunk"
	}

	// Defines a speed per OSM way tag.
	speedMapper := func(_ int, tags map[string]string) int {
		return 10
	}

	// Defines the custom profile.
	p := goroutingkit.NewProfile(
		"customVehicle", // Profile name
		// TransportationMode, other values: BikeMode, PedestrianMode.
		goroutingkit.VehicleMode,
		// Prevent left turns, only for TransportationMode "VehicleMode".
		false,
		// Prevent U-turns, only for TransportationMode "VehicleMode".
		false,
		filter,
		speedMapper,
	)

	// Defines a routingkit measure using the custom profile.
	m, err := routingkit.DurationByPoint(
		"testdata/rk_test.osm.pbf",
		1000,
		1<<30,
		p,
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := m.Cost(

		route.Point{7.34375, 52.16391},
		route.Point{7.32165, 52.15834},
	)
	fmt.Println(int(cost))
	// Output:
	// 1194
}
