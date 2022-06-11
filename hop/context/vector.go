package context

import (
	"encoding/json"
	"reflect"
)

// A Vector stores a contiguous list of some type.
type Vector[T any] interface {
	// Append one or more values to the end of a vector.
	Append(value T, values ...T) Change
	// Data underlying a vector.
	Data(Context) []T
	// Get an index of a vector.
	Get(Context, int) T
	// Insert one or more values at an index in a vector.
	Insert(index int, value T, values ...T) Change
	// Len returns the length of a vector,
	Len(Context) int
	// Prepend one or more values at the beginning of a vector.
	Prepend(value T, values ...T) Change
	// Remove a subvector from a start to an end index.
	Remove(start, end int) Change
	// Set a value by index.
	Set(int, T) Change
}

// NewVector returns a new Vector.
func NewVector[T any](ctx Context, values ...T) Vector[T] {
	connect()
	return vectorProxy[T]{vector: newVectorFunc(ctx, anySlice(values)...)}
}

type vectorProxy[T any] struct {
	vector Vector[any]
}

func (v vectorProxy[T]) Append(value T, values ...T) Change {
	return v.vector.Append(value, anySlice(values)...)
}

func (v vectorProxy[T]) Data(ctx Context) []T {
	sliceAny := v.vector.Data(ctx)
	sliceT := make([]T, len(sliceAny))
	for i, v := range sliceAny {
		sliceT[i] = v.(T)
	}
	return sliceT
}

func (v vectorProxy[T]) Get(ctx Context, index int) T {
	return v.vector.Get(ctx, index).(T)
}

func (v vectorProxy[T]) Insert(index int, value T, values ...T) Change {
	return v.vector.Insert(index, value, anySlice(values)...)
}

func (v vectorProxy[T]) Len(ctx Context) int {
	return v.vector.Len(ctx)
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

func (d vectorProxy[T]) String() string {
	var x []T
	return reflect.TypeOf(x).String()
}

func (d vectorProxy[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

var newVectorFunc func(Context, ...any) Vector[any]

func anySlice[T any](source []T) []any {
	dest := make([]any, len(source))
	for i, v := range source {
		dest[i] = v
	}
	return dest
}
