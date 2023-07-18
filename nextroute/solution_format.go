package nextroute

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	runSchema "github.com/nextmv-io/sdk/run/schema"
	"github.com/nextmv-io/sdk/run/statistics"
)

// FormatOptions are the options that influence the format of the output.
type FormatOptions struct {
	Disable struct {
		Progression bool `json:"progression" usage:"disable the progression series"`
	} `json:"disable"`
}

// CustomStatistics is an example of custom statistics that can be added to the
// output and used in experiments.
type CustomStatistics struct {
	UsedVehicles      int `json:"used_vehicles,omitempty"`
	UnplannedStops    int `json:"unplanned_stops,omitempty"`
	MaxTravelDuration int `json:"max_travel_duration,omitempty"`
	MaxDuration       int `json:"max_duration,omitempty"`
	MinTravelDuration int `json:"min_travel_duration,omitempty"`
	MinDuration       int `json:"min_duration,omitempty"`
	MaxStopsInRoute   int `json:"max_stops_in_route,omitempty"`
	MinStopsInRoute   int `json:"min_stops_in_route,omitempty"`
}

// Format formats a solution in a basic format.
func Format(
	ctx context.Context,
	options any,
	progressioner alns.Progressioner,
	solutions ...Solution,
) runSchema.Output {
	solutionOutputs := make([]schema.SolutionOutput, 0, len(solutions))
	quit := make(chan struct{})
	defer close(quit)
	for _, s := range solutions {
		// Process solutions of vehicles.
		solutionVehicles := s.Vehicles()
		vehicles := make([]schema.VehicleOutput, len(solutionVehicles))
		for v, state := range solutionVehicles {
			vehicles[v] = toVehicleOutput(state)
		}

		// Process unplanned stops.
		unplanned := make([]schema.StopOutput, 0)

		for u := range s.UnPlannedPlanUnits().Iterator(quit) {
			for _, modelStop := range u.ModelPlanUnit().Stops() {
				stop := schema.StopOutput{
					ID: modelStop.ID(),
					Location: schema.Location{
						Lon: modelStop.Location().Longitude(),
						Lat: modelStop.Location().Latitude(),
					},
				}
				inputStop := modelStop.Data().(schema.Stop)
				if inputStop.CustomData != nil {
					stop.CustomData = inputStop.CustomData
				}
				unplanned = append(unplanned, stop)
			}
		}

		objective := makeObjective(s)
		solutionOutput := schema.SolutionOutput{
			Unplanned: unplanned,
			Vehicles:  vehicles,
			Objective: objective,
		}
		solutionOutputs = append(solutionOutputs, solutionOutput)
	}

	output := runSchema.NewOutput(options, solutionOutputs...)

	// initialize statistics
	output.Statistics = statistics.NewStatistics()

	// set run duration if available
	if start, ok := ctx.Value(run.Start).(time.Time); ok {
		duration := time.Since(start).Seconds()
		output.Statistics.Run = &statistics.Run{
			Duration: &duration,
		}
	}

	progressionValues := progressioner.Progression()
	if len(progressionValues) == 0 {
		return output
	}

	seriesData := make([]statistics.DataPoint, 0, len(progressionValues))
	iterationsSeriesData := make([]statistics.DataPoint, 0, len(progressionValues))
	for _, progression := range progressionValues {
		seriesData = append(seriesData, statistics.DataPoint{
			X: statistics.Float64(progression.ElapsedSeconds),
			Y: statistics.Float64(progression.Value),
		})
		iterationsSeriesData = append(iterationsSeriesData, statistics.DataPoint{
			X: statistics.Float64(progression.ElapsedSeconds),
			Y: statistics.Float64(progression.Iterations),
		})
	}
	lastProgressionElement := progressionValues[len(progressionValues)-1]
	lastProgressionValue := statistics.Float64(lastProgressionElement.Value)

	output.Statistics.Result = &statistics.Result{
		Duration: &lastProgressionElement.ElapsedSeconds,
		Value:    &lastProgressionValue,
	}

	r := reflect.ValueOf(options)
	f := reflect.Indirect(r).FieldByName("Format")
	if f.IsValid() && f.CanInterface() {
		if format, ok := f.Interface().(FormatOptions); ok {
			if format.Disable.Progression {
				return output
			}
		}
	}

	output.Statistics.SeriesData = &statistics.SeriesData{
		Value: statistics.Series{
			Name:       output.Solutions[len(output.Solutions)-1].(schema.SolutionOutput).Objective.Name,
			DataPoints: seriesData,
		},
	}
	output.Statistics.SeriesData.Custom = append(output.Statistics.SeriesData.Custom, statistics.Series{
		Name:       "iterations",
		DataPoints: iterationsSeriesData,
	})

	return output
}

// DefaultStatistics creates default custom statistics for a given solution.
func DefaultStatistics(solution Solution) CustomStatistics {
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

	return CustomStatistics{
		UsedVehicles:      vehicleCount,
		UnplannedStops:    solution.UnPlannedPlanUnits().Size(),
		MaxTravelDuration: maxTravelDuration,
		MaxDuration:       maxDuration,
		MinTravelDuration: minTravelDuration,
		MinDuration:       minDuration,
		MaxStopsInRoute:   maxStops,
		MinStopsInRoute:   minStops,
	}
}

// toVehicleOutput constructs the output state of a vehicle.
func toVehicleOutput(vehicle SolutionVehicle) schema.VehicleOutput {
	solutionStops := vehicle.SolutionStops()
	vehicleName := vehicle.ModelVehicle().ID()
	hasUserDefinedStartTime := vehicle.First().Start() !=
		vehicle.ModelVehicle().Model().Epoch()

	stops := make([]schema.PlannedStopOutput, 0, len(solutionStops))
	cumulativeStopsDuration := 0
	for _, solutionStop := range solutionStops {
		if !solutionStop.ModelStop().Location().IsValid() {
			continue
		}

		stop := schema.PlannedStopOutput{
			Stop: schema.StopOutput{
				ID: solutionStop.ModelStop().ID(),
				Location: schema.Location{
					Lon: solutionStop.ModelStop().Location().Longitude(),
					Lat: solutionStop.ModelStop().Location().Latitude(),
				},
			},
		}

		if inputStop, ok := solutionStop.ModelStop().Data().(schema.Stop); ok {
			if inputStop.CustomData != nil {
				stop.Stop.CustomData = inputStop.CustomData
			}
		}

		stop = setTimes(solutionStop, stop, hasUserDefinedStartTime)
		stops = append(stops, stop)
		cumulativeStopsDuration += int(solutionStop.End().Sub(solutionStop.Start()).Seconds())
	}

	vehicleOutput := schema.VehicleOutput{
		ID:                  vehicleName,
		Route:               stops,
		RouteDuration:       int(vehicle.Duration().Seconds()),
		RouteTravelDuration: int(vehicle.Last().CumulativeTravelDuration().Seconds()),
		RouteStopsDuration:  cumulativeStopsDuration,
	}

	inputVehicle := vehicle.ModelVehicle().Data().(schema.Vehicle)
	if inputVehicle.CustomData != nil {
		vehicleOutput.CustomData = inputVehicle.CustomData
	}

	vehicleOutput.RouteWaitingDuration = vehicleOutput.RouteDuration -
		vehicleOutput.RouteTravelDuration -
		vehicleOutput.RouteStopsDuration

	return vehicleOutput
}

// setTimes sets all the time-related fields of a stop in the output.
func setTimes(
	solutionStop SolutionStop,
	stopOutput schema.PlannedStopOutput,
	hasUserDefinedStartTime bool,
) schema.PlannedStopOutput {
	// we need to access the timezone via the vehicle of the model
	timezoneLocation := solutionStop.ModelStop().Model().
		Vehicle(solutionStop.VehicleIndex()).Start().Location()
	arrival := solutionStop.Arrival().In(timezoneLocation)
	departure := solutionStop.End().In(timezoneLocation)
	service := solutionStop.Start().In(timezoneLocation)
	if hasUserDefinedStartTime {
		stopOutput.ArrivalTime = &arrival
		stopOutput.EndTime = &departure
		stopOutput.StartTime = &service
	}

	stopOutput.TravelDuration = int(solutionStop.TravelDuration().Seconds())
	stopOutput.Duration = int(solutionStop.End().Sub(solutionStop.Start()).Seconds())
	stopOutput.WaitingDuration = int(solutionStop.Start().Sub(solutionStop.Arrival()).Seconds())
	stopOutput.CumulativeTravelDuration = int(solutionStop.CumulativeTravelDuration().Seconds())

	data := solutionStop.ModelStop().Data()
	inputStop, ok := data.(schema.Stop)
	if !ok || inputStop.TargetArrivalTime == nil {
		return stopOutput
	}
	targetArrivalTime := inputStop.TargetArrivalTime.In(timezoneLocation)
	stopOutput.TargetArrivalTime = &targetArrivalTime

	if inputStop.EarlyArrivalTimePenalty != nil {
		earliness := int(math.Max(inputStop.TargetArrivalTime.Sub(arrival).Seconds(), 0.0))
		stopOutput.EarlyArrivalDuration = earliness
	}

	if inputStop.LateArrivalTimePenalty != nil {
		lateness := int(math.Max(arrival.Sub(*inputStop.TargetArrivalTime).Seconds(), 0.0))
		stopOutput.LateArrivalDuration = lateness
	}

	return stopOutput
}

func makeObjective(s Solution) schema.ObjectiveOutput {
	sumObjective := s.Model().Objective()
	terms := make([]schema.ObjectiveOutput, len(sumObjective.Terms()))
	for i, term := range sumObjective.Terms() {
		value := s.ObjectiveValue(term.Objective())
		terms[i] = schema.ObjectiveOutput{
			Name:   fmt.Sprintf("%v", term.Objective()),
			Factor: term.Factor(),
			Base:   value / term.Factor(),
			Value:  value,
		}
	}

	return schema.ObjectiveOutput{
		Name:       fmt.Sprintf("%v", sumObjective),
		Objectives: terms,
		Value:      s.ObjectiveValue(sumObjective),
	}
}
