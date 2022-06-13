package store

import (
	"encoding/json"
	"reflect"
)

// Slice manages an immutable slice container of some type in a Store.
type Slice[T any] interface {
	// Append one or more values to the end of a slice.
	Append(value T, values ...T) Change
	// Get an index of a slice.
	Get(Store, int) T
	// Insert one or more values at an index in a slice.
	Insert(index int, value T, values ...T) Change
	// Len returns the length of a slice.
	Len(Store) int
	// Prepend one or more values at the beginning of a slice.
	Prepend(value T, values ...T) Change
	// Remove a subslice from a start to an end index.
	Remove(start, end int) Change
	// Set a value by index.
	Set(int, T) Change
	// Slice representation that is mutable.
	Slice(Store) []T
}

// NewSlice returns a new Slice and stores it in a Store.
func NewSlice[T any](s Store, values ...T) Slice[T] {
	return sliceProxy[T]{slice: newSliceFunc(s, anySlice(values)...)}
}

type sliceProxy[T any] struct {
	slice Slice[any]
}

// Implements Slice

func (s sliceProxy[T]) Append(value T, values ...T) Change {
	return s.slice.Append(value, anySlice(values)...)
}

func (s sliceProxy[T]) Get(store Store, index int) T {
	return s.slice.Get(store, index).(T)
}

func (s sliceProxy[T]) Insert(index int, value T, values ...T) Change {
	return s.slice.Insert(index, value, anySlice(values)...)
}

func (s sliceProxy[T]) Len(store Store) int {
	return s.slice.Len(store)
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

func (s sliceProxy[T]) Slice(store Store) []T {
	sliceAny := s.slice.Slice(store)
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

func anySlice[T any](source []T) []any {
	dest := make([]any, len(source))
	for i, v := range source {
		dest[i] = v
	}
	return dest
}

var newSliceFunc func(Store, ...any) Slice[any]
