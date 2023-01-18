package kmeans

import (
	"github.com/nextmv-io/sdk/connect"
)

// An Option configures a k-means model.
type Option func(Model) error

// MaximumPoints defines the maximum number of points that can be assigned to
// a cluster. maximumPoints must be equal to the number of slices in the model.
func MaximumPoints(maximumPoints []int) Option {
	connect.Connect(con, &maximumPointsFunc)
	return maximumPointsFunc(maximumPoints)
}

// ExcludedPoints defines the points that are excluded from the
// cluster. The excludedPoints contain the index of the point in the model to be
// excluded for each cluster. The length of excludedPoints must be equal to the
// number of clusters in the model. The slice of excludedPoints for a
// cluster can be any size of indices of points.
func ExcludedPoints(excludedPoints [][]int) Option {
	connect.Connect(con, &excludedPointsFunc)
	return excludedPointsFunc(excludedPoints)
}

// MaximumSumValue defines the maximum value constraint for the cluster. The
// maximum value constraint limits the number of points that can be assigned to
// the cluster based on the value of a point defined by values. Values must be
// the same length as points in the model.
func MaximumSumValue(
	maximumValue []int,
	values [][]int,
) Option {
	connect.Connect(con, &maximumSumValueFunc)
	return maximumSumValueFunc(maximumValue, values)
}

var (
	maximumPointsFunc   func([]int) Option
	excludedPointsFunc  func([][]int) Option
	maximumSumValueFunc func([]int, [][]int) Option
)
