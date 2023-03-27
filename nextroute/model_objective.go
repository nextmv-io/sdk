package nextroute

import "github.com/nextmv-io/sdk/connect"

// ObjectiveDataUpdater is the interface than can be used by an objective if
// it wants to store data with each stop in a solution.
type ObjectiveDataUpdater interface {
	// UpdateObjectiveData is called when a stop is added to a solution. The
	// solutionStop has all it's expression values set and this function can use
	// them to update the objective data for the stop. The data returned can be
	// used by the estimate function and can be retrieved by the
	// SolutionStop.ObjectiveValue function.
	UpdateObjectiveData(s SolutionStop) Copier
}

// ModelObjective is an objective function that can be used to optimize a
// solution.
type ModelObjective interface {
	RegisteredModelExpressions

	// EstimateDeltaValue returns the estimated change in the score if the given
	// visit positions are changed.
	EstimateDeltaValue(visitPositions StopPositions) float64

	// Index returns the index of the objective. The index is unique for each
	// objective. The index is used to identify the objective in the
	// solution.
	Index() int

	// Value returns the value of the objective for the given solution.
	Value(solution Solution) float64
}

// ModelObjectives is a slice of model objectives.
type ModelObjectives []ModelObjective

// ModelObjectiveSum is a sum of model objectives.
type ModelObjectiveSum interface {
	ModelObjective

	// Add adds an objective to the sum.
	Add(factor float64, objective ModelObjective) error

	// ModelObjectives returns the model objectives that are part of the sum.
	ModelObjectives() ModelObjectives
}

// NewModelObjectiveIndex returns the next unique objective index.
// This function can be used to create a unique index for a custom
// objective.
func NewModelObjectiveIndex() int {
	connect.Connect(con, &newModelObjectiveIndex)
	return newModelObjectiveIndex()
}
