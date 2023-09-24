package mip

import (
	"github.com/nextmv-io/sdk/run/schema"
	"github.com/nextmv-io/sdk/run/statistics"
)

// CustomResultStatistics is an example of custom statistics that can be added
// to the output and used in experiments.
type CustomResultStatistics struct {
	// Constraints in the matrix, i.e. the number of constraints.
	Constraints int `json:"constraints,omitempty"`
	// Provider of the solution.
	Provider SolverProvider `json:"provider,omitempty"`
	// Status of the solution.
	Status string `json:"status,omitempty"`
	// Variables in the matrix, i.e. the number of variables.
	Variables int `json:"variables,omitempty"`
}

// Format the MIP solution into the output format that the runner expects. The
// solution represents the decisions that were made by the solver, translated
// into the domain of the problem.
func Format(
	options any,
	solution any,
	solverSolution Solution,
) schema.Output {
	output := schema.NewOutput(options, []any{})
	output.Statistics = statistics.NewStatistics()
	output.Statistics.Run = &statistics.Run{}
	output.Statistics.Result = &statistics.Result{}
	if solverSolution == nil || !solverSolution.HasValues() || solverSolution.IsInfeasible() {
		return output
	}

	duration := solverSolution.RunTime().Seconds()
	output.Statistics.Result.Duration = &duration
	output.Statistics.Run.Duration = &duration

	value := statistics.Float64(solverSolution.ObjectiveValue())
	output.Statistics.Result.Value = &value

	output.Solutions = []any{solution}

	return output
}

// DefaultCustomResultStatistics creates default custom statistics for a given
// solution.
func DefaultCustomResultStatistics(model Model, solution Solution) CustomResultStatistics {
	status := "unknown"
	switch {
	case solution.IsOptimal():
		status = "optimal"
	case solution.IsUnbounded():
		status = "unbounded"
	case solution.IsSubOptimal():
		status = "suboptimal"
	case solution.IsInfeasible():
		status = "infeasible"
	}

	return CustomResultStatistics{
		Status:      status,
		Variables:   len(model.Vars()),
		Constraints: len(model.Constraints()),
		Provider:    solution.Provider(),
	}
}
