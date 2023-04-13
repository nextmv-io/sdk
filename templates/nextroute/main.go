// © 2019-2021 nextmv.io inc. All rights reserved.
// nextmv.io, inc. CONFIDENTIAL
//
// This file includes unpublished proprietary source code of nextmv.io, inc.
// The copyright notice above does not evidence any actual or intended
// publication of such source code. Disclosure of this source code or any
// related proprietary information is strictly prohibited without the express
// written permission of nextmv.io, inc.

package main

import (
	"context"
	"log"

	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/nextroute"
)

func main() {
	runner := run.CLI(solver)
	err := runner.Run(context.Background())
	if err != nil {
		log.Fatalf("could not run solver: %v", err)
	}
}

func solver(i nextroute.Input, opts nextroute.EngineOptions) ([]nextroute.Solution, error) {
	model, err := nextroute.NewModel(i, opts.ModelOptions)
	if err != nil {
		return nil, err
	}

	// Create a solver.
	solver, err := nextroute.NewSolver(model, opts.SolverOptions)
	if err != nil {
		return nil, err
	}

	// Solve the problem.
	solutions, err := solver.Solve(opts.SolveOptions)
	if err != nil {
		return nil, err
	}

	return solutions, nil
}
