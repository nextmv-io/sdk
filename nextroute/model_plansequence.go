package nextroute

// ModelPlanSequence is a plan that has multiple stops which all need to be
// planned together on the same vehicle. The precedence has to be equal to the
// order of the stops in the cluster. Other stops are allowed to be planned
// in between the stops in the cluster.
type ModelPlanSequence interface {
	ModelPlanCluster

	// IsFirst returns true if the stop is the first stop in the sequence.
	IsFirst(stop ModelStop) bool
	// IsLast returns true if the stop is the last stop in the sequence.
	IsLast(stop ModelStop) bool

	// Next returns the next stop in the sequence. If the stop is the last
	// stop in the sequence, the stop itself is returned. If the stop is not
	// part of the sequence, nil is returned.
	Next(stop ModelStop) ModelStop
	// Previous returns the previous stop in the sequence. If the stop is the
	// first stop in the sequence, the stop itself is returned. If the stop is
	// not part of the sequence, nil is returned.
	Previous(stop ModelStop) ModelStop
}
