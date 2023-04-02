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

// ImmutableSolutionPlanClusterCollection is a collection of solution plan clusters.
type ImmutableSolutionPlanClusterCollection interface {
	// Iterator returns a channel that can be used to iterate over the solution
	// plan clusters in the collection.
	// If you break out of the for loop before the channel is closed,
	// the goroutine launched by the Iterator() method will be blocked forever,
	// waiting to send the next element on the channel. This can lead to a
	// goroutine leak and potentially exhaust the system resources. Therefore,
	// it is recommended to always use the following pattern:
	//    iter := collection.Iterator()
	//    for {
	//        element, ok := <-iter
	//        if !ok {
	//            break
	//        }
	//        // do something with element, potentially break out of the loop
	//    }
	//    close(iter)
	Iterator(quit <-chan struct{}) <-chan SolutionPlanCluster
	// RandomDraw returns a random sample of n different solution plan clusters.
	RandomDraw(n int) SolutionPlanClusters
	// RandomElement returns a random solution plan cluster.
	RandomElement() SolutionPlanCluster
	// Size return the number of solution plan clusters in the collection.
	Size() int
	// SolutionStops returns the solution stops of all the clusters in the
	// collection.
	SolutionStops() SolutionStops
	// SolutionPlanCluster returns the solution plan clusters in the collection
	// which correspond to the given model plan cluster. If no such solution
	// plan cluster is found, nil is returned.
	SolutionPlanCluster(modelPlanCluster ModelPlanCluster) SolutionPlanCluster
	// SolutionPlanClusters returns the solution plan clusters in the collection.
	// The returned slice is a defensive copy of the internal slice, so
	// modifying it will not affect the collection.
	SolutionPlanClusters() SolutionPlanClusters
}

// SolutionPlanClusterCollection is a collection of solution plan clusters.
type SolutionPlanClusterCollection interface {
	ImmutableSolutionPlanClusterCollection
	// Add adds a solution plan cluster to the collection.
	Add(solutionPlanCluster SolutionPlanCluster)
	// Remove removes a solution plan cluster from the collection.
	Remove(cluster SolutionPlanCluster)
}
