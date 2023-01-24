package measure

import (
	"encoding/json"
)

// DependentIndexed is a measure uses a custom cost func to calculate parameter
// dependent costs for connecting to points by index.
func DependentIndexed(
	cost func(
		from,
		to int,
		times Times,
		id string,
		route []int,
		value float64,
	) float64,
) DependentByIndex {
	return &dependentIndexed{
		cost: cost,
	}
}

type dependentIndexed struct {
	cost func(
		from,
		to int,
		times Times,
		id string,
		route []int,
		value float64,
	) float64
}

func (b *dependentIndexed) Cost(
	from,
	to int,
	times Times,
	id string,
	route []int,
	value float64,
) float64 {
	return b.cost(from, to, times, id, route, value)
}

func (b *dependentIndexed) MarshalJSON() ([]byte, error) {
	measures := []map[string]any{}

	return json.Marshal(map[string]any{
		"measures": measures,
		"type":     "dependentIndexed",
	})
}
