package nextroute

// ModelPlanSequence is a plan that has multiple stops which all need to be
// planned together on the same vehicle. The precedence has to be equal to the
// order of the stops in the cluster. Other stops are allowed to be planned
// in between the stops in the cluster.
type ModelPlanSequence interface {
	ModelPlanCluster
}
