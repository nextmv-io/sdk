package factory

import (
	"context"
	"github.com/nextmv-io/sdk/connect"
	"math"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/schema"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

// Format formats a solution in a basic format using the [schema.Output] to
// format a solution.
func Format(
	ctx context.Context,
	options any,
	progressioner alns.Progressioner,
	solutions ...nextroute.Solution,
) runSchema.Output {
	connect.Connect(con, &format)
	return format(ctx, options, progressioner, solutions...)
}

// DefaultCustomResultStatistics creates default custom statistics for a given
// solution.
func DefaultCustomResultStatistics(solution nextroute.Solution) schema.CustomResultStatistics {
	vehicleCount := 0
	maxTravelDuration := 0
	minTravelDuration := math.MaxInt64
	maxDuration := 0
	minDuration := math.MaxInt64
	maxStops := 0
	minStops := math.MaxInt64
	for _, vehicle := range solution.Vehicles() {
		if vehicle.IsEmpty() {
			continue
		}

		vehicleCount++
		duration := vehicle.Duration().Seconds()
		if int(duration) > maxDuration {
			maxDuration = int(duration)
		}
		if int(duration) < minDuration {
			minDuration = int(duration)
		}

		travelDuration := int(vehicle.Last().CumulativeTravelDuration().Seconds())
		if travelDuration > maxTravelDuration {
			maxTravelDuration = travelDuration
		}
		if travelDuration < minTravelDuration {
			minTravelDuration = travelDuration
		}

		stops := vehicle.NumberOfStops()
		if stops > maxStops {
			maxStops = stops
		}
		if stops < minStops {
			minStops = stops
		}
	}

	return schema.CustomResultStatistics{
		ActivatedVehicles: vehicleCount,
		UnplannedStops:    solution.UnPlannedPlanUnits().Size(),
		MaxTravelDuration: maxTravelDuration,
		MaxDuration:       maxDuration,
		MinTravelDuration: minTravelDuration,
		MinDuration:       minDuration,
		MaxStopsInVehicle: maxStops,
		MinStopsInVehicle: minStops,
	}
}
