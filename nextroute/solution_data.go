package nextroute

// SolutionData is a data interface available on several solution constructs. It
// allows to attach arbitrary data to a solution construct.
type SolutionData interface {
	// Data returns the data.
	Data() Copier
	// SetData sets the data.
	SetData(Copier)
	// CopyData returns a deep copy of the data.
	CopyData() SolutionData
}
