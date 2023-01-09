package kmeans

import (
	"time"

	"github.com/nextmv-io/sdk/measure"
)

// Cluster is a cluster of points. A cluster is defined by a centroid
// and a set of points. The centroid is the center of the cluster.
type Cluster interface {
	// Centroid returns the centroid of the cluster.
	Centroid() measure.Point

	// ClusterModel returns the cluster model of which the invoking
	// cluster is an instance.
	ClusterModel() ClusterModel

	// Indices returns the indices of the points in the cluster. The
	// index refers to the index of the point in the model.
	Indices() []int

	// Points returns the points in the cluster. The points are the
	// points in the model at the position Indices.
	Points() []measure.Point

	// WithinClusterSumOfSquares returns the sum of the squared
	// distances between each point and the cluster centroid.
	WithinClusterSumOfSquares() float64
}

// Solution is a solution to a k-means clustering problem.
type Solution interface {
	// Clusters returns the clusters derived from the solution.
	Clusters() []Cluster
	// Feasible returns true if the solution is feasible. A solution is
	// feasible if the solver was able to find a solution that satisfied
	// the constraints of the model.
	Feasible() bool
	// RunTime returns the time it took to derive the solution.
	RunTime() time.Duration
	// Unassigned returns the points that were not assigned to any
	// cluster.
	Unassigned() []measure.Point
	// UnassignedIndices returns the indices of the points that were
	// not assigned to any cluster.
	UnassignedIndices() []int
}
