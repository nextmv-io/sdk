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

		// Process unassigned stops.
		unassigned := make([]schema.StopOutput, 0)

		for u := range s.UnPlannedPlanUnits().Iterator(quit) {
			for _, v := range u.ModelPlanUnit().Stops() {
				unassigned = append(unassigned, schema.StopOutput{
					ID: v.ID(),
					Location: schema.Location{
						Lon: v.Location().Longitude(),
						Lat: v.Location().Latitude(),
					},
				})
			}
		}

		objective := makeObjective(s)
		solutionOutput := schema.SolutionOutput{
			Unplanned: unassigned,
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

	r := reflect.ValueOf(options)
	f := reflect.Indirect(r).FieldByName("Format")
	if !f.IsValid() || !f.CanInterface() {
		return output
	}
	if format, ok := f.Interface().(FormatOptions); ok {
		if format.Disable.Progression {
			return output
		}
	}

	seriesData := make([]statistics.DataPoint, 0)
	for _, progression := range progressionValues {
		seriesData = append(seriesData, statistics.DataPoint{
			X: statistics.Float64(progression.ElapsedSeconds),
			Y: statistics.Float64(progression.Value),
		})
	}
	lastProgressionElement := progressionValues[len(progressionValues)-1]
	lastProgressionValue := statistics.Float64(lastProgressionElement.Value)
	output.Statistics.SeriesData = &statistics.SeriesData{
		Value: statistics.Series{
			Name:       output.Solutions[len(output.Solutions)-1].(schema.SolutionOutput).Objective.Name,
			DataPoints: seriesData,
		},
	}
	output.Statistics.Result = &statistics.Result{
		Duration: &lastProgressionElement.ElapsedSeconds,
		Value:    &lastProgressionValue,
	}

	return output
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
