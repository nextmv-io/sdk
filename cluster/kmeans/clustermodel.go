package kmeans

// MaximumSumValueConstraint is a constraint on the sum of the values of
// points in a cluster. The values of the points are defined by the
// PointToInt interface.
type MaximumSumValueConstraint interface {
	// MaximumValue returns the maximum value of the sum of the values
	// of the points in the cluster.
	MaximumValue() int
	// Values returns the values of the points in the cluster to use to
	// calculate the sum of the values.
	Values() []int
}

// ClusterModel is a model of a cluster.
type ClusterModel interface {
	// ExcludedPointIndices returns the points that were excluded from the
	// cluster.
	ExcludedPointIndices() []int

	// MaximumPoints returns the maximum number of points that can be
	// assigned to the cluster.
	MaximumPoints() int
	// MaximumSumValueConstraints returns the constraints on the sum of
	// the values of the points in the cluster.
	MaximumSumValueConstraints() []MaximumSumValueConstraint

	// SetExcludedPointIndices sets the points that were excluded from the
	// cluster. The indices contain the index of the point in the
	// model to be excluded.
	SetExcludedPointIndices(indices []int) error
	// SetMaximumPoints sets the maximum number of points that can be
	// assigned to the cluster.
	SetMaximumPoints(maximumPoints int)
	// SetMaximumSumValue sets the maximum value constraint for
	// the cluster. The maximum value constraint limits the number of
	// points that can be assigned to the cluster based on the value
	// of a point defined by values. Values must be the same length as
	// points in the model.
	SetMaximumSumValue(
		maximumValue int,
		values []int,
	) (MaximumSumValueConstraint, error)
}
