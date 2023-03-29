// package main holds the implementation of the cloud-routing template.
package main

import (
	"log"
	"time"

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
		route.Velocities(routerInput.Velocities),
		route.Unassigned(routerInput.Penalties),
		route.InitializationCosts(routerInput.InitializationCosts),
		route.ValueFunctionMeasures(timeMeasures),
		route.LimitDurations(
			routerInput.DurationLimits,
			true, /*ignoreTriangular*/
		),
	)
	if err != nil {
		return nil, err
	}

	if len(routerInput.Windows) > 0 {
		err = router.Options(route.Windows(routerInput.Windows))
		if err != nil {
			return nil, err
		}
	}
	if len(routerInput.Shifts) > 0 {
		err = router.Options(route.Shifts(routerInput.Shifts))
		if err != nil {
			return nil, err
		}
	}

	if len(routerInput.StopAttributes) > 0 &&
		len(routerInput.VehicleAttributes) > 0 {
		err = router.Options(
			route.Attribute(
				routerInput.VehicleAttributes,
				routerInput.StopAttributes,
			),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(routerInput.ServiceTimes) > 0 {
		err = router.Options(route.Services(routerInput.ServiceTimes))
		if err != nil {
			return nil, err
		}
	}
	for kind := range routerInput.Kinds {
		err = router.Options(
			route.Capacity(
				routerInput.Quantities[kind],
				routerInput.Capacities[kind],
			),
		)
		if err != nil {
			return nil, err
		}
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

func (i *input) prepareInputData() error {
	i.Stops = append(i.Stops, i.AlternateStops...)
	i.applyVehicleDefaults()
	i.applyStopDefaults()
	// Handle dynamic fields
	err := i.handleDynamics()
	if err != nil {
		return err
	}
	i.autoConfigureUnassigned()
	err = i.makeDurationGroups()
	if err != nil {
		return err
	}
	return nil
}
