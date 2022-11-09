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
		"testdata/maryland.osm.pbf",
		1000,
		1<<30,
		carProfile,
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := byPointDistance.Cost(
		route.Point{-123.1041788, 43.9965908},
		route.Point{-123.0774532, 44.044393},
	)
	fmt.Println(int(cost))
	// Output:
	// 7447
}

func ExampleDurationByPoint() {
	carProfile := routingkit.Car()
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.2)
	byPointDuration, err := routingkit.DurationByPoint(
		"testdata/maryland.osm.pbf",
		1000,
		1<<30,
		carProfile,
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := byPointDuration.Cost(
		route.Point{-123.1041788, 43.9965908},
		route.Point{-123.0774532, 44.044393},
	)
	fmt.Println(int(cost))
	// Output:
	// 6874
}

func ExampleMatrix() {
	srcs := []route.Point{
		{-123.1041788, 43.9965908},
		{-123.1117056, 44.0568198},
	}
	dests := []route.Point{
		{-123.0774532, 44.044393},
		{-123.0399395, 44.0276874},
	}
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.3)
	byIndexDistance, err := routingkit.Matrix(
		"testdata/maryland.osm.pbf",
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
	// 8050
}

func ExampleDurationMatrix() {
	srcs := []route.Point{
		{-123.1041788, 43.9965908},
		{-123.1117056, 44.0568198},
	}
	dests := []route.Point{
		{-123.0774532, 44.044393},
		{-123.0399395, 44.0276874},
	}
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 1.2)
	byIndexDistance, err := routingkit.DurationMatrix(
		"testdata/maryland.osm.pbf",
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
	// 7431
}

// The following code example shows how to create your own vehicle profile and
// use it with routingkit. It customizes the vehicle speed and makes it depend
// on the tags present. In this example the speed is fixed to a single value
// (defined in `customVehicleSpeedMapper`). Furthermore, only ways are allowed
// to be used by the `customVehicle` which have the highway tag and its value
// is motorway (defined in `customVehicleTagMapFilter`). Please refer to the
// OpenStreetMaps [documentation on ways][osm-docs] to learn more about [tags
// and their values][osm-ways].
// [osm-docs]: https://wiki.openstreetmap.org/wiki/Way
// [osm-ways]: https://wiki.openstreetmap.org/wiki/Key:highway
func Example_customProfile() {
	fallbackMeasure := route.ScaleByPoint(route.HaversineByPoint(), 2.1)

	// Restricts ways to defined OSM way tags.
	filter := func(id int, tags map[string]string) bool {
		highway := tags["highway"]
		return highway == "motorway"
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
		"testdata/maryland.osm.pbf",
		1000,
		1<<30,
		p,
		fallbackMeasure,
	)
	if err != nil {
		panic(err)
	}
	cost := m.Cost(
		route.Point{-123.1041788, 43.9965908},
		route.Point{-123.0774532, 44.044393},
	)
	fmt.Println(int(cost))
	// Output:
	// 12030
}
