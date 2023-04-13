// package main holds the implementation of the cloud-routing template.
package main

import (
	"log"
	"time"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/encode"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	err := run.Run(solver,
		run.Encode[run.CLIRunnerConfig, input](
			GenericEncoder[store.Solution, store.Options](encode.JSON()),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// solver takes the input and solver options and constructs a routing solver.
// All route features/options depend on the input format. Depending on your
// goal you can add, delete or fix options or add more input validations. Please
// see the [route package
// documentation](https://pkg.go.dev/github.com/nextmv-io/sdk/route) for further
// information on the options available to you.
var solver = func(i input, opts store.Options) (store.Solver, error) {
	// In case you directly expose the solver to untrusted, external input,
	// it is advisable from a security point of view to add strong
	// input validations before passing the data to the solver.

	err := i.prepareInputData()
	if err != nil {
		return nil, err
	}

	routerInput := i.transform()

	timeMeasures, err := buildTravelTimeMeasures(
		routerInput.Velocities,
		i.Vehicles,
		i.Stops,
		routerInput.Starts,
		routerInput.Ends,
		i.durationGroupsDomains,
	)
	if err != nil {
		return nil, err
	}

	// Define base router.
	router, err := route.NewRouter(
		routerInput.Stops,
		routerInput.Vehicles,
		route.Starts(routerInput.Starts),
		route.Ends(routerInput.Ends),
		route.Shifts(routerInput.Shifts),
		route.Precedence(routerInput.Precedences),
		route.Velocities(routerInput.Velocities),
		route.Unassigned(routerInput.Penalties),
		route.ValueFunctionMeasures(timeMeasures),
		route.Capacity(
			routerInput.Quantities["default"],
			routerInput.Capacities["default"],
		),
		// route.Services(routerInput.ServiceTimes),
	)
	if err != nil {
		return nil, err
	}

	// You can also fix solver options like the expansion limit below.
	opts.Diagram.Expansion.Limit = 1
	options := i.applySolveOptions(opts)

	// If the duration limit is unset, we set it to 10s. You can configure
	// longer solver run times here. For local runs there is no time limitation.
	// If you want to make cloud runs for longer than 5 minutes, please contact:
	// support@nextmv.io
	if options.Limits.Duration == 0 {
		options.Limits.Duration = 10 * time.Second
	}

	return router.Solver(options)
}
