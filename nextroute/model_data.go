package nextroute

// ModelData is a data interface available on several model constructs. It
// allows to attach arbitrary data to a model construct.
type ModelData interface {
	// Data returns the data.
	Data() any
	// SetData sets the data.
	SetData(any)
}
