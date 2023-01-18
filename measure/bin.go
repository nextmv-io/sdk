package measure

import (
	"encoding/json"
	"sync/atomic"
)

// Bin is a measure that selects from a slice of indexed measures. Logic
// defined in the selector function determines which measure is used in the
// cost calculation.
func Bin(
	measures []ByIndex,
	selector func(from, to int) int,
) ByIndex {
	return &bin{
		measures:      measures,
		selector:      selector,
		measuresCount: make([]uint64, len(measures)),
	}
}

type bin struct {
	measures      []ByIndex
	selector      func(from, to int) int
	measuresCount []uint64
}

func (b *bin) Cost(from, to int) float64 {
	i := b.selector(from, to)
	atomic.AddUint64(&b.measuresCount[i], 1)
	return b.measures[i].Cost(from, to)
}

func (b *bin) MarshalJSON() ([]byte, error) {
	measures := []map[string]any{}

	for i, count := range b.measuresCount {
		measure := b.measures[i]
		measures = append(measures, map[string]any{
			"measure": measure,
			"count":   count,
		})
	}

	return json.Marshal(map[string]any{
		"measures": measures,
		"type":     "bin",
	})
}
