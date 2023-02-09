package measure

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
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
	Index      int
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

// ByIndexAndTime holds a measure and an endTime up until this measure is to be
// used. ByIndexAndTime is to be used with NewTimeDependentMeasure which a slice
// of ByIndexAndTime.
type byIndexAndTime struct {
	measure ByIndex
	endTime int
}

// TimeDependentCostFunc returns a new TimeDependentCostFunc.
func TimeDependentCostFunc(
	measures []ByIndex,
	endTimes []time.Time,
) func(from, to int, data VehicleData) float64 {
	m := make([]byIndexAndTime, len(measures))
	for i := range measures {
		m[i] = byIndexAndTime{
			measure: measures[i],
			endTime: int(endTimes[i].Unix()),
		}
	}
	sort.SliceStable(m, func(i, j int) bool {
		return m[i].endTime < m[j].endTime
	})
	return func(from, to int, data VehicleData) float64 {
		time := data.Times.EstimatedDeparture[data.Index]
		for _, measure := range m {
			if time < measure.endTime {
				return measure.measure.Cost(from, to)
			}
		}
		panic(fmt.Sprintf("no measure for time '%d' was provided", time))
	}
}
