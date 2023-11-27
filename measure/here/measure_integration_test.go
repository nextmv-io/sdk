//go:build integration
// +build integration

package here

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var apiKey = os.Getenv("HERE_API_KEY")

// This is not a classical "test" so much as a script to get a gut-check that we're passing
// props on to the HERE measure correctly
func TestHereMeasuresIntegration(t *testing.T) {
	points := []route.Point{
		// Boulder
		{-105.276995, 40.023353},
		// Boulder
		{-105.272575, 40.004848},
		// Longmont
		{-105.099540, 40.178720},
		// Denver
		{-104.989334, 39.734760},
	}
	loc, err := time.LoadLocation("MST")
	if err != nil {
		t.Fatalf("loading location: %v", err)
	}
	tests := []struct {
		description string
		points      []route.Point
		opts        []MatrixOption
	}{
		{
			description: "normal",
			points:      points,
			opts: []MatrixOption{
				WithDepartureTime(time.Time{}),
			},
		},
		{
			description: "departure time",
			points:      points,
			opts: []MatrixOption{
				WithDepartureTime(time.Date(2021, 12, 23, 5, 30, 0, 0, loc)),
			},
		},
		{
			description: "bike",
			points:      points,
			opts: []MatrixOption{
				WithTransportMode(TransportModeBicycle),
			},
		},
		{
			description: "avoid features and areas",
			points:      points,
			opts: []MatrixOption{
				WithAvoidFeatures([]Feature{
					TollRoad,
					ControlledAccessHighway,
				}),
				WithAvoidAreas([]BoundingBox{
					{
						North: 39.902022,
						South: 39.857233,
						East:  -104.971970,
						West:  -105.122689,
					},
				}),
			},
		},
		{
			description: "truck profile with axle groups",
			points:      points,
			opts: []MatrixOption{
				WithTransportMode(TransportModeTruck),
				WithTruckProfile(Truck{
					ShippedHazardousGoods: []HazardousGood{Poison, Explosive},
					GrossWeight:           36287,
					TunnelCategory:        TunnelCategoryNone,
					Type:                  TruckTypeTractor,
					AxleCount:             5,
					Height:                411,
					Width:                 259,
					Length:                2194,
					TrailerCount:          2,
					WeightPerAxleGroup: &WeightPerAxleGroup{
						Tandem: 14514,
					},
				}),
			},
		},
		{
			description: "truck profile with weight per axle",
			points:      points,
			opts: []MatrixOption{
				WithTransportMode(TransportModeTruck),
				WithTruckProfile(Truck{
					ShippedHazardousGoods: []HazardousGood{Poison, Explosive},
					GrossWeight:           36287,
					TunnelCategory:        TunnelCategoryNone,
					Type:                  TruckTypeTractor,
					AxleCount:             5,
					WeightPerAxle:         7257,
					Height:                411,
					Width:                 259,
					Length:                2194,
					TrailerCount:          2,
				}),
			},
		},
	}
	for i, test := range tests {
		cli := NewClient(apiKey)
		ctx := context.Background()
		distances, durations, err := cli.DistanceDurationMatrices(ctx, test.points, test.opts...)
		if err != nil {
			t.Errorf("[%d] %s: getting matrices: %v", i, test.description, err)
		}
		log.Println(test.description)
		log.Println("Distances")
		for i := 0; i < len(test.points); i++ {
			for j := 0; j < len(test.points); j++ {
				fmt.Printf("%.0f ", distances.Cost(i, j))
			}
			fmt.Println("")
		}
		log.Println("Durations")
		for i := 0; i < len(test.points); i++ {
			for j := 0; j < len(test.points); j++ {
				fmt.Printf("%.0f ", durations.Cost(i, j))
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}
