package measure

import (
	"encoding/json"
	"fmt"

	"github.com/nextmv-io/sdk/model"
)

// Location measure returns the sum of the cost computed by the passed in
// measure and the specified cost of the 'to' location. This cost is read from
// the passed in costs slice.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Location(
	m ByIndex,
	costs []float64,
	durationGroups DurationGroups,
) (ByIndex, error) {
	if durationGroups == nil {
		durationGroups = make(DurationGroups, 0)
	}
	indexToGroup, err := durationGroups.ToIndexGroup()
	if err != nil {
		return nil, err
	}
	return location{
		costs:          costs,
		m:              m,
		durationGroups: durationGroups,
		indexToGroup:   indexToGroup,
	}, nil
}

type location struct {
	m              ByIndex
	indexToGroup   map[int]DurationGroup // maps a location index to duration group
	costs          []float64
	durationGroups []DurationGroup
}

func (l location) Cost(from, to int) float64 {
	// Additional group cost only apply when driving 'to' a location
	// belonging to a group 'from' a location not belonging to the group.
	costGroup := 0
	if dg, ok := l.indexToGroup[to]; ok {
		if !dg.Group.Contains(from) {
			costGroup = dg.Duration
		}
	}
	return l.m.Cost(from, to) + l.costs[to] + float64(costGroup)
}

func (l location) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"costs":   l.costs,
		"measure": l.m,
		"type":    "location",
	})
}

// DurationGroup groups stops by index which have additional service costs
// attached to them.
//
// Deprecated: This package is deprecated and will be removed in the future.
type DurationGroup struct {
	Group    model.Domain
	Duration int
}

// DurationGroups represents a slice of duration groups. Each duration group is
// used to account for additional service costs whenever a stop of a group is
// approached first.
//
// Deprecated: This package is deprecated and will be removed in the future.
type DurationGroups []DurationGroup

// ToIndexGroup maps a location index to duration group to quickly access the
// corresponding group of a given location.
//
// Deprecated: This package is deprecated and will be removed in the future.
func (dg DurationGroups) ToIndexGroup() (map[int]DurationGroup, error) {
	indexToGroup := make(map[int]DurationGroup)
	// Check if groups are overlapping
	for i, gd1 := range dg {
		for j, gd2 := range dg {
			if i != j {
				if gd1.Group.Overlaps(gd2.Group) {
					return nil, fmt.Errorf("group durations at index %d and %d overlap", i, j)
				}
			}
		}
		// Set up indexToGroup map which maps each location of a stop group
		// to the group it belongs to.
		for _, index := range gd1.Group.Slice() {
			indexToGroup[index] = gd1
		}
	}
	return indexToGroup, nil
}
