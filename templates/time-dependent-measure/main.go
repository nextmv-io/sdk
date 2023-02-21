// package main holds the implementation of the routing template.
package main

import (
	"log"
	"time"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	err := run.Run(solver)
	if err != nil {
		log.Fatal(err)
	}
}

// This struct describes the expected json input by the runner.
// Features not needed can simply be deleted or commented out, but make
// sure that the corresponding option in `solver` is also commented out.
// In case you would like to support a different input format you can
// change the struct as you see fit. You may need to change some code in
// `solver` to use the new structure.
type input struct {
	Stops    []route.Stop       `json:"stops"`
	Vehicles []string           `json:"vehicles"`
	Shifts   []route.TimeWindow `json:"shifts"`
	Starts   []route.Position   `json:"starts"`
	Ends     []route.Position   `json:"ends"`
}

// solver takes the input and solver options and constructs a routing solver.
// All route features/options depend on the input format. Depending on your
// goal you can add, delete or fix options or add more input validations. Please
// see the [route package
// documentation](https://pkg.go.dev/github.com/nextmv-io/sdk/route) for further
// information on the options available to you.
func solver(i input, opts store.Options) (store.Solver, error) {
	// In case you directly expose the solver to untrusted, external input,
	// it is advisable from a security point of view to add strong
	// input validations before passing the data to the solver.

	// Create two measures with different costs.
	points := make([]measure.Point, len(i.Stops)+2*len(i.Vehicles))
	for i, s := range i.Stops {
		points[i] = measure.Point{
			s.Position.Lon, s.Position.Lat,
		}
	}
	m1 := measure.Indexed(measure.HaversineByPoint(), points)
	m1 = measure.Override(m1, measure.Constant(0), func(from, to int) bool {
		return from > len(points)-1 || to > len(points)-1
	})
	m2 := measure.Scale(m1, 2)

	// Create a byIndexAndTime to combine each measure with an end time for up
	// until (exclusive) to use it.
	byIndexAndTime := make([][]route.ByIndexAndTime, len(i.Shifts))
	for shiftIndex, t := range i.Shifts {
		byIndexAndTime[shiftIndex] = make([]route.ByIndexAndTime, 2)
		byIndexAndTime[shiftIndex][0] = route.ByIndexAndTime{
			Measure: m1,
			EndTime: int(t.Start.Add(30 * time.Second).Unix()),
		}
		byIndexAndTime[shiftIndex][1] = route.ByIndexAndTime{
			Measure: m2,
			EndTime: int(t.Start.Add(60 * time.Second).Unix()),
		}
	}

	// Use the time dependent measures for each vehicle to adhere to the
	// ValueFunctionMeasures option.
	dependentMeasures := make([]measure.DependentByIndex, len(i.Vehicles))
	for index := range i.Vehicles {
		// Create a time dependent measure client.
		dependentMeasure, err := route.NewTimeDependentMeasure(
			int(i.Shifts[index].Start.Unix()),
			byIndexAndTime[index],
			m1,
		)
		if err != nil {
			panic("could not create dependent measure client")
		}
		dependentMeasures[index] = dependentMeasure
	}

	// Define base router.
	router, err := route.NewRouter(
		i.Stops,
		i.Vehicles,
		route.Shifts(i.Shifts),
		route.ValueFunctionDependentMeasures(dependentMeasures),
	)
	if err != nil {
		return nil, err
	}

	// You can also fix solver options like the expansion limit below.
	opts.Diagram.Expansion.Limit = 1
	// A duration limit of 0 is treated as infinity. For cloud runs you need to
	// set an explicit duration limit which is why it is currently set to 10s
	// here in case no duration limit is set. For local runs there is no time
	// limitation. If you want to make cloud runs for longer than 5 minutes,
	// please contact: support@nextmv.io
	if opts.Limits.Duration == 0 {
		opts.Limits.Duration = 10 * time.Second
	}

	return router.Solver(opts)
}
