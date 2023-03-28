package nextroute

import (
	"math/rand"

	"github.com/nextmv-io/sdk/connect"
)

// NewSolutionPlanClusterCollection creates a new SolutionPlanClusterCollection.
// A SolutionPlanClusterCollection is a collection of solution plan clusters. It
// can be used to randomly draw a sample of solution plan clusters and remove
// solution plan clusters from the collection.
func NewSolutionPlanClusterCollection(
	source *rand.Rand,
	planClusters SolutionPlanClusters,
) SolutionPlanClusterCollection {
	connect.Connect(con, &newSolutionPlanClusterCollection)
	return newSolutionPlanClusterCollection(source, planClusters)
}

// SolutionPlanClusterCollection is a collection of solution plan clusters.
type SolutionPlanClusterCollection interface {
	// RandomDraw returns a random sample of n different solution plan clusters.
	RandomDraw(n int) SolutionPlanClusters
	// RandomElement returns a random solution plan cluster.
	RandomElement() SolutionPlanCluster
	// Remove removes a solution plan cluster from the collection.
	Remove(cluster SolutionPlanCluster)
	// Size return the number of solution plan clusters in the collection.
	Size() int
}
