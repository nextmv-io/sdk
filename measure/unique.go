package measure

import "encoding/json"

// Unique returns a ByIndex that uses a reference slice to map the indices of a
// point to the index of the measure.
// m represents a matrix of unique points.
// references maps a stop (by index) to an index in m.
func Unique(m ByIndex, references []int) ByIndex {
	return &unique{m: m, references: references}
}

// Struct that implements the ByIndex interface and holds additional data to be
// able to do so.
type unique struct {
	m          ByIndex
	references []int
}

// Cost returns the cost between two locations.
func (u *unique) Cost(from, to int) float64 {
	return u.m.Cost(u.references[from], u.references[to])
}

// Triangular returns whether the measure is triangular.
func (u *unique) Triangular() bool {
	return IsTriangular(u.m)
}

// MarshalJSON returns the JSON encoding of the measure.
func (u *unique) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"measure":    u.m,
		"type":       "unique",
		"references": u.references,
	})
}
