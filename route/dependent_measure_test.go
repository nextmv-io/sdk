package route_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/route"
)

func ExampleDependentIndexed() {
	t := time.Now()
	indexed1 := route.Constant(1000)
	indexed2 := route.Scale(indexed1, 2)
	measures := []measure.ByIndex{indexed1, indexed2}

	endTimes := []time.Time{t.Add(2000 * time.Second), t.Add(5000 * time.Second)}
	etds := []int{
		int(t.Add(500 * time.Second).Unix()),
		int(t.Add(3000 * time.Second).Unix()),
	}
	c := route.NewTimeDependentMeasuresClient(measures, endTimes, measures[0])
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

func TestDependentIndexed(te *testing.T) {
	t := time.Now()
	indexed1 := route.Constant(100)
	indexed2 := route.Scale(indexed1, 2)
	indexed3 := route.Scale(indexed2, 2)
	measures := []measure.ByIndex{indexed1, indexed2, indexed3}

	endTimes := []time.Time{
		t.Add(150 * time.Second),
		t.Add(175 * time.Second),
		t.Add(5000 * time.Second),
	}
	etds := []int{
		int(t.Add(100 * time.Second).Unix()),
		int(t.Add(3000 * time.Second).Unix()),
	}
	c := route.NewTimeDependentMeasuresClient(measures, endTimes, measures[0])
	dependentMeasure := c.DependentByIndex()

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
		te.Errorf("overlapping dependent measure, got:%f, want:%f", got1, want1)
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
		te.Errorf("overlapping dependent measure, got:%f, want:%f", got2, want2)
	}
}
