// package main holds the implementation of the routing template.
package main

import (
	"log"
	"time"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	_, err := run.Run(solver)
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
	Stops               []route.Stop         `json:"stops"`
	Vehicles            []string             `json:"vehicles"`
	InitializationCosts []float64            `json:"initialization_costs"`
	Starts              []route.Position     `json:"starts"`
	Ends                []route.Position     `json:"ends"`
	Quantities          []int                `json:"quantities"`
	Capacities          []int                `json:"capacities"`
	Precedences         []route.Job          `json:"precedences"`
	Windows             []route.Window       `json:"windows"`
	Shifts              []route.TimeWindow   `json:"shifts"`
	Penalties           []int                `json:"penalties"`
	Backlogs            []route.Backlog      `json:"backlogs"`
	VehicleAttributes   []route.Attributes   `json:"vehicle_attributes"`
	StopAttributes      []route.Attributes   `json:"stop_attributes"`
	Velocities          []float64            `json:"velocities"`
	Groups              [][]string           `json:"groups"`
	ServiceTimes        []route.Service      `json:"service_times"`
	AlternateStops      []route.Alternate    `json:"alternate_stops"`
	Limits              []route.Limit        `json:"limits"`
	DurationLimits      []float64            `json:"duration_limits"`
	DistanceLimits      []float64            `json:"distance_limits"`
	ServiceGroups       []route.ServiceGroup `json:"service_groups"`
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

	// Define base router.
	router, err := route.NewRouter(
		i.Stops,
		i.Vehicles,
		route.Threads(2),
		route.Velocities(i.Velocities),
		route.Starts(i.Starts),
		route.Ends(i.Ends),
		route.Capacity(i.Quantities, i.Capacities),
		route.Precedence(i.Precedences),
		route.Services(i.ServiceTimes),
		route.Shifts(i.Shifts),
		route.Windows(i.Windows),
		route.Unassigned(i.Penalties),
		route.InitializationCosts(i.InitializationCosts),
		route.Backlogs(i.Backlogs),
		route.LimitDurations(
			i.DurationLimits,
			true, /*ignoreTriangular*/
		),
		route.LimitDistances(
			i.DistanceLimits,
			true, /*ignoreTriangular*/
		),
		route.Attribute(i.VehicleAttributes, i.StopAttributes),
		route.Grouper(i.Groups),
		route.Alternates(i.AlternateStops),
		route.ServiceGroups(i.ServiceGroups),
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
