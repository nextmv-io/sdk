package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewSolveObserver returns a new solve observer. The solve observer can be
// used to observe the solve process. The solve observer writes to the given
// file name.
//
// The solve observer writes the following columns:
// - The first column is the type of event. The type of event can be one of:
//   - `+` for a plan event.
//   - `-` for an unplan event.
//   - `~` for a failed event.
//   - `b` for a new best solution event.
//   - 'o' for the objective definition event.
//   - `r` for a reset event.
//   - The second column is the time since the solve process started in
//     nanoseconds.
//   - The third column is the step. The step is used to group events
//     together. The step is incremented for each unplan event.
//   - The next columns are dependent on the event type.
//   - For an objective definition event, the next columns are:
//   - The number of terms of the objective.
//   - For each term of the objective the factor and name of the objective.
//   - For a plan event, the next columns are:
//   - The previous stop ID.
//   - The stop ID.
//   - The next stop ID.
//   - For each term of the objective the score
//   - The score after planning.
//   - The estimated impact on the objective of planning.
//   - For an unplan event, the next columns are:
//   - The previous stop ID.
//   - The stop ID.
//   - The next stop ID.
//   - For each term of the objective the score
//   - The score after un-planning.
//   - For a failed event, the next column is the reason for the failure.
//   - For a reset event, the next columns are:
//   - The score of the work solution.
//   - The score of the solution resetting to.
//   - For a new best solution event, the next columns are:
//   - For each term of the objective the score
//   - Last column is the score of the new best solution
//
// This observer should not be used in combination with parallel solving. The
// output will be garbled. To use the observer add it to the model the
// following way where
//
//	solver, err := nextroute.NewSolver(model, solverOptions)
//	solveObserver, err := observers.NewSolveObserver("solve.log")
//	solveObserver.Register(solver)
func NewSolveObserver(fileName string) (SolveObserver, error) {
	connect.Connect(con, &newSolveObserver)
	return newSolveObserver(fileName)
}

// SolveObserver is an observer for the solve process.
type SolveObserver interface {
	SolutionObserver
	SolutionUnPlanObserver

	// Register registers the solver to the solve observer.
	Register(solver Solver) error

	// FileName returns the file name of the solve observer.
	FileName() string
}
