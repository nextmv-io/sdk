package model

import (
	"hash/maphash"
	"strconv"
)

// Identifier needs to be implemented by any type to be used with MultiMap. The
// value returned by ID() needs to be unique for every instance of Identifier.
type Identifier interface {
	ID() string
}

// MultiMap is a map with an n-dimensional index.
type MultiMap[T any, T2 Identifier] interface {
	// Get an element from the MultiMap, given an n-dimensional index or
	// adds it by using the create function if it hasn't been created previously.
	Get(identifiers ...T2) T
	// Length of the elements in the MultiMap.
	Length() int
}

// multiMap implements MultiMap.
type multiMap[T any, T2 Identifier] struct {
	hash   maphash.Hash
	m      map[uint64]T
	create func(...T2) T
	sets   [][]T2
}

// NewMultiMap creates a new MultiMap. It takes a create function that is
// responsible for creating a new entity of T based on a given n-dimensional index.
// The second argument is a variable number of sets, one set per dimension of
// the index.
func NewMultiMap[T any, T2 Identifier](
	create func(...T2) T,
	sets ...[]T2,
) MultiMap[T, T2] {
	return &multiMap[T, T2]{
		m:      map[uint64]T{},
		sets:   sets,
		create: create,
	}
}
func (m *multiMap[T, T2]) Get(identifiers ...T2) T {
	m.hash.Reset()
	for i, id := range identifiers {
		_, err := m.hash.WriteString(id.ID())
		if err != nil {
			panic(err)
		}
		_, err = m.hash.WriteString(strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
	}
	index := m.hash.Sum64()
	if v, ok := m.m[index]; ok {
		return v
	}
	variable := m.create(identifiers...)
	m.m[index] = variable
	return variable
}
func (m *multiMap[T, T2]) Length() int {
	return len(m.m)
}
