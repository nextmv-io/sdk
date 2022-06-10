package context

// A Vector stores a contiguous list of some type.
type Vector[T any] interface {
	// Get an index of a vector.
	Get(Context, int) T
	// Len returns the length of a vector,
	Len(Context) int
	// Slice representation of a vector.
	Slice(Context) []T

	// Append one or more values to the end of a vector.
	Append(value T, values ...T) Change
	// Insert one or more values at an index in a vector.
	Insert(index int, value T, values ...T) Change
	// Prepend one or more values at the beginning of a vector.
	Prepend(value T, values ...T) Change
	// Remove a subvector from a start to an end index.
	Remove(start, end int) Change
	// Set a value by index.
	Set(int, T) Change
}

// NewVector returns a new Vector.
func NewVector[T any](values ...T) Vector[T] {
	connect()
	return vectorProxy[T]{vector: newVectorFunc()}
}

type vectorProxy[T any] struct {
	vector Vector[any]
}

func (v vectorProxy[T]) Get(ctx Context, index int) T {
	return v.vector.Get(ctx, index).(T)
}

func (v vectorProxy[T]) Len(ctx Context) int {
	return v.vector.Len(ctx)
}

func (v vectorProxy[T]) Slice(ctx Context) []T {
	sliceAny := v.vector.Slice(ctx)
	sliceT := make([]T, len(sliceAny))
	for i, v := range sliceAny {
		sliceT[i] = v.(T)
	}
	return sliceT
}

func (v vectorProxy[T]) Append(value T, values ...T) Change {
	return v.vector.Append(value, anySlice(values)...)
}

func (v vectorProxy[T]) Insert(index int, value T, values ...T) Change {
	return v.vector.Insert(index, value, anySlice(values)...)
}

func (v vectorProxy[T]) Prepend(value T, values ...T) Change {
	return v.vector.Prepend(value, anySlice(values)...)
}

func (v vectorProxy[T]) Remove(start, end int) Change {
	return v.vector.Remove(start, end)
}

func (v vectorProxy[T]) Set(index int, value T) Change {
	return v.vector.Set(index, value)
}

var newVectorFunc func(...any) Vector[any]

func anySlice[T any](source []T) []any {
	dest := make([]any, len(source))
	for i, v := range source {
		dest[i] = v
	}
	return dest
}
