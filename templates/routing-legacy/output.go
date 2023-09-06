package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/common"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/statistics"
)

// FleetOutput is the root output structure.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [SolutionOutput] and it's
// entities instead.
type FleetOutput struct {
	Hop struct {
		Version string `json:"version"`
	} `json:"hop"`
	Options    schema.Options `json:"options"`
	Solutions  []any          `json:"solutions"`
	Statistics Statistics     `json:"statistics"`
}

// FleetState is the solution output structure.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [SolutionOutput] and it's
// entities instead.
type FleetState struct {
	Vehicles   []FleetOutputVehicle `json:"vehicles"`
	Unassigned []FleetOutputStop    `json:"unassigned"`
	Summary    FleetOutputSummary   `json:"value_summary"` //nolint:tagliatelle
}

// Statistics adds run stats to the output.
type Statistics struct {
	Time  Time `json:"time"`
	Value int  `json:"value"`
}

// Time holds runtime information.
type Time struct {
	Start          time.Time `json:"start"`
	Elapsed        string    `json:"elapsed"`
	ElapsedSeconds float64   `json:"elapsed_seconds"`
}

// FleetOutputSummary contains a breakdown of the objective value.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [SolutionOutput] and it's
// entities instead.
type FleetOutputSummary struct {
	Value                           float64 `json:"value"`
	TotalTravelDistance             float64 `json:"total_travel_distance"`
	TotalTravelTime                 int     `json:"total_travel_time"`
	TotalEarlinessPenalty           *int    `json:"total_earliness_penalty,omitempty"`
	TotalLatenessPenalty            *int    `json:"total_lateness_penalty,omitempty"`
	TotalUnassignedPenalty          float64 `json:"total_unassigned_penalty"`
	TotalVehicleInitializationCosts int     `json:"total_vehicle_initialization_costs"`
}

// FleetOutputVehicle contains the route for one vehicle.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [SolutionOutput] and it's
// entities instead.
type FleetOutputVehicle struct {
	ID                         string            `json:"id"`
	Value                      int               `json:"value"`
	TravelDistance             float64           `json:"travel_distance"`
	TravelTime                 int               `json:"travel_time"`
	EarlinessPenalty           *int              `json:"earliness_penalty,omitempty"`
	LatenessPenalty            *int              `json:"lateness_penalty,omitempty"`
	VehicleInitializationCosts int               `json:"vehicle_initialization_costs"`
	Route                      []FleetOutputStop `json:"route,omitempty"`
	Polyline                   string            `json:"polyline,omitempty"`
}

// FleetOutputStop contains all data of a stop.
// DEPRECATION NOTICE: this part of the API is deprecated and is no longer
// maintained. It will be deleted soon. Please use [SolutionOutput] and it's
// entities instead.
type FleetOutputStop struct {
	ID               string     `json:"id"`
	Lon              float64    `json:"lon"`
	Lat              float64    `json:"lat"`
	Distance         *float64   `json:"distance,omitempty"`
	ETA              *time.Time `json:"eta,omitempty"`
	ETS              *time.Time `json:"ets,omitempty"`
	ETD              *time.Time `json:"etd,omitempty"`
	EarlinessPenalty *int       `json:"earliness_penalty,omitempty"`
	LatenessPenalty  *int       `json:"lateness_penalty,omitempty"`
	Polyline         string     `json:"polyline,omitempty"`
}

func format(
	ctx context.Context,
	duration time.Duration,
	progressioner alns.Progressioner,
	toSolutionOutputFn func(nextroute.Solution) (any, error),
	solutions ...nextroute.Solution,
) (FleetOutput, error) {
	mappedSolutions, err := mapWithError(solutions, toSolutionOutputFn)
	if err != nil {
		return FleetOutput{}, err
	}

	output := FleetOutput{
		Solutions: mappedSolutions,
		Options: schema.Options{
			Solver: &schema.SolverOptions{
				Limits: &schema.Limits{
					Duration: duration.String(),
				},
			},
		},
		Hop: struct {
			Version string "json:\"version\""
		}{
			Version: sdk.VERSION,
		},
	}

	startTime := time.Time{}
	if start, ok := ctx.Value(run.Start).(time.Time); ok {
		startTime = start
	}

	progressionValues := progressioner.Progression()

	if len(progressionValues) == 0 {
		return output, errors.New("no solution values or elapsed time values found")
	}

	seriesData := common.Map(
		progressionValues,
		func(progressionEntry alns.ProgressionEntry) statistics.DataPoint {
			return statistics.DataPoint{
				X: statistics.Float64(progressionEntry.ElapsedSeconds),
				Y: statistics.Float64(progressionEntry.Value),
			}
		},
	)

	if len(output.Solutions) == 1 {
		seriesData = seriesData[len(seriesData)-1:]
	}

	if len(output.Solutions) != len(seriesData) {
		return output, errors.New("more or less solution values than solutions found")
	}
	for idx, data := range seriesData {
		if _, ok := output.Solutions[idx].(FleetState); ok {
			output.Statistics.Time.Start = startTime
			output.Statistics.Value = int(data.Y)
			output.Statistics.Time.ElapsedSeconds = float64(data.X)
			output.Statistics.Time.Elapsed = time.Duration(data.X * statistics.Float64(time.Second)).String()
		}
	}

	return output, nil
}

// ToFleetSolutionOutput is a transformation function to create a legacy fleet
// output format.
func ToFleetSolutionOutput(solution nextroute.Solution) (any, error) {
	fleetOutputVehicles, err := mapWithError(solution.Vehicles(), toFleetVehicleOutput)
	if err != nil {
		return FleetState{}, err
	}

	return FleetState{
		Unassigned: common.MapSlice(
			solution.UnPlannedPlanUnits().SolutionPlanUnits(),
			func(solutionPlanUnit nextroute.SolutionPlanUnit) []FleetOutputStop {
				if solutionPlanStopsUnit, ok := solutionPlanUnit.(nextroute.SolutionPlanStopsUnit); ok {
					return common.Map(
						solutionPlanStopsUnit.SolutionStops(),
						func(solutionStop nextroute.SolutionStop) FleetOutputStop {
							return toFleetStopOutput(solutionStop.ModelStop())
						},
					)
				}
				return []FleetOutputStop{}
			},
		),
		Vehicles: fleetOutputVehicles,
		Summary:  toFleetObjectiveOutput(solution),
	}, nil
}

func toFleetStopOutput(modelStop nextroute.ModelStop) FleetOutputStop {
	return FleetOutputStop{
		ID:  modelStop.ID(),
		Lon: modelStop.Location().Longitude(),
		Lat: modelStop.Location().Latitude(),
	}
}

func toPlannedStopOutput(solutionStop nextroute.SolutionStop) FleetOutputStop {
	timezoneLocation := solutionStop.
		Vehicle().
		ModelVehicle().
		Start().
		Location()

	stop := toFleetStopOutput(solutionStop.ModelStop())

	plannedStopOutput := FleetOutputStop{
		ID:  stop.ID,
		Lon: stop.Lon,
		Lat: stop.Lat,
	}

	arrival := solutionStop.Arrival().In(timezoneLocation)
	end := solutionStop.End().In(timezoneLocation)
	start := solutionStop.Start().In(timezoneLocation)

	if solutionStop.Vehicle().First().Start() !=
		solutionStop.Vehicle().ModelVehicle().Model().Epoch() {
		plannedStopOutput.ETA = &arrival
		plannedStopOutput.ETD = &end
		plannedStopOutput.ETS = &start
	}

	initPenalty := 0
	plannedStopOutput.EarlinessPenalty = &initPenalty
	plannedStopOutput.LatenessPenalty = &initPenalty
	if inputStop, ok := solutionStop.ModelStop().Data().(schema.Stop); ok {
		if inputStop.EarlyArrivalTimePenalty != nil && inputStop.TargetArrivalTime != nil {
			penaltyValue := int(*inputStop.EarlyArrivalTimePenalty)
			duration := int(math.Max(inputStop.TargetArrivalTime.Sub(arrival).Seconds(), 0.0))
			penalty := penaltyValue * duration
			plannedStopOutput.EarlinessPenalty = &penalty
		}

		if inputStop.LateArrivalTimePenalty != nil && inputStop.TargetArrivalTime != nil {
			penaltyValue := int(*inputStop.LateArrivalTimePenalty)
			duration := int(math.Max(arrival.Sub(*inputStop.TargetArrivalTime).Seconds(), 0.0))
			penalty := penaltyValue * duration
			plannedStopOutput.LatenessPenalty = &penalty
		}
	}

	if data, ok := solutionStop.Vehicle().ModelVehicle().VehicleType().Data().(distanceData); ok {
		distance := 0.0
		if solutionStop.Previous().ModelStop().Location().IsValid() &&
			solutionStop.ModelStop().Location().IsValid() {
			distance = data.distance.Value(
				solutionStop.Vehicle().ModelVehicle().VehicleType(),
				solutionStop.Previous().ModelStop(),
				solutionStop.ModelStop(),
			)
		}
		plannedStopOutput.Distance = &distance
	}

	return plannedStopOutput
}

func toFleetVehicleOutput(vehicle nextroute.SolutionVehicle) (FleetOutputVehicle, error) {
	solutionStops := common.Filter(
		vehicle.SolutionStops(),
		func(solutionStop nextroute.SolutionStop) bool {
			return solutionStop.ModelStop().Location().IsValid()
		},
	)

	route := common.Map(
		solutionStops,
		toPlannedStopOutput,
	)

	earliness := 0
	lateness := 0
	routeTravelDistance := 0.0
	polylinePoints := []measure.Point{}
	for idx, stop := range route {
		if stop.EarlinessPenalty != nil {
			earliness += *stop.EarlinessPenalty
		}
		if stop.LatenessPenalty != nil {
			lateness += *stop.LatenessPenalty
		}
		routeTravelDistance += *stop.Distance
		route[idx].Distance = &routeTravelDistance

		polylinePoints = append(polylinePoints, measure.Point{
			stop.Lon, stop.Lat,
		})
	}

	polyline, legs := haversinePolyline(polylinePoints)

	for i := range route {
		if i < len(legs) {
			route[i].Polyline = legs[i]
		}
	}

	initCosts := 0
	if inputVehicle, ok := vehicle.ModelVehicle().Data().(schema.Vehicle); ok {
		if len(route) > 0 && inputVehicle.ActivationPenalty != nil {
			initCosts = *inputVehicle.ActivationPenalty
		}
	}

	vehicleOutput := FleetOutputVehicle{
		ID:                         vehicle.ModelVehicle().ID(),
		Route:                      route,
		TravelTime:                 int(vehicle.Duration().Seconds()),
		TravelDistance:             routeTravelDistance,
		Value:                      int(vehicle.Duration().Seconds()) + earliness + lateness,
		EarlinessPenalty:           &earliness,
		LatenessPenalty:            &lateness,
		VehicleInitializationCosts: initCosts,
		Polyline:                   polyline,
	}

	return vehicleOutput, nil
}

func toFleetObjectiveOutput(solution nextroute.Solution) FleetOutputSummary {
	earlinessPenalty := 0
	latenessPenatly := 0
	unassignedPenalty := 0.0
	travelTime := 0.0
	initCosts := 0

	for _, t := range solution.Model().Objective().Terms() {
		switch fmt.Sprintf("%v", t.Objective()) {
		case "early_arrival_penalty":
			earlinessPenalty = int(solution.ObjectiveValue(t.Objective()))
		case "unplanned_penalty":
			unassignedPenalty = solution.ObjectiveValue(t.Objective())
		case "late_arrival_penalty":
			latenessPenatly = int(solution.ObjectiveValue(t.Objective()))
		case "vehicles_duration":
			travelTime = solution.ObjectiveValue(t.Objective())
		case "vehicle_activation_penalty":
			initCosts = int(solution.ObjectiveValue(t.Objective()))
		}
	}

	travelDistance := 0.0
	for _, v := range solution.Vehicles() {
		if data, ok := v.ModelVehicle().VehicleType().Data().(distanceData); ok {
			for _, s := range v.SolutionStops() {
				distance := data.distance.Value(
					s.Vehicle().ModelVehicle().VehicleType(),
					s.Previous().ModelStop(),
					s.ModelStop(),
				)
				travelDistance += distance
			}
		}
	}

	return FleetOutputSummary{
		Value:                           solution.ObjectiveValue(solution.Model().Objective()),
		TotalTravelTime:                 int(travelTime),
		TotalEarlinessPenalty:           &earlinessPenalty,
		TotalLatenessPenalty:            &latenessPenatly,
		TotalUnassignedPenalty:          unassignedPenalty,
		TotalTravelDistance:             travelDistance,
		TotalVehicleInitializationCosts: initCosts,
	}
}
