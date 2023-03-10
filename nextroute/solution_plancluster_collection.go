package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"math/rand"
)

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

func NewSolutionPlanClusterCollection(
	source *rand.Rand,
	planClusters SolutionPlanClusters,
) SolutionPlanClusterCollection {
	connect.Connect(con, &newSolutionPlanClusterCollection)
	return newSolutionPlanClusterCollection(source, planClusters)
}
