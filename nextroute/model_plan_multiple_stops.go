package nextroute

// ModelPlanMultipleStops is a [ModelPlanCluster] that has multiple stops that
// must be planned together on the same vehicle. The plan contains a
// [DirectedAcyclicGraph] (DAG) that restricts the sequence in which the stops
// can be planned on the vehicle.
type ModelPlanMultipleStops interface {
	ModelPlanCluster

	// DirectedAcyclicGraph returns the [DirectedAcyclicGraph] of the plan
	// cluster.
	DirectedAcyclicGraph() DirectedAcyclicGraph
}
