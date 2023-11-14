package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// Cluster is both a constraint and an objective that limits/prefers the
// vehicles a plan cluster will be added to. If used as a constraint a plan
// cluster can only be added to a vehicles whose centroid is closer to the plan
// cluster than the centroid of any other vehicle. In case of using it as an
// objective, those vehicles will be preferred.
type Cluster interface {
	ConstraintStopDataUpdater
	ModelConstraint
	ModelObjective

	// IncludeFirst returns whether the first stop of the vehicle is included in the
	// centroid calculation. The centroid is used to determine the distance
	// between a new stop and the cluster.
	IncludeFirst() bool
	// IncludeLast returns whether the last stop of the vehicle is included in
	// the centroid calculation. The centroid is used to determine the distance
	// between a new stop and the cluster.
	IncludeLast() bool

	// SetIncludeFirst sets whether the first stop of the vehicle is included in
	// the centroid calculation. The centroid is used to determine the distance
	// between a new stop and the cluster.
	SetIncludeFirst(includeFirst bool)
	// SetIncludeLast sets whether the last stop of the vehicle is included in
	// the centroid calculation. The centroid is used to determine the distance
	// between a new stop and the cluster.
	SetIncludeLast(includeLast bool)
}

// NewCluster creates a new cluster component. It needs to be added as a
// constraint or as an objective to the model to be taken into account.
// By default, the first and last stop of a vehicle are not included in the
// centroid calculation.
func NewCluster() (Cluster, error) {
	connect.Connect(con, &newCluster)
	return newCluster()
}
