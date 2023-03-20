package nextroute

import "github.com/nextmv-io/sdk/nextroute/common"

// PlanClusterType is the type of the plan cluster.
type PlanClusterType int

const (
	// SingleStop is a cluster that contains a single stop.
	SingleStop PlanClusterType = iota
)

// ModelPlanCluster is a cluster of stops in a plan. A cluster is a set of stops
// that are required to be planned together. For example, a cluster can be a
// pickup and a delivery stop that are required to be planned together.
type ModelPlanCluster interface {
	// Centroid returns the centroid of the cluster. The centroid is the
	// average location of all stops in the cluster.
	Centroid() common.Location

	// Index returns the index of the cluster.
	Index() int

	// Stops returns the stops in the cluster.
	Stops() ModelStops

	// Type returns the type of the cluster.
	Type() PlanClusterType
}

// ModelPlanClusters is a slice of plan clusters.
type ModelPlanClusters []ModelPlanCluster
