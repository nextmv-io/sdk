package context

import (
	"encoding/json"
	"reflect"
)

// A Slice manages an immutable slice of some type.
type Slice[T any] interface {
	// Append one or more values to the end of a slice.
	Append(value T, values ...T) Change
	// Get an index of a slice.
	Get(Context, int) T
	// Insert one or more values at an index in a slice.
	Insert(index int, value T, values ...T) Change
	// Len returns the length of a slice,
	Len(Context) int
	// Prepend one or more values at the beginning of a slice.
	Prepend(value T, values ...T) Change
	// Remove a subslice from a start to an end index.
	Remove(start, end int) Change
	// Set a value by index.
	Set(int, T) Change
	// Slice representation that is mutable.
	Slice(Context) []T
}

// NewSlice returns an imutable slice container.
func NewSlice[T any](ctx Context, values ...T) Slice[T] {
	return sliceProxy[T]{slice: newSliceFunc(ctx, anySlice(values)...)}
}

type sliceProxy[T any] struct {
	slice Slice[any]
}

// Implements Slice

func (s sliceProxy[T]) Append(value T, values ...T) Change {
	return s.slice.Append(value, anySlice(values)...)
}

func (s sliceProxy[T]) Get(ctx Context, index int) T {
	return s.slice.Get(ctx, index).(T)
}

func (s sliceProxy[T]) Insert(index int, value T, values ...T) Change {
	return s.slice.Insert(index, value, anySlice(values)...)
}

func (s sliceProxy[T]) Len(ctx Context) int {
	return s.slice.Len(ctx)
}

func (s sliceProxy[T]) Prepend(value T, values ...T) Change {
	return s.slice.Prepend(value, anySlice(values)...)
}

func (s sliceProxy[T]) Remove(start, end int) Change {
	return s.slice.Remove(start, end)
}

func (s sliceProxy[T]) Set(index int, value T) Change {
	return s.slice.Set(index, value)
}

func (s sliceProxy[T]) Slice(ctx Context) []T {
	sliceAny := s.slice.Slice(ctx)
	sliceT := make([]T, len(sliceAny))
	for i, s := range sliceAny {
		sliceT[i] = s.(T)
	}
	return sliceT
}

// Implements fmt.Stringer

func (s sliceProxy[T]) String() string {
	var x []T
	return reflect.TypeOf(x).String()
}

// Implements json.Marshaler

func (s sliceProxy[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

var newSliceFunc func(Context, ...any) Slice[any]

func anySlice[T any](source []T) []any {
	dest := make([]any, len(source))
	for i, v := range source {
		dest[i] = v
	}
	return dest
}
