package nextroute

// ModelPlanMultipleStops is a plan that has multiple stops that must be planned
// together on the same vehicle. Unlike [ModelPlanSequence], the order of the
// stops is not enforced on the vehicle, meaning the only behavior enforced is
// that the stops must be planned together.
type ModelPlanMultipleStops interface {
	ModelPlanCluster
}
