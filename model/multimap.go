package model

import (
	"hash/maphash"
)

// Identifier needs to be implemented by any type to be used with MultiMap. The
// value returned by ID() needs to be unique for every instance of Identifier.
type Identifier interface {
	ID() int
}

// MultiMap is a map with an n-dimensional index.
type MultiMap[T any] interface {
	Get(identifiers ...Identifier) T
	Slice(identifiers ...Identifier) []T
	Length() int
}

// multiMap is a map with an n-dimensional index.
type multiMap[T any] struct {
	hash   *maphash.Hash
	m      map[uint64]T
	create func(...Identifier) (T, error)
	sets   [][]Identifier
}

// NewMultiMap creates a new MultiMap. It takes a create function that is
// responsible to create a new entity of T based on a given n-dimensional index.
// The second argument is a variable number of sets, one set per dimension of
// the index.
func NewMultiMap[T any](
	create func(...Identifier) (T, error),
	sets ...[]Identifier,
) MultiMap[T] {
	return multiMap[T]{
		m:      map[uint64]T{},
		hash:   &maphash.Hash{},
		sets:   sets,
		create: create,
	}
}

// Get retrieves an element from the MultiMap, given an n-dimensional index.
func (m multiMap[T]) Get(identifiers ...Identifier) T {
	m.hash.Reset()
	for _, id := range identifiers {
		m.hash.WriteString(string(id.ID()))
	}
	index := m.hash.Sum64()
	v, ok := m.m[index]
	if !ok {
		variable, err := m.create(identifiers...)
		if err != nil {
			panic(err) // or return the error?
		}
		m.m[index] = variable
		v = variable
	}
	return v
}

// Length returns the number of elements in the MultiMap.
func (m multiMap[T]) Length() int {
	return len(m.m)
}

// Slice works as follows:
// Assume your index is based on (vehicles,stops) and you want all elements
// dealing with stop[0], then call Slice(nil,stop[0]).
func (m multiMap[T]) Slice(identifiers ...Identifier) []T {
	sets := make([][]Identifier, len(m.sets))
	for i, identifier := range identifiers {
		// nil is sentinal for "get all identifiers in the set"
		if identifier == nil {
			sets[i] = m.sets[i]
		} else {
			sets[i] = []Identifier{identifier}
		}
	}
	indices := cartN(sets...)

	returnList := []T{}
	for _, index := range indices {
		returnList = append(returnList, m.Get(index...))
	}
	return returnList
}

// cartN computes the cartesian product of the given sets.
func cartN[T any](a ...[]T) [][]T {
	c := 1
	for _, a := range a {
		c *= len(a)
	}
	if c == 0 {
		return nil
	}
	p := make([][]T, c)
	b := make([]T, c*len(a))
	n := make([]int, len(a))
	s := 0
	for i := range p {
		e := s + len(a)
		pi := b[s:e]
		p[i] = pi
		s = e
		for j, n := range n {
			pi[j] = a[j][n]
		}
		for j := len(n) - 1; j >= 0; j-- {
			n[j]++
			if n[j] < len(a[j]) {
				break
			}
			n[j] = 0
		}
	}
	return p
}
