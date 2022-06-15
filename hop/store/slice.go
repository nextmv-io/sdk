package store

import (
	"encoding/json"
	"reflect"

	"github.com/nextmv-io/sdk/hop/store/types"
)

// Slice returns a new Slice and stores it in a Store.
func Slice[T any](s types.Store, values ...T) types.Slice[T] {
	return sliceProxy[T]{slice: sliceFunc(s, anySlice(values)...)}
}

type sliceProxy[T any] struct {
	slice types.Slice[any]
}

// Implements types.Slice.

func (s sliceProxy[T]) Append(value T, values ...T) types.Change {
	return s.slice.Append(value, anySlice(values)...)
}

func (s sliceProxy[T]) Get(store types.Store, index int) T {
	return s.slice.Get(store, index).(T)
}

func (s sliceProxy[T]) Insert(index int, value T, values ...T) types.Change {
	return s.slice.Insert(index, value, anySlice(values)...)
}

func (s sliceProxy[T]) Len(store types.Store) int {
	return s.slice.Len(store)
}

func (s sliceProxy[T]) Prepend(value T, values ...T) types.Change {
	return s.slice.Prepend(value, anySlice(values)...)
}

func (s sliceProxy[T]) Remove(start, end int) types.Change {
	return s.slice.Remove(start, end)
}

func (s sliceProxy[T]) Set(index int, value T) types.Change {
	return s.slice.Set(index, value)
}

func (s sliceProxy[T]) Slice(store types.Store) []T {
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

var sliceFunc func(types.Store, ...any) types.Slice[any]
