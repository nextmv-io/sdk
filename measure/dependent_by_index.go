package measure

import (
	"encoding/json"
	"fmt"
	"sort"
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

// ByIndexAndTime holds a measure and an endTime up until this measure is to be
// used. ByIndexAndTime is to be used with NewTimeDependentMeasure which a slice
// of ByIndexAndTime.
type ByIndexAndTime struct {
	m       ByIndex
	endTime int
}

type timeDependentMeasures []ByIndexAndTime

// NewTimeDependentMeasure is used to create a time dependent cost function. The
// passed slice is sorted by endTime first.
func NewTimeDependentMeasure(measures []ByIndexAndTime) timeDependentMeasures {
	sort.SliceStable(measures, func(i, j int) bool {
		return measures[i].endTime < measures[j].endTime
	})
	return measures
}

// TimeDependentCosts returns the cost for connecting to points by index. It
// selects the measure by looping of the sorted MeasureByTime slice and picks
// the first measure that satisfies the condition endTime < time. If no such
// endTime is given, the function panics.
func (t timeDependentMeasures) TimeDependentCosts() func(
	from,
	to int,
	data VehicleData,
) float64 {
	return func(from, to int, data VehicleData) float64 {
		time := data.Times.EstimatedDeparture[from]
		for _, measure := range t {
			if measure.endTime < time {
				return measure.m.Cost(from, to)
			}
		}
		panic(fmt.Sprintf("no measure for time '%d' was provided", time))
	}
}
