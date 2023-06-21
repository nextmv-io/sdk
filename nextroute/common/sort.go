package common

import (
	"sort"
)

// Compare compares two numbers and returns -1 if a < b, 0 if a == b,
// or 1 if a > b.
func Compare[T Numbers](a, b T) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// CompareFunction is generic type for a function that compares two values
// and returns -1 if a < b, 0 if a == b, or 1 if a > b.
type CompareFunction[T any] func(a, b T) int

// CompareOpposite returns a CompareFunction that compares two values
// in opposite order of the given CompareFunction.
func CompareOpposite[T any](c CompareFunction[T]) CompareFunction[T] {
	return func(a, b T) int {
		return c(b, a)
	}
}

// CompareFunctionDefined returns a CompareFunction that compares two values
// by applying the given function to each value and comparing the results.
func CompareFunctionDefined[N Numbers, T any](a, b T, f func(T) N) int {
	return Compare(f(a), f(b))
}

// CompareBy returns a CompareFunction that compares two values
// by applying the given function to each value and comparing the results.
func CompareBy[N Numbers, T any](f func(T) N) CompareFunction[T] {
	return func(a, b T) int {
		return Compare(f(a), f(b))
	}
}

// ComposeCompare returns a CompareFunction that compares two values by
// applying the given CompareFunction to each value and comparing the results.
// If the CompareFunction returns 0, then the next CompareFunction is applied.
// If all CompareFunctions return 0, then 0 is returned.
func ComposeCompare[T any](c CompareFunction[T], compares ...CompareFunction[T]) CompareFunction[T] {
	return func(a, b T) int {
		if r := c(a, b); r != 0 {
			return r
		}
		for _, compare := range compares {
			if r := compare(a, b); r != 0 {
				return r
			}
		}
		return 0
	}
}

// Sort sorts the elements of a slice in increasing order using one or more
// CompareFunctions.
func Sort[T any](
	elements []T,
	compare CompareFunction[T],
	compares ...CompareFunction[T],
) []T {
	e := DefensiveCopy(elements)
	if len(e) < 2 {
		return e
	}
	b := byCompareSort[T]{
		elements: e,
		compare:  ComposeCompare(compare, compares...),
	}
	sort.Sort(b)
	return b.elements
}

type byCompareSort[T any] struct {
	elements []T
	compare  CompareFunction[T]
}

// Len is the number of elements in the collection.
func (a byCompareSort[T]) Len() int {
	return len(a.elements)
}

// Swap swaps the elements with indexes i and j.
func (a byCompareSort[T]) Swap(i, j int) {
	a.elements[i], a.elements[j] = a.elements[j], a.elements[i]
}

// Less reports whether the element with index i should sort before the
// element with index j.
func (a byCompareSort[T]) Less(i, j int) bool {
	return a.compare(a.elements[i], a.elements[j]) < 0
}
