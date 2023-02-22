package measure

import (
	"encoding/json"
)

// DependentIndexed is a measure that uses a custom cost function to calculate
// parameter dependent costs for connecting two points by index. If the measures
// are time dependent all future stops in the sequence will be fully
// recalculated. Otherwise there will be a constant shift to achieve better
// performance.
func DependentIndexed(
	timeDependent bool,
	cost func(
		from,
		to int,
		data *VehicleData,
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
	VehicleID  string
	Times      Times
	Route      []int
	Index      int
	RouteValue int
}

type dependentIndexed struct {
	cost func(
		from,
		to int,
		data *VehicleData,
	) float64
	timeDependent bool
}

func (b *dependentIndexed) Cost(
	from,
	to int,
	data *VehicleData,
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
