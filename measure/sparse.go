package measure

import (
	"encoding/json"
	"sync/atomic"
)

// Sparse measure returns pre-computed costs between two locations without
// requiring a full data set. If two locations do not have an associated cost,
// then a backup measure is used.
func Sparse(m ByIndex, arcs map[int]map[int]float64) ByIndex {
	return &sparse{m: m, arcs: arcs}
}

// DebugSparse returns a Sparse that when marshalled will include debugging
// information describing the number of queries for elements included in (and
// not included in) the matrix.
func DebugSparse(m ByIndex, arcs map[int]map[int]float64) ByIndex {
	return &sparse{m: m, arcs: arcs, debugMode: true}
}

type sparse struct {
	m         ByIndex
	debugMode bool
	arcs      map[int]map[int]float64

	counts struct {
		hit  uint64
		miss uint64
	}
}

func (s *sparse) Cost(from, to int) float64 {
	if m, ok := s.arcs[from]; ok {
		if c, ok := m[to]; ok {
			if s.debugMode {
				atomic.AddUint64(&s.counts.hit, 1)
			}
			return c
		}
	}

	if s.debugMode {
		atomic.AddUint64(&s.counts.miss, 1)
	}
	return s.m.Cost(from, to)
}

func (s sparse) MarshalJSON() ([]byte, error) {
	output := map[string]any{
		"measure": s.m,
		"arcs":    s.arcs,
		"type":    "sparse",
	}

	if s.debugMode {
		output["counts"] = map[string]uint64{
			"hit":  s.counts.hit,
			"miss": s.counts.miss,
		}
	}

	return json.Marshal(output)
}
