package store

import (
	"encoding/json"
	"reflect"
)

// Slice manages an immutable slice container of some type in a Store.
type Slice[T any] interface {
	/*
		Append one or more values to the end of a Slice.

			s1 := store.New()
			x := store.NewSlice(s1, 1, 2, 3) // [1, 2, 3]
			s2 := s1.Apply(x.Append(4, 5))
			x.Slice(s2) // [1, 2, 3, 4, 5]
	*/
	Append(value T, values ...T) Change

	/*
		Get an index of a Slice.

			s := store.New()
			x := store.NewSlice(s, 1, 2, 3)
			x.Get(s, 2) // 3
	*/
	Get(Store, int) T

	/*
		Insert one or more values at an index in a Slice.

			s1 := store.New()
			x := store.NewSlice(s, "a", "b", "c")
			s2 := s1.Apply(s.Insert(2, "d", "e")) // [a, b, d, e, c]
	*/
	Insert(index int, value T, values ...T) Change

	/*
		Len returns the length of a Slice.

			s := store.New()
			x := store.NewSlice(s, 1, 2, 3)
			x.Len(s) // 3
	*/
	Len(Store) int

	/*
		Prepend one or more values at the beginning of a Slice.

			s1 := store.New()
			x := store.NewSlice(s, 1, 2, 3)    // [1, 2, 3]
			s2 := s1.Apply(x.Prepend(4, 5)) // [4, 5, 1, 2, 3]
	*/
	Prepend(value T, values ...T) Change

	/*
		Remove a sub-Slice from a starting to an ending index.

			s1 := store.New()
			x := store.NewSlice(s, 1, 2, 3) // [1, 2, 3]
			s2 := s1.Apply(x.Remove(1))  // [1, 3]
	*/
	Remove(start, end int) Change

	/*
		Set a value by index.

			s1 := store.New()
			x := store.NewSlice(s, "a", "b", "c") // [a, b, c]
			s2 := s1.Apply(x.Set(1, "d"))      // [a, d, c]
	*/
	Set(int, T) Change

	/*
		Slice representation that is mutable.

			s := store.New()
			x := store.NewSlice(s, 1, 2, 3)
			x.Slice(s) // []int{1, 2, 3}
	*/
	Slice(Store) []T
}

/*
NewSlice returns a new NewSlice and stores it in a Store.

	s := store.New()
	x := store.NewSlice[int](s)        // []int{}
	y := store.NewSlice(s, 3.14, 2.72) // []float64{3.14, 2.72}
*/
func NewSlice[T any](s Store, values ...T) Slice[T] {
	connect()
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
