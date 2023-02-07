package measure

import (
	"encoding/json"
)

// DependentIndexed is a measure uses a custom cost func to calculate parameter
// dependent costs for connecting to points by index.
func DependentIndexed(
	timeDependent bool,
	cost func(
		from,
		to int,
		data VehicleData,
	) float64,
) DependentByIndex {
	return &dependentIndexed{
		cost:          cost,
		timeDependent: timeDependent,
	}
}

// VehicleData holds vehicle specific data, including times by index (ETA, ETD
// and ETS), a vehicle id, the vehicle's route and the solution value for that
// vehicle.
type VehicleData struct {
	Times      Times
	VehicleID  string
	Route      []int
	RouteValue int
}

type dependentIndexed struct {
	cost func(
		from,
		to int,
		data VehicleData,
	) float64
	timeDependent bool
}

func (b *dependentIndexed) Cost(
	from,
	to int,
	data VehicleData,
) float64 {
	return b.cost(from, to, data)
}

func (b *dependentIndexed) TimeDependent() bool {
	return b.timeDependent
}

func (b *dependentIndexed) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":          "dependentIndexed",
		"timeDependent": b.timeDependent,
	})
}
