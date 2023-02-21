package route

import (
	"errors"
	"sort"

	"github.com/nextmv-io/sdk/model"
)

// ByIndexAndTime holds a Measure and an EndTime (exclusive) up until this
// measure is to be used as a Unix timestamp. ByIndexAndTime is to be used with
// NewTimeDependentMeasure which a slice of ByIndexAndTime.
type ByIndexAndTime struct {
	Measure ByIndex
	EndTime int
}

type client struct {
	measures        []ByIndexAndTime
	fallbackMeasure ByIndexAndTime
	cache           map[int]ByIndexAndTime
}

// NewTimeDependentMeasure returns a new NewTimeDependentMeasure
// which implements a cost function.
// It takes a startTime (e.g. vehicle start) byIndexAndTime measures, where each
// measure is given with an endTime (exclusive) up until the measure will be
// used and a fallback measure.
func NewTimeDependentMeasure(
	startTime int,
	measures []ByIndexAndTime,
	fallback ByIndex,
) (DependentByIndex, error) {
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
		cache: make(map[int]ByIndexAndTime),
	}

	cacheTimes(startTime, c)

	return DependentIndexed(true, c.cost()), nil
}

func (c *client) cost() func(
	from,
	to int,
	data *VehicleData,
) float64 {
	return func(from, to int, data *VehicleData) float64 {
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
	if m, ok := c.cache[startTime]; ok {
		measure = m
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

func cacheTimes(startTime int, c *client) {
	counter := 0
	for _, measure := range c.measures {
		for i := startTime; i < measure.EndTime; i++ {
			c.cache[i] = measure
			counter++
		}
		startTime = measure.EndTime
	}
}
