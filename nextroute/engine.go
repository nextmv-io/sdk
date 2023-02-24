// Package nextroute is a package
package nextroute

import "github.com/nextmv-io/sdk/connect"

// Input is the input to a model.
type Input struct {
	Name string `json:"name"`
}

// Model is the input to a solver.
type Model struct {}

// ModelOptions configure the model.
type ModelOptions struct {}

// SolverOptions configure the solver.
type SolverOptions struct {}

// SolveOptions configure a run of a solver, e.g. runtime or gap.
type SolveOptions struct {}

// EngineOptions configure the engine.
type EngineOptions struct {
	ModelOptions
	SolverOptions
	SolveOptions
}

// Solver is able to solve a Model.
type Solver interface {
    Solve(SolveOptions) ([]Solution, error)
}

// Solution is the result of invoking Solve() on a Solver.
type Solution struct {}

// NewModel creates a new model.
func NewModel(input Input, modelOptions ModelOptions) (Model, error) {
    connect.Connect(con, &newModelFunc)
	return newModelFunc(input, modelOptions)
}

// NewSolver creates a new solver.
func NewSolver(model Model, solverOptions SolverOptions) (Solver, error) {
    connect.Connect(con, &newSolverFunc)
	return newSolverFunc(model, solverOptions)
}

var (
	con           = connect.NewConnector("sdk", "Nextroute")
	newSolverFunc func(Model, SolverOptions) (Solver, error)
	newModelFunc  func(Input, ModelOptions) (Model, error)
)
