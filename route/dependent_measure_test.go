package route

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/store"
)

func ExampleDependentIndexed() {
	t := time.Now()
	indexed1 := Constant(1000)
	indexed2 := Scale(indexed1, 2)
	measures := []measure.ByIndex{indexed1, indexed2}

	endTimes := []time.Time{
		t.Add(2000 * time.Second),
		t.Add(5000 * time.Second),
	}

	byIndex := make([]ByIndexAndTime, len(measures))
	for i, m := range measures {
		byIndex[i] = ByIndexAndTime{
			Measure: m,
			EndTime: int(endTimes[i].Unix()),
		}
	}

	etds := []int{
		int(t.Add(500 * time.Second).Unix()),
		int(t.Add(3000 * time.Second).Unix()),
	}

	c, err := NewTimeDependentMeasuresClient(byIndex, measures[0])
	if err != nil {
		panic(err)
	}
	dependentMeasure := c.DependentByIndex()
	fmt.Println(dependentMeasure.Cost(0, 1, measure.VehicleData{
		Index: 0,
		Times: measure.Times{
			EstimatedDeparture: etds,
		},
	}))
	fmt.Println(dependentMeasure.Cost(1, 0, measure.VehicleData{
		Index: 1,
		Times: measure.Times{
			EstimatedDeparture: etds,
		},
	}))
	// Output:
	// 1000
	// 2000
}

func TestDependentIndexed(t *testing.T) {
	etds, dependentMeasure := dependentMeasures(t)

	// The ETD is selected by index (here: 0) and because of short time ranges
	// by end times the costs are calculated of 3 measures.
	// First 50% of the first measure -> 50
	// Second: 25% of the second measure -> 25
	// Third: 25% of the third measure -> 100
	want1 := 175.0
	got1 := dependentMeasure.Cost(0, 1, measure.VehicleData{
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
	got2 := dependentMeasure.Cost(1, 0, measure.VehicleData{
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

	c, err := NewTimeDependentMeasuresClient(byIndex, measures[0])
	if err != nil {
		t.Errorf(err.Error())
	}
	dependentMeasure := c.DependentByIndex()

	return etds, dependentMeasure
}

func TestValueFuncOption(t *testing.T) {
	stops := []Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	_, dependentMeasure := dependentMeasures(t)
	dependentMeasures := make([]DependentByIndex, len(vehicles))
	for i := 0; i < len(vehicles); i++ {
		dependentMeasures[i] = dependentMeasure
	}

	// Declare the router and its solver.
	router, err := NewRouter(
		stops,
		vehicles,
		Threads(1),
		ValueFunctionMeasures(dependentMeasures),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
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
		cache:           make(map[int]*ByIndexAndTime),
	}

	cacheTimes(startTime, &c)
	want := 15
	if len(c.cache) != want {
		t.Errorf("cached items, got:%d, want:%d", len(c.cache), want)
	}

	for k, v := range c.cache {
		if k < 5 {
			if v != &c.measures[0] {
				t.Errorf(
					"caches measure is not correct, got:%d, want:%d", v,
					&c.measures[0],
				)
			}
			continue
		}
		if k < 10 {
			if v != &c.measures[1] {
				t.Errorf(
					"caches measure is not correct, got:%d, want:%d", v,
					&c.measures[1],
				)
			}
			continue
		}
		if k < 15 {
			if v != &c.measures[2] {
				t.Errorf(
					"caches measure is not correct, got:%d, want:%d", v,
					&c.measures[2],
				)
			}
			continue
		}
	}
}
