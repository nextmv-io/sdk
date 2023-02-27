package nextroute

// ModelObjective is an objective function that can be used to optimize a
// solution.
type ModelObjective interface {
	// EstimateDeltaScore returns the estimated change in the score if the given
	// visit positions are changed.
	EstimateDeltaScore(visitPositions StopPositions) float64

	// Factor returns the factor of the objective.
	Factor() float64

	// Value returns the value of the objective for the given solution.
	Value(solution Solution) float64
}

// ModelObjectives is a slice of model objectives.
type ModelObjectives []ModelObjective

// ModelObjectiveSum is a sum of model objectives.
type ModelObjectiveSum interface {
	ModelObjective

	// Add adds an objective to the sum.
	Add(objective ModelObjective)

	// ModelObjectives returns the model objectives that are part of the sum.
	ModelObjectives() ModelObjectives
}
