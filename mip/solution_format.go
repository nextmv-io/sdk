package mip

import (
	"github.com/nextmv-io/sdk/run/schema"
	"github.com/nextmv-io/sdk/run/statistics"
)

// CustomResult is an example of custom statistics that can be added to the
// output and used in experiments.
type CustomResult struct {
	Columns int    `json:"columns,omitempty"`
	Rows    int    `json:"rows,omitempty"`
	Status  string `json:"status,omitempty"`
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

// DefaultStatistics creates default custom statistics for a given solution.
func DefaultStatistics(model Model, solution Solution) CustomResult {
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

	return CustomResult{
		Status:  status,
		Columns: len(model.Vars()),
		Rows:    len(model.Constraints()),
	}
}
