package nextroute

// SolutionPlanCluster is a cluster of stops that are planned to be visited by
// a vehicle.
type SolutionPlanCluster interface {
	// IsPlanned returns true if the cluster is planned.
	IsPlanned() bool

	// ModelPlanCluster returns the ModelPlanCluster this clustes is
	// based upon.
	ModelPlanCluster() ModelPlanCluster

	// Solution returns the solution this cluster is part of.
	Solution() Solution
	// SolutionStop returns the solution stop for the given model stop.
	// Will panic if the stop is not part of the cluster.
	SolutionStop(stop ModelStop) SolutionStop
	// SolutionStops returns the solution stops in this cluster.
	SolutionStops() SolutionStops

	// UnPlan un-plans the cluster by removing the underlying solution stops
	// from the solution. Returns true if the cluster was unplanned
	// successfully, false if the cluster was not unplanned successfully. A
	// cluster is not successful if it did not result in a change in the
	// solution without violating any hard constraints.
	UnPlan() (bool, error)
}

// SolutionPlanClusters is a slice of SolutionPlanCluster.
type SolutionPlanClusters []SolutionPlanCluster
