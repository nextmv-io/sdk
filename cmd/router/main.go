package main

import (
	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	run.Run(handler)
}

// Struct to read from JSON in.
type input struct {
	Stops      []route.Stop `json:"stops,omitempty"`
	Vehicles   []string     `json:"vehicles,omitempty"`
	Quantities []int        `json:"quantities,omitempty"`
	Capacities []int        `json:"capacities,omitempty"`
}

func handler(i input, opt store.Options) (store.Solver, error) {
	router, err := route.NewRouter(
		i.Stops,
		i.Vehicles,
		route.Capacity(i.Quantities, i.Capacities),
	)
	if err != nil {
		return nil, err
	}

	return router.Solver(opt)
}
