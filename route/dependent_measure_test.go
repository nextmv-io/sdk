package route

import (
	"sync"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/measure"
)

func TestDependentIndexed(t *testing.T) {
	etds, dependentMeasure := dependentMeasures(t)

	// The ETD is selected by index (here: 0) and because of short time ranges
	// by end times the costs are calculated of 3 measures.
	// First 50% of the first measure -> 50
	// Second: 25% of the second measure -> 25
	// Third: 25% of the third measure -> 100
	want1 := 175.0
	got1 := dependentMeasure.Cost(0, 1, &measure.VehicleData{
		Index: 0,
		Times: measure.Times{
			EstimatedDeparture: etds,
		},
	})

	if got1 != want1 {
		t.Errorf("overlapping dependent measure, got:%f, want:%f", got1, want1)
	}

	// The ETD selected is a late start so that only the third measure is needed
	// to calculate the costs:
	// 100% of the third measure --> 400
	want2 := 400.0
	got2 := dependentMeasure.Cost(1, 0, &measure.VehicleData{
		Index: 1,
		Times: measure.Times{
			EstimatedDeparture: etds,
		},
	})

	if got2 != want2 {
		t.Errorf("overlapping dependent measure, got:%f, want:%f", got2, want2)
	}
}

func dependentMeasures(t *testing.T) ([]int, measure.DependentByIndex) {
	startTime := time.Now()
	indexed1 := Constant(100)
	indexed2 := Scale(indexed1, 2)
	indexed3 := Scale(indexed2, 2)
	measures := []measure.ByIndex{indexed1, indexed2, indexed3}

	endTimes := []time.Time{
		startTime.Add(150 * time.Second),
		startTime.Add(175 * time.Second),
		startTime.Add(5000 * time.Second),
	}

	byIndex := make([]ByIndexAndTime, len(measures))
	for i, m := range measures {
		byIndex[i] = ByIndexAndTime{
			Measure: m,
			EndTime: int(endTimes[i].Unix()),
		}
	}

	etds := []int{
		int(startTime.Add(100 * time.Second).Unix()),
		int(startTime.Add(3000 * time.Second).Unix()),
	}

	dependentMeasure, err := NewTimeDependentMeasure(byIndex, measures[0])
	if err != nil {
		t.Errorf(err.Error())
	}

	return etds, dependentMeasure
}

func TestCache(t *testing.T) {
	startTime := time.Now()
	indexed1 := Constant(100)
	indexed2 := Scale(indexed1, 2)
	indexed3 := Scale(indexed2, 2)
	measures := []measure.ByIndex{indexed1, indexed2, indexed3}

	endTimes := []time.Time{
		startTime.Add(5 * time.Second),
		startTime.Add(10 * time.Second),
		startTime.Add(15 * time.Second),
	}

	byIndex := make([]ByIndexAndTime, len(measures))
	for i, m := range measures {
		byIndex[i] = ByIndexAndTime{
			Measure: m,
			EndTime: int(endTimes[i].Unix()),
		}
	}

	c := client{
		measures:        byIndex,
		fallbackMeasure: byIndex[0],
		cache:           sync.Map{},
	}

	cacheTimes(startTime, &c)
	want := 15
	length := 0
	c.cache.Range(func(key, value any) bool {
		length++
		return true
	})

	if length != want {
		t.Errorf("cached items, got:%d, want:%d", length, want)
	}
}
