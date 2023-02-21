package route

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/model"
)

// ByIndexAndTime holds a Measure and an EndTime (exclusive) up until this
// measure is to be used as a Unix timestamp. ByIndexAndTime is to be used with
// NewTimeDependentMeasure which a slice of ByIndexAndTime.
type ByIndexAndTime struct {
	Measure ByIndex
	EndTime int
}

// ClientOption can pass options to be used with a TimeDependentMeasure client.
type ClientOption func(*client)

type client struct {
	measures        []ByIndexAndTime
	fallbackMeasure ByIndexAndTime
	cache           sync.Map
}

// NewTimeDependentMeasure returns a new NewTimeDependentMeasure
// which implements a cost function.
// It takes byIndexAndTime measures, where each measure is given with an endTime
// (exclusive) up until the measure will be used and a fallback measure.
func NewTimeDependentMeasure(
	measures []ByIndexAndTime,
	fallback ByIndex,
	opts ...ClientOption,
) (measure.DependentByIndex, error) {
	measuresCopy := make([]ByIndexAndTime, len(measures))
	copy(measuresCopy, measures)
	sort.SliceStable(measuresCopy, func(i, j int) bool {
		return measuresCopy[i].EndTime < measuresCopy[j].EndTime
	})

	if fallback == nil {
		return nil, errors.New("a fallback measure must be given")
	}

	c := &client{
		measures: measuresCopy,
		// The fallback measure will also be used if we are getting a very late
		// ETA for the last stop. To achieve this we max out the time.Time
		// endTime as int.
		fallbackMeasure: ByIndexAndTime{
			Measure: fallback,
			EndTime: model.MaxInt,
		},
		cache: sync.Map{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return DependentIndexed(true, c.Cost()), nil
}

func (c *client) Cost() func(
	from,
	to int,
	data *measure.VehicleData,
) float64 {
	return func(from, to int, data *measure.VehicleData) float64 {
		if data.Index == -1 || data.Times.EstimatedDeparture == nil {
			return c.fallbackMeasure.Measure.Cost(from, to)
		}
		etd := data.Times.EstimatedDeparture[data.Index]
		cost := c.interpolate(from, to, etd, 0, 1)
		return cost
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
	if m, ok := c.cache.Load(startTime); ok {
		measure = m.(ByIndexAndTime)
	} else {
		for _, m := range c.measures {
			if startTime < m.EndTime {
				measure = m
				c.cache.Store(startTime, m)
				break
			}
		}
	}

	// Get the drive costs for current measure. The new total costs depend on
	// the previous costs and the part needs to be calculated with new measure.
	rawDriveTime := measure.Measure.Cost(from, to)
	interpolatedDriveTime := partialFactor * rawDriveTime
	driveEnd := float64(startTime) + interpolatedDriveTime
	if driveEnd < float64(measure.EndTime) || measure.EndTime == model.MaxInt {
		return prevIncurredCosts + interpolatedDriveTime
	}
	newPartialFactor := float64(measure.EndTime-startTime) /
		(driveEnd - float64(startTime))
	newCosts := prevIncurredCosts + newPartialFactor*interpolatedDriveTime
	return c.interpolate(
		from,
		to,
		measure.EndTime,
		newCosts,
		partialFactor-newPartialFactor,
	)
}

// WithFullCache creates a full cache up front for each second in the
// byIndexAndTime measure. Otherwise this cache will be built on the fly.
func WithFullCache(startTime time.Time) ClientOption {
	return func(c *client) {
		cacheTimes(startTime, c)
	}
}

func cacheTimes(startTime time.Time, c *client) {
	time := int(startTime.Unix())
	for _, measure := range c.measures {
		for i := time; i < measure.EndTime; i++ {
			c.cache.Store(i, measure)
		}
		time = measure.EndTime
	}
}
