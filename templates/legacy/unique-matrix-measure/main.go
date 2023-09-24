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

type stop struct {
	ID        string `json:"id"`
	Reference int    `json:"reference"`
}

// This struct describes the expected json input by the runner.
// Features not needed can simply be deleted or commented out, but make
// sure that the corresponding option in `solver` is also commented out.
// In case you would like to support a different input format you can
// change the struct as you see fit. You may need to change some code in
// `solver` to use the new structure.
type input struct {
	Stops        []stop           `json:"stops"`
	Vehicles     []string         `json:"vehicles"`
	Starts       []stop           `json:"starts"`
	Ends         []stop           `json:"ends"`
	UniquePoints []route.Position `json:"unique_points"`
	UniqueMatrix [][]float64      `json:"unique_matrix"`
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

	// Make sure we have a matrix of NxN where N is the number of unique points.
	checkInput(i)

	// Create a full matrix from the reduced unique matrix. The full matrix also
	// has start/end locations for each vehicle incorporated. In this example we
	// do not use any start/end locations.
	// We go over each stop in each row of the full matrix and lookup it's stop
	// references to access the cost from the reduced size matrix.
	matrixSize := len(i.Stops) + 2*len(i.Vehicles)
	fullMatrix := make([][]float64, matrixSize)

	// Make all stops, starts and ends accessible from one slice.
	joinedStops := joinStops(i, matrixSize)

	for stopIndex1 := range fullMatrix {
		// Create an row filled with 0.
		fullMatrix[stopIndex1] = make([]float64, matrixSize)

		// Get first reference index.
		refIndex1 := joinedStops[stopIndex1].Reference
		for stopIndex2 := range fullMatrix[stopIndex1] {
			// Get second reference index and look up costs.
			refIndex2 := joinedStops[stopIndex2].Reference
			if refIndex1 != -1 && refIndex2 != -1 {
				cost := i.UniqueMatrix[refIndex1][refIndex2]
				fullMatrix[stopIndex1][stopIndex2] = cost
			}
		}
	}

	// Create a by index measure from the full matrix to be used in the routing
	// engine and use it for each vehicle.
	byIndex := measure.Matrix(fullMatrix)
	byIndexMeasures := make([]measure.ByIndex, len(i.Vehicles))
	for v := range i.Vehicles {
		byIndexMeasures[v] = byIndex
	}

	// For the routing engine we now need to create a compatible stop data format.
	stops := make([]route.Stop, len(i.Stops))
	for stopIndex, s := range i.Stops {
		stops[stopIndex] = route.Stop{
			ID:       s.ID,
			Position: i.UniquePoints[s.Reference],
		}
	}

	starts := make([]route.Position, len(i.Starts))
	for stopIndex, s := range i.Starts {
		starts[stopIndex] = i.UniquePoints[s.Reference]
	}

	ends := make([]route.Position, len(i.Ends))
	for stopIndex, s := range i.Ends {
		ends[stopIndex] = i.UniquePoints[s.Reference]
	}

	// Define base router.
	router, err := route.NewRouter(
		stops,
		i.Vehicles,
		route.Starts(starts),
		route.Ends(ends),
		route.ValueFunctionMeasures(byIndexMeasures),
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

func checkInput(i input) {
	if len(i.UniqueMatrix) != len(i.UniquePoints) {
		panic("number of matrix rows must match unique points size")
	}
	for _, row := range i.UniqueMatrix {
		if len(row) != len(i.UniquePoints) {
			panic("matrix rows size must match unique points size")
		}
	}

	if len(i.Starts) > 0 && len(i.Starts) != len(i.Vehicles) {
		panic("if starts are given they must match the number of vehicles")
	}

	if len(i.Ends) > 0 && len(i.Ends) != len(i.Vehicles) {
		panic("if ends are given they must match the number of vehicles")
	}
}

func joinStops(i input, matrixSize int) []stop {
	joinedStops := make([]stop, matrixSize)
	// Set default reference to -1.
	for index := range joinedStops {
		joinedStops[index].Reference = -1
	}

	// Copy stops to slice.
	copy(joinedStops, i.Stops)

	// Fill joined stops alternating with Starts if applicable.
	if len(i.Starts) > 0 {
		for index := len(i.Stops); index < matrixSize; index += 2 {
			joinedStops[index] = i.Starts[(index-len(i.Stops))/2]
		}
	}

	// Fill joined stops alternating with Ends if applicable.
	if len(i.Ends) > 0 {
		for index := len(i.Stops) + 1; index < matrixSize; index += 2 {
			joinedStops[index] = i.Ends[(index-len(i.Stops)-1)/2]
		}
	}
	return joinedStops
}
