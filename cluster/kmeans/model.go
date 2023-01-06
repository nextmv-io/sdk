package kmeans

import (
	"github.com/nextmv-io/sdk/measure"
)

// Model is a model of a k-means clustering problem.
type Model interface {
	// ClusterModels returns the cluster models. The cluster models
	// define the constraints for the clusters.
	ClusterModels() []ClusterModel
	// Points returns the points to be clustered.
	Points() []measure.Point
}
