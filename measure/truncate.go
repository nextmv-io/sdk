package measure

import "encoding/json"

// Truncate the cost of some other measure.
func Truncate(m ByIndex, lower, upper float64) ByIndex {
	return truncate{byIndex: m, lower: lower, upper: upper}
}

type truncate struct {
	byIndex ByIndex
	lower   float64
	upper   float64
}

func (t truncate) Cost(from, to int) float64 {
	cost := t.byIndex.Cost(from, to)
	if cost < t.lower {
		return t.lower
	}
	if cost > t.upper {
		return t.upper
	}
	return cost
}

func (t truncate) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "truncate",
		"measure": t.byIndex,
		"lower":   t.lower,
		"upper":   t.upper,
	})
}
