package kmeans

import "github.com/nextmv-io/sdk/measure"

// Cluster is a cluster of points. A cluster is defined by a centroid
// and a set of points. The centroid is the center of the cluster.
type Cluster interface {
	// Centroid returns the centroid of the cluster.
	Centroid() measure.Point

	// Indices returns the indices of the points in the cluster. The
	// index refers to the index of the point in the model.
	Indices() []int

	// Points returns the points in the cluster. The points are the
	// points in the model at the position Indices.
	Points() []measure.Point

	// ClusterModel returns the cluster model of which the invoking
	// cluster is an instance.
	ClusterModel() ClusterModel

	// WithinClusterSumOfSquares returns the sum of the squared
	// distances between each point and the cluster centroid.
	WithinClusterSumOfSquares() float64
}

// Solution is a solution to a k-means clustering problem.
type Solution interface {
	// Clusters returns the clusters derived from the solution.
	Clusters() []Cluster
	// Unassigned returns the points that were not assigned to any
	// cluster.
	Unassigned() []measure.Point
}
