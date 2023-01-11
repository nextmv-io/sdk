// package main holds the implementation of the routing template.
package main

import (
	"time"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	run.Run(solver)
}

// This struct describes the expected json input by the runner.
// Features not needed can simply be deleted or commented out, but make
// sure that the corresponding option in `solver` is also commented out.
// In case you would like to support a different input format you can
// change the struct as you see fit. You may need to change some code in
// `solver` to use the new structure.
type input struct {
	Stops          []route.Stop       `json:"stops"`
	Vehicles       []string           `json:"vehicles"`
	Starts         []route.Position   `json:"starts"`
	Ends           []route.Position   `json:"ends"`
	Shifts         []route.TimeWindow `json:"shifts"`
	DurationMatrix [][]float64        `json:"duration_matrix"`
	DistanceMatrix [][]float64        `json:"distance_matrix"`
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

	distanceMatrices := make([]route.ByIndex, len(i.Vehicles))
	durationMatrices := make([]route.ByIndex, len(i.Vehicles))
	for j := range i.Vehicles {
		distanceMatrices[j] = route.Matrix(i.DistanceMatrix)
		durationMatrices[j] = route.Matrix(i.DurationMatrix)
	}

	// Define base router.
	router, err := route.NewRouter(
		i.Stops,
		i.Vehicles,
		route.ValueFunctionMeasures(distanceMatrices),
		route.TravelTimeMeasures(durationMatrices),
		route.Shifts(i.Shifts),
	)
	if err != nil {
		return nil, err
	}

	if len(i.Starts) > 0 {
		err = router.Options(route.Starts(i.Starts))
		if err != nil {
			panic("error using starts option")
		}
	}
	if len(i.Ends) > 0 {
		err = router.Options(route.Ends(i.Ends))
		if err != nil {
			panic("error using ends option")
		}
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
