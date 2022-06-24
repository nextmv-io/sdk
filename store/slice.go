package store

import (
	"encoding/json"
	"reflect"
)

/*
NewSlice returns a new NewSlice and stores it in a Store.

	s := store.New()
	x := store.NewSlice[int](s)        // []int{}
	y := store.NewSlice(s, 3.14, 2.72) // []float64{3.14, 2.72}
*/
func NewSlice[T any](s Store, values ...T) Slice[T] {
	return sliceProxy[T]{slice: newSliceFunc(s, anySlice(values)...)}
}

type sliceProxy[T any] struct {
	slice Slice[any]
}

// Implements Slice.

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

// Implements fmt.Stringer.

func (s sliceProxy[T]) String() string {
	var x []T
	return reflect.TypeOf(x).String()
}

// Implements json.Marshaler.

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
