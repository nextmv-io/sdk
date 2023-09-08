// package main holds the implementation of the routing template.
package main

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/nextmv-io/sdk/cluster/kmeans"
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
	Stops             []route.Stop `json:"stops"`
	Vehicles          []string     `json:"vehicles"`
	StopWeight        []int        `json:"stop_weight"`
	ClusterCapacities []int        `json:"cluster_capacities"`
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

	// Creates an evenly sized cluster for every vehicles and creates
	// compatibility attributes for each stop/vehicle such that every cluster
	// must be served by 1 vehicle.
	vehicleAttributes, stopAttributes, err := cluster(i)
	if err != nil {
		return nil, err
	}

	// Define base router.
	router, err := route.NewRouter(
		i.Stops,
		i.Vehicles,
		route.Attribute(vehicleAttributes, stopAttributes),
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

// Creates a set of big clusters where 1 vehicles serves exactly 1 cluster.
func cluster(input input) ([]route.Attributes, []route.Attributes, error) {
	solution, err := clusterSolution(input)
	if err != nil {
		return nil, nil, err
	}

	vehicleAttributes := make([]route.Attributes, 0)
	stopAttributes := make([]route.Attributes, 0)
	for i, v := range input.Vehicles {
		vehicleAttributes = append(vehicleAttributes, route.Attributes{
			ID:         v,
			Attributes: []string{strconv.Itoa(i)},
		})
	}

	for clusterIndex, c := range solution.Clusters() {
		for _, stopIndex := range c.Indices() {
			attr := strconv.Itoa(clusterIndex)
			stopAttributes = append(
				stopAttributes,
				route.Attributes{
					ID:         input.Stops[stopIndex].ID,
					Attributes: []string{attr},
				})
		}
	}
	return vehicleAttributes, stopAttributes, nil
}

func clusterSolution(input input) (kmeans.Solution, error) {
	maximumPoints := make([]int, len(input.Vehicles))
	maximumValues := make([]int, len(input.Vehicles))
	values := make([][]int, len(input.Vehicles))
	points := make([]measure.Point, len(input.Stops))
	weights := make([]int, len(input.Stops))

	for i, w := range input.Stops {
		weights[i] = input.StopWeight[i]
		points[i] = measure.Point{w.Position.Lat, w.Position.Lon}
	}

	for idx := 0; idx < len(input.Vehicles); idx++ {
		maximumPoints[idx] = int(
			math.Ceil(float64(len(input.Stops)) / float64(len(input.Vehicles))),
		)
		maximumValues[idx] = input.ClusterCapacities[idx]
		values[idx] = weights
	}

	// We create a kmeans model using options.
	model, err := kmeans.NewModel(
		points,
		len(input.Vehicles),
		kmeans.MaximumPoints(maximumPoints),
		kmeans.MaximumSumValue(maximumValues, values),
	)
	if err != nil {
		return nil, err
	}

	// We create a solver with the model.
	solver, err := kmeans.NewSolver(model)
	if err != nil {
		return nil, err
	}
	solveOptions := kmeans.NewSolveOptions().
		SetMeasure(measure.EuclideanByPoint()).
		SetMaximumDuration(10 * time.Second)

	solution, err := solver.Solve(solveOptions)
	if err != nil {
		panic(err)
	}
	return solution, nil
}
