package route

import (
	"sort"
	"time"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/model"
)

// ByIndexAndTime holds a measure and an endTime up until this measure is to be
// used. ByIndexAndTime is to be used with NewTimeDependentMeasure which a slice
// of ByIndexAndTime.
type byIndexAndTime struct {
	measure ByIndex
	endTime int
}

// ClientOption can pass options to be used with a TimeDependentMeasure client.
type ClientOption func(*client)

// TimeDependentMeasuresClient is an interface to handle time dependent
// measures. It implements a Cost function that takes time into account to
// calculate costs.
type TimeDependentMeasuresClient interface {
	Cost() func(from, to int, data measure.VehicleData) float64
	DependentByIndex() measure.DependentByIndex
}

type client struct {
	measures        []byIndexAndTime
	fallbackMeasure byIndexAndTime
	cache           map[int]*byIndexAndTime
}

// NewTimeDependentMeasuresClient returns a new NewTimeDependentMeasuresClient
// which implements a cost function.
func NewTimeDependentMeasuresClient(
	measures []ByIndex,
	endTimes []time.Time,
	fallback ByIndex,
	opts ...ClientOption,
) TimeDependentMeasuresClient {
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

	c := &client{
		measures: m,
		// The fallback measure will also be used if we are getting a very late
		// ETA for the last stop. To achieve this we max out the time.Time
		// endTime as int.
		fallbackMeasure: byIndexAndTime{
			measure: fallback,
			endTime: model.MaxInt,
		},
		cache: map[int]*byIndexAndTime{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *client) Cost() func(
	from,
	to int,
	data measure.VehicleData,
) float64 {
	return func(from, to int, data measure.VehicleData) float64 {
		if data.Index == -1 {
			return c.fallbackMeasure.measure.Cost(from, to)
		}
		etd := data.Times.EstimatedDeparture[data.Index]
		return c.interpolate(from, to, etd, 0, 1)
	}
}

func (c *client) interpolate(
	from,
	to,
	startTime int,
	prevIncurredCosts,
	partialFactor float64,
) float64 {
	// Use default measure and look for a better one afterwards
	measure := c.fallbackMeasure
	if m, ok := c.cache[startTime]; ok {
		measure = *m
	} else {
		for _, m := range c.measures {
			if startTime < m.endTime {
				measure = m
				c.cache[startTime] = &m
				break
			}
		}
	}

	// Get the drive costs for current measure. The new total costs depend on
	// the previous costs and the part needs to be calculated with new measure.
	rawDriveTime := measure.measure.Cost(from, to)
	interpolatedDriveTime := partialFactor * rawDriveTime
	driveEnd := startTime + int(interpolatedDriveTime)
	if driveEnd <= measure.endTime || measure.endTime == model.MaxInt {
		return prevIncurredCosts + interpolatedDriveTime
	}
	newPartial := float64(measure.endTime-startTime) /
		float64(driveEnd-startTime)
	newCosts := prevIncurredCosts + newPartial*interpolatedDriveTime
	return c.interpolate(
		from,
		to,
		measure.endTime+1,
		newCosts,
		1-newPartial,
	)
}

func (c *client) DependentByIndex() measure.DependentByIndex {
	return DependentIndexed(true, c.Cost())
}

// WithFullCache creates a full cache up front for each second in the
// byIndexAndTime measure. Otherwise this cache will be built on the fly.
func WithFullCache(startTime time.Time) ClientOption {
	return func(c *client) {
		time := int(startTime.Unix())
		for _, measure := range c.measures {
			for i := time; i < measure.endTime+1; i++ {
				c.cache[i] = &measure
			}
		}
	}
}
