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
	ConstraintDataUpdater
	ModelConstraint
	ModelObjective
}

// NewCluster creates a new cluster component. It needs to be added as a
// constraint or as an objective to the model to be taken into account.
func NewCluster() (Cluster, error) {
	connect.Connect(con, &newCluster)
	return newCluster()
}
