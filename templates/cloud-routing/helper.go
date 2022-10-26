package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/store"
)

func (i *input) applyVehicleDefaults() {
	// Specify some defaults for missing values
	capacity := 0
	vehicleCompatibilities := []string{}
	maxStops := -1
	maxDistance := model.MaxInt
	maxDuration := model.MaxInt

	// Apply all vehicle defaults, if no explicit value is set
	// Apply general default, if neither explicit nor default value is given
	for v := range i.Vehicles {
		vehicleDefaults := i.Defaults != nil && i.Defaults.Vehicles != nil
		if i.Vehicles[v].Start == nil {
			if vehicleDefaults && i.Defaults.Vehicles.Start != nil {
				i.Vehicles[v].Start = i.Defaults.Vehicles.Start
			}
		}

		if i.Vehicles[v].End == nil {
			if vehicleDefaults && i.Defaults.Vehicles.End != nil {
				i.Vehicles[v].End = i.Defaults.Vehicles.End
			}
		}

		if i.Vehicles[v].Speed == nil {
			if vehicleDefaults && i.Defaults.Vehicles.Speed != nil {
				i.Vehicles[v].Speed = i.Defaults.Vehicles.Speed
			} else {
				s := 10.0
				i.Vehicles[v].Speed = &(s)
			}
		}

		if i.Vehicles[v].Capacity == nil {
			if vehicleDefaults && i.Defaults.Vehicles.Capacity != nil {
				i.Vehicles[v].Capacity = i.Defaults.Vehicles.Capacity
			} else {
				i.Vehicles[v].Capacity = &capacity
			}
		}

		if i.Vehicles[v].ShiftStart == nil {
			if vehicleDefaults && i.Defaults.Vehicles.ShiftStart != nil {
				i.Vehicles[v].ShiftStart = i.Defaults.Vehicles.ShiftStart
			}
		}

		if i.Vehicles[v].ShiftEnd == nil {
			if vehicleDefaults && i.Defaults.Vehicles.ShiftEnd != nil {
				i.Vehicles[v].ShiftEnd = i.Defaults.Vehicles.ShiftEnd
			}
		}

		if i.Vehicles[v].CompatibilityAttributes == nil {
			if vehicleDefaults && i.Defaults.Vehicles.CompatibilityAttributes != nil {
				ca := i.Defaults.Vehicles.CompatibilityAttributes
				i.Vehicles[v].CompatibilityAttributes = ca
			} else {
				i.Vehicles[v].CompatibilityAttributes = vehicleCompatibilities
			}
		}

		if i.Vehicles[v].MaxStops == nil {
			if vehicleDefaults && i.Defaults.Vehicles.MaxStops != nil {
				i.Vehicles[v].MaxStops = i.Defaults.Vehicles.MaxStops
			} else {
				i.Vehicles[v].MaxStops = &maxStops
			}
		}

		if i.Vehicles[v].MaxDistance == nil {
			if vehicleDefaults && i.Defaults.Vehicles.MaxDistance != nil {
				i.Vehicles[v].MaxDistance = i.Defaults.Vehicles.MaxDistance
			} else {
				i.Vehicles[v].MaxDistance = &maxDistance
			}
		}

		if i.Vehicles[v].MaxDuration == nil {
			if vehicleDefaults && i.Defaults.Vehicles.MaxDuration != nil {
				i.Vehicles[v].MaxDuration = i.Defaults.Vehicles.MaxDuration
			} else {
				i.Vehicles[v].MaxDuration = &maxDuration
			}
		}

		if i.Vehicles[v].StopDurationMultiplier == nil {
			multiplier := 1.0
			i.Vehicles[v].StopDurationMultiplier = &multiplier
		}
	}
}

// applyStopDefaults applies the given default values for all values of vehicles
// and stops not explicitly defined.
func (i *input) applyStopDefaults() {
	// Specify some defaults for missing values
	unassignedPenalty := 0
	quantity := 0
	hardWindow := []time.Time{}
	maxWait := -1
	stopDuration := 0
	stopCompatibilities := []string{}

	// Use default earliness/lateness penalties only if any penalty is given
	earlinessPenalty, latenessPenalty := 0.0, 0.0
	earlinessUsed, latenessUsed := false, false
	for _, stop := range i.Stops {
		if stop.EarlinessPenalty != nil {
			earlinessUsed = true
		}
		if stop.LatenessPenalty != nil {
			latenessUsed = true
		}
	}

	// Apply all vehicle defaults, if no explicit value is set
	// Apply general default, if neither explicit nor default value is given
	for s := range i.Stops {
		stopDefaults := i.Defaults != nil && i.Defaults.Stops != nil

		if i.Stops[s].UnassignedPenalty == nil {
			if stopDefaults && i.Defaults.Stops.UnassignedPenalty != nil {
				i.Stops[s].UnassignedPenalty = i.Defaults.Stops.UnassignedPenalty
			} else {
				i.Stops[s].UnassignedPenalty = &unassignedPenalty
			}
		}

		if i.Stops[s].Quantity == nil {
			if stopDefaults && i.Defaults.Stops.Quantity != nil {
				i.Stops[s].Quantity = i.Defaults.Stops.Quantity
			} else {
				i.Stops[s].Quantity = &quantity
			}
		}

		if i.Stops[s].HardWindow == nil {
			if stopDefaults && i.Defaults.Stops.HardWindow != nil {
				i.Stops[s].HardWindow = i.Defaults.Stops.HardWindow
			} else {
				i.Stops[s].HardWindow = &hardWindow
			}
		}

		if i.Stops[s].MaxWait == nil {
			if stopDefaults && i.Defaults.Stops.MaxWait != nil {
				i.Stops[s].MaxWait = i.Defaults.Stops.MaxWait
			} else {
				i.Stops[s].MaxWait = &maxWait
			}
		}

		if i.Stops[s].StopDuration == nil {
			if stopDefaults && i.Defaults.Stops.StopDuration != nil {
				i.Stops[s].StopDuration = i.Defaults.Stops.StopDuration
			} else {
				i.Stops[s].StopDuration = &stopDuration
			}
		}

		if i.Stops[s].TargetTime == nil {
			if stopDefaults && i.Defaults.Stops.TargetTime != nil {
				i.Stops[s].TargetTime = i.Defaults.Stops.TargetTime
			}
		}

		if i.Stops[s].EarlinessPenalty == nil {
			if stopDefaults && i.Defaults.Stops.EarlinessPenalty != nil {
				i.Stops[s].EarlinessPenalty = i.Defaults.Stops.EarlinessPenalty
			} else if earlinessUsed {
				i.Stops[s].EarlinessPenalty = &earlinessPenalty
			}
		}

		if i.Stops[s].LatenessPenalty == nil {
			if stopDefaults && i.Defaults.Stops.LatenessPenalty != nil {
				i.Stops[s].LatenessPenalty = i.Defaults.Stops.LatenessPenalty
			} else if latenessUsed {
				i.Stops[s].LatenessPenalty = &latenessPenalty
			}
		}

		if i.Stops[s].CompatibilityAttributes == nil {
			if stopDefaults && i.Defaults.Stops.CompatibilityAttributes != nil {
				ca := i.Defaults.Stops.CompatibilityAttributes
				i.Stops[s].CompatibilityAttributes = ca
			} else {
				i.Stops[s].CompatibilityAttributes = &stopCompatibilities
			}
		}
	}
}

// applySolveOptions applies all solver options given in the input data to the
// given solver options structure.
func (i *input) applySolveOptions(o store.Options) store.Options {
	// Diagram
	if i.Options == nil || i.Options.Solver == nil {
		return o
	}
	diagram := i.Options.Solver.Diagram
	if diagram != nil && diagram.Width != nil {
		o.Diagram.Width = *i.Options.Solver.Diagram.Width
	}
	if diagram != nil && diagram.Expansion != nil &&
		diagram.Expansion.Limit != nil {
		o.Diagram.Expansion.Limit = *diagram.Expansion.Limit
	}
	// Limits
	limits := i.Options.Solver.Limits
	if limits != nil && limits.Duration != nil {
		o.Limits.Duration = i.Options.Solver.Limits.Duration.Duration
	}
	if limits != nil && limits.Nodes != nil {
		o.Limits.Nodes = *i.Options.Solver.Limits.Nodes
	}
	if limits != nil && limits.Solutions != nil {
		o.Limits.Solutions = *i.Options.Solver.Limits.Solutions
	}
	return o
}

// dynamicDefaultKey is used when only scalar quantities / capacities are given.
const dynamicDefaultKey = "default"

// handleDynamics converts all dynamic input fields (fields with not only one
// explicitly expected type). The values will be written to un-exported fields.
func (i *input) handleDynamics() error {
	for s, stop := range i.Stops {
		// Handle multi capacity fields
		if reflect.ValueOf(stop.Quantity).Kind() == reflect.Map {
			// We detected a map and need to handle it
			quantities, ok := stop.Quantity.(map[string]any)
			if !ok {
				return fmt.Errorf(
					"expecting string to int map for quantity at stop %s, but got %v",
					stop.ID,
					stop.Quantity,
				)
			}
			// Init map
			i.Stops[s].quantity = map[string]int{}
			// Set all values
			for kind, quant := range quantities {
				rounded, ok := convertToInt(quant)
				if !ok {
					return fmt.Errorf(
						"expecting integer value for quantity at stop %s, but got %s: %v",
						stop.ID,
						kind,
						quant,
					)
				}
				i.Stops[s].quantity[kind] = rounded
			}
		} else {
			// Since it is not a map, we expect a scalar value
			quantity, ok := convertToInt(stop.Quantity)
			if !ok {
				return fmt.Errorf(
					"expecting string to int map for quantity at stop %s, but got %v",
					stop.ID,
					stop.Quantity,
				)
			}
			i.Stops[s].quantity = map[string]int{
				dynamicDefaultKey: quantity,
			}
		}
		// Handle multi precedes fields
		if stop.Precedes == nil {
			i.Stops[s].precedes = []string{}
		} else if reflect.ValueOf(stop.Precedes).Kind() == reflect.Slice {
			slice := stop.Precedes.([]any)
			precedes := make([]string, len(slice))
			for i, p := range slice {
				val, ok := p.(string)
				if !ok {
					return fmt.Errorf(
						"expecting slice of string for precedes at stop %s, but got %v",
						stop.ID,
						stop.Precedes,
					)
				}
				precedes[i] = val
			}
			i.Stops[s].precedes = precedes
		} else {
			// Since it is not an array, we expect a string
			precedes, ok := stop.Precedes.(string)
			if !ok {
				return fmt.Errorf(
					"expecting string for precedes at stop %s, but got %v",
					stop.ID,
					stop.Precedes,
				)
			}
			i.Stops[s].precedes = []string{precedes}
		}
		// Handle multi succeeds fields
		if stop.Succeeds == nil {
			i.Stops[s].succeeds = []string{}
		} else if reflect.ValueOf(stop.Succeeds).Kind() == reflect.Slice {
			slice := stop.Succeeds.([]any)
			succeeds := make([]string, len(slice))
			for i, p := range slice {
				val, ok := p.(string)
				if !ok {
					return fmt.Errorf(
						"expecting slice of string for succeeds at stop %s, but got %v",
						stop.ID,
						stop.Succeeds,
					)
				}
				succeeds[i] = val
			}
			i.Stops[s].succeeds = succeeds
		} else {
			// Since it is not an array, we expect a string
			succeeds, ok := stop.Succeeds.(string)
			if !ok {
				return fmt.Errorf("expecting string for succeeds at stop %s, but got %v",
					stop.ID,
					stop.Succeeds)
			}
			i.Stops[s].succeeds = []string{succeeds}
		}
	}
	for v, vehicle := range i.Vehicles {
		if reflect.ValueOf(vehicle.Capacity).Kind() == reflect.Map {
			// We detected a map and need to handle it
			capacities, ok := vehicle.Capacity.(map[string]any)
			if !ok {
				return fmt.Errorf(
					"expecting string to int map for capacity at vehicle %s, but got %v",
					vehicle.ID,
					vehicle.Capacity,
				)
			}
			// Init map
			i.Vehicles[v].capacity = map[string]int{}
			// Set all values
			for kind, cap := range capacities {
				rounded, ok := convertToInt(cap)
				if !ok {
					return fmt.Errorf(
						"expecting integer value for capacity at vehicle %s, but got %s: %v",
						vehicle.ID,
						kind,
						cap,
					)
				}
				i.Vehicles[v].capacity[kind] = rounded
			}
		} else {
			// Since it is not a map, we expect a scalar value
			capacity, ok := convertToInt(vehicle.Capacity)
			if !ok {
				return fmt.Errorf(
					"expecting string to int map for capacity at vehicle %s, but got %v",
					vehicle.ID,
					vehicle.Capacity,
				)
			}
			i.Vehicles[v].capacity = map[string]int{
				dynamicDefaultKey: capacity,
			}
		}
	}

	// Set all dynamic fields - we are done here
	return nil
}

func (i *input) makeDurationGroups() error {
	durationGroupDomains := make([]route.DurationGroup, len(i.DurationGroups))
	stopIDToIndex := make(map[string]int, len(i.Stops))
	for s, stop := range i.Stops {
		if stop.ID == "" {
			return errors.New("stop has an empty ID")
		}
		stopIDToIndex[stop.ID] = s
	}

	for idx, dg := range i.DurationGroups {
		group := model.NewDomain()
		for _, stop := range dg.Group {
			if val, ok := stopIDToIndex[stop]; ok {
				group = group.Add(val)
			} else {
				return fmt.Errorf(
					"could not find the stop %s in duration group %d",
					stop,
					idx,
				)
			}
		}
		durationGroupDomains[idx] = route.DurationGroup{
			Duration: dg.Duration,
			Group:    group,
		}
	}
	i.durationGroupsDomains = durationGroupDomains
	return nil
}

// intTolerance specified the absolute offset tolerable for dynamic fields that
// are naturally float, but where integer values are expected.
const intTolerance = 0.00001

// convertToInt tries to convert the given unknown typed value to int. Returns
// true, iff the conversion was successful.
func convertToInt(unknown any) (int, bool) {
	// Convert unknown to float
	floatType := reflect.TypeOf(float64(0))
	v := reflect.ValueOf(unknown)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, false
	}
	fv := v.Convert(floatType)
	floatValue := fv.Float()
	// Round to next int (not ok if offset out of tolerance)
	rounded := math.Round(floatValue)
	if math.Abs(rounded-floatValue) > intTolerance {
		return 0, false
	}
	return int(rounded), true
}

const autoPenaltySlack = 0.01

// autoConfigureUnassigned sets automatically configured unassigned penalties,
// if all unassigned penalties are 0. Returns a bool indicating whether
// penalties were auto-configured.
func (i *input) autoConfigureUnassigned() bool {
	// Check whether we need to auto-configure
	for _, stop := range i.Stops {
		if *stop.UnassignedPenalty != 0 {
			return false
		}
	}
	// >> Determine suitable penalty
	// Big-m as per stop fraction of max int
	// Assumption: costs for servicing a stop never exceed the max int fraction
	// Caveats:
	//  - Needs to be revised, if introducing non-stop-related costs
	//  - Homogeneous penalty per stop is larger for travel time than individual
	//    would be
	//    - Should still be fine regarding overflows
	// Advantages:
	//  - Much faster than computing all stop to stop costs for all vehicles
	//  - Much simpler than correctly respecting all features
	//    (e.g.: lateness/earliness)
	autoPenalty := 0
	if len(i.Stops) > 0 {
		// each stop gets homogeneous fraction of MaxInt
		autoPenalty = int(float64(model.MaxInt) / float64(len(i.Stops)) *
			// reserve some
			(1 - autoPenaltySlack))
	}

	// Set auto penalty
	for s := range i.Stops {
		i.Stops[s].UnassignedPenalty = &autoPenalty
	}
	return true
}

// buildHaversineMeasures returns an array of Haversine measures based on stops
// and vehicles' start and end locations.
func buildHaversineMeasures(
	vehicles []vehicle,
	stops []stop,
	starts []route.Position,
	ends []route.Position,
) ([]route.ByIndex, []route.Point) {
	count := len(stops)
	points := make([]route.Point, count+2*len(vehicles))

	for s, stop := range stops {
		points[s] = point(stop.Position)
	}

	if len(starts) > 0 {
		for v, start := range starts {
			points[count+v*2] = point(start)
		}
	}

	if len(ends) > 0 {
		for v, end := range ends {
			points[count+v*2+1] = point(end)
		}
	}

	// Haversine measure and override cost of going to/from an empty point.
	m := route.Indexed(route.HaversineByPoint(), points)
	m = route.Override(
		m,
		route.Constant(0),
		func(from, to int) bool {
			return points[from] == nil || points[to] == nil
		},
	)
	measures := make([]route.ByIndex, len(vehicles))
	for v := range vehicles {
		measures[v] = m
	}

	return measures, points
}

// buildTravelTimeMeasures returns an array of duration measures with a velocity
// of 10m/s if neither Velocities option nor TravelTimeMeasure was not used.
func buildTravelTimeMeasures(
	velocities []float64,
	vehicles []vehicle,
	stops []stop,
	starts []route.Position,
	ends []route.Position,
	durationGroups route.DurationGroups,
) ([]route.ByIndex, error) {
	timeMeasures := make([]route.ByIndex, len(vehicles))
	measures, points := buildHaversineMeasures(vehicles, stops, starts, ends)
	for i := range measures {
		timeMeasurer := route.Scale(measures[i], 1/velocities[i])
		timeMeasures[i] = timeMeasurer
	}

	timeMeasures, err := addLocationMeasure(
		timeMeasures,
		points,
		vehicles,
		stops, durationGroups,
	)
	if err != nil {
		return nil, err
	}

	return timeMeasures, nil
}

// addLocationMeasure creates a measure for each vehicle containing travel time
// and stop durations.
func addLocationMeasure(
	ms []route.ByIndex,
	points []route.Point,
	vehicles []vehicle,
	stops []stop,
	durationGroupsDomains route.DurationGroups,
) ([]route.ByIndex, error) {
	newMs := make([]route.ByIndex, len(ms))
	copy(newMs, ms)

	// Get durations of stops
	anyDuration := false
	vehicleDurations := make([][]float64, len(vehicles))
	if len(durationGroupsDomains) > 0 {
		anyDuration = true
	}
	for d := range vehicleDurations {
		durations := make([]float64, len(points))
		for s, stop := range stops {
			if stop.StopDuration != nil && *stop.StopDuration != 0 {
				anyDuration = true
				multiplier := vehicles[d].StopDurationMultiplier
				durations[s] = float64(*stop.StopDuration) * *multiplier
			}
		}
		vehicleDurations[d] = durations
	}

	// Add stop durations
	if anyDuration {
		for v := range vehicles {
			newMeasure, err := route.Location(
				ms[v], vehicleDurations[v],
				durationGroupsDomains,
			)
			if err != nil {
				return nil, err
			}
			newMs[v] = newMeasure
		}
	}

	return newMs, nil
}

func point(p route.Position) route.Point {
	if p.Lon == 0 || p.Lat == 0 {
		return nil
	}
	return route.Point{
		p.Lon, p.Lat,
	}
}
