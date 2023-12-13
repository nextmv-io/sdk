// Package decode implements a parser for TSPLIB instances.
package decode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run/decode"
)

// TSPLIB creates a TSPLIB decoder.
func TSPLIB() decode.Decoder {
	return TSPLIBDecoder{}
}

// TSPLIBDecoder is a Decoder that decodes a json into a struct.
type TSPLIBDecoder struct {
}

const (
	inUndefined         = iota
	inCoords            = iota
	inEdgeWeight        = iota
	inDemand            = iota
	inDepot             = iota
	inPickupAndDelivery = iota
)

// Decode a tsplib instance from the given reader to the nextroute input format.
func (j TSPLIBDecoder) Decode(reader io.Reader, anyInput any) error {
	input, ok := anyInput.(*schema.Input)
	if !ok {
		return fmt.Errorf("input is not of type schema.Input")
	}
	stopDefaults := schema.StopDefaults{}
	vehicleDefaults := schema.VehicleDefaults{}
	defaults := schema.Defaults{
		Stops:    &stopDefaults,
		Vehicles: &vehicleDefaults,
	}
	input.Defaults = &defaults

	numberOfVehicles := 0
	// Prepare
	scanner := bufio.NewScanner(reader)
	var byPoint measure.ByPoint
	var matrix *[][]float64

	// Parse file
	numbers := []int{}
	stopsByNumber := map[int]schema.Stop{}
	// these are used to scale the coordinates to the range -180, 180 and -90,
	// 90 because the TSPLIB format allow arbitrary coordinates but nextroute
	// only supports coordinates in that range.
	maxLon, maxLat := 0.0, 0.0
	depots := map[int]struct{}{}
	section := inUndefined
	explicit := false
	measureGiven := false
	matrix = &[][]float64{}
scanLoop: // Label scanner loop to break out of nested switch statements
	for scanner.Scan() {
		// Get line and make sure there is whitespace around the delimiter
		line := strings.Replace(scanner.Text(), ":", " : ", 1)
		// Split line by whitespace
		s := strings.Fields(line)
		if len(s) < 1 {
			continue
		}

		// Check for new section start
		switch s[0] {
		case "VEHICLES":
			nrOfVehicles, err := strconv.Atoi(s[2])
			if err != nil {
				return errors.New("error parsing number of vehicles: " + err.Error())
			}
			numberOfVehicles = nrOfVehicles
		case "CAPACITY":
			capacity, err := strconv.Atoi(s[2])
			if err != nil {
				return errors.New("error parsing capacity: " + err.Error())
			}
			input.Defaults.Vehicles.Capacity = capacity
			continue
		case "DISTANCE":
			// parse float
			floatDuration, err := strconv.ParseFloat(s[2], 64)
			if err != nil {
				return errors.New("error parsing distance: " + err.Error())
			}
			maxDuration := int(floatDuration)
			input.Defaults.Vehicles.MaxDuration = &maxDuration
			continue
		case "SERVICE_TIME":
			serviceTime, err := strconv.Atoi(s[2])
			if err != nil {
				return errors.New("error parsing service time: " + err.Error())
			}
			input.Defaults.Stops.Duration = &serviceTime
			continue
		case "EDGE_WEIGHT_TYPE":
			measureGiven = true
			switch s[2] {
			case "EXPLICIT":
				explicit = true
			case "MAN_2D":
				byPoint = measure.TaxicabByPoint()
			case "EUC_2D":
				byPoint = measure.EuclideanByPoint()
			case "GEO":
				byPoint = measure.HaversineByPoint()
			default:
				return fmt.Errorf("unsupported measure: %s", s[2])
			}
			continue
		case "EDGE_WEIGHT_FORMAT":
			if s[2] != "FULL_MATRIX" {
				return fmt.Errorf("unsupported edge weight format: %s", s[2])
			}
			continue
		case "NODE_COORD_SECTION":
			// Start of node coordinate section
			section = inCoords
			continue
		case "EDGE_WEIGHT_SECTION":
			// Start of edge weight section (alternative to node coords)
			section = inEdgeWeight
			continue
		case "DEMAND_SECTION":
			// Start of section defining demands per node
			section = inDemand
			continue
		case "DEPOT_SECTION":
			// Start of section defining the depots
			section = inDepot
			continue
		case "PICKUP_AND_DELIVERY_SECTION":
			// Start of section defining pickup and delivery requests
			section = inPickupAndDelivery
			continue
		case "END", "EOF":
			// End of file
			break scanLoop
		}

		// Skip, if in unknown territory
		if section == inUndefined {
			continue
		}

		switch section {
		case inCoords:
			// Get the number of the associated request
			number, err := strconv.Atoi(s[0])
			if err != nil {
				return errors.New("error parsing number in NODE_COORD_SECTION: " + err.Error())
			}
			// Create a request for the coordinates
			x, err := strconv.ParseFloat(s[1], 64)
			if err != nil {
				return errors.New("error parsing x coordinate in NODE_COORD_SECTION: " + err.Error())
			}
			y, err := strconv.ParseFloat(s[2], 64)
			if err != nil {
				return errors.New("error parsing y coordinate in NODE_COORD_SECTION: " + err.Error())
			}
			stopsByNumber[number] = schema.Stop{
				Location: schema.Location{
					Lon: x,
					Lat: y,
				},
				ID: fmt.Sprint(number),
			}
			numbers = append(numbers, number)
			xAbs := math.Abs(x)
			if xAbs > maxLon {
				maxLon = xAbs
			}
			yAbs := math.Abs(y)
			if yAbs > maxLat {
				maxLat = yAbs
			}
		case inEdgeWeight:
			// Create requests from edges
			row := make([]float64, len(s))
			for index, v := range s {
				weight, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return errors.New("error parsing edge weight in EDGE_WEIGHT_SECTION: " + err.Error())
				}
				row[index] = weight
			}
			*matrix = append(*matrix, row)
		case inDemand:
			// Get the number of the associated request
			number, err := strconv.Atoi(s[0])
			if err != nil {
				return errors.New("error parsing number in DEMAND_SECTION: " + err.Error())
			}
			// Set demand for request
			l := stopsByNumber[number]
			quantity, err := strconv.Atoi(s[1])
			if err != nil {
				return errors.New("error parsing demand in DEMAND_SECTION: " + err.Error())
			}
			l.Quantity = -quantity // negative quantity indicates demand
			stopsByNumber[number] = l
			// If demand is zero, it must be a depot - mark it
			if l.Quantity == 0 {
				depots[number] = struct{}{}
			}
		case inDepot:
			// Get the number of the associated request
			number, err := strconv.Atoi(s[0])
			if err != nil {
				return errors.New("error parsing number in DEPOT_SECTION: " + err.Error())
			}
			// Skip section terminal
			if number == -1 {
				continue
			}
			// Check whether depot was already marked correctly by demand
			// section
			if _, ok := depots[number]; !ok {
				// Skip depot
				err := fmt.Errorf(
					"error parsing instance: depot %d in DEPOT section was not represented with 0 demand in DEMAND section", number)
				return err
			}
		case inPickupAndDelivery:
			return fmt.Errorf("pickup and delivery requests are not supported")
		}
	}

	// Measure must be defined at this point
	if !measureGiven {
		return fmt.Errorf("no measure definition provided")
	}
	// If explicit was defined but no distances were given, that is an error
	if explicit && matrix == nil {
		return fmt.Errorf("explicit edge weights indicated but not found")
	}

	if numberOfVehicles == 0 {
		return fmt.Errorf("no vehicles defined")
	}

	// Determine depot 'request' (choose min index, if multiple are given)
	depotNumber := math.MaxInt32
	for n := range depots {
		if n < depotNumber {
			depotNumber = n
		}
	}
	// Use first location, if no depots are given
	if depotNumber == math.MaxInt32 {
		for _, n := range numbers {
			if n < depotNumber {
				depotNumber = n
			}
		}
	}

	// Use 1 as depot number, if no locations are given
	if depotNumber == math.MaxInt32 {
		depotNumber = 1
	}
	depot := stopsByNumber[depotNumber]

	// Aggregate stops
	if len(numbers) != 0 && len(stopsByNumber) == len(numbers) {
		for _, number := range numbers {
			stop := stopsByNumber[number]
			if _, ok := depots[number]; ok {
				// Skip depot locations
				continue
			}
			// Add request
			input.Stops = append(input.Stops, stop)
		}
	} else {
		// the matrix defines how many stops there are
		for i := 0; i < len(*matrix); i++ {
			// matrix is 0 indexed, but stops are 1 indexed
			if i+1 == depotNumber {
				continue
			}
			stop := schema.Stop{
				ID: fmt.Sprint(i + 1),
				// location 1,1 is a dummy location
				Location: schema.Location{
					Lat: 1,
					Lon: 1,
				},
			}
			if stopByNumber, ok := stopsByNumber[i+1]; ok {
				stop.Quantity = stopByNumber.Quantity
			}
			input.Stops = append(input.Stops, stop)
		}
	}

	// Create estimated / sufficient vehicles
	for i := 0; i < numberOfVehicles; i++ {
		input.Vehicles = append(input.Vehicles, schema.Vehicle{
			ID:            fmt.Sprint(i),
			StartLocation: &depot.Location,
			EndLocation:   &depot.Location,
		})
	}

	size := len(input.Stops) + 2*len(input.Vehicles)
	floats := make([][]float64, size)

	if byPoint != nil {
		for i := 0; i < len(input.Stops); i++ {
			floats[i] = make([]float64, size)
			for j := 0; j < len(input.Stops); j++ {
				floats[i][j] = byPoint.Cost(
					measure.Point{
						input.Stops[i].Location.Lon,
						input.Stops[i].Location.Lat,
					},
					measure.Point{
						input.Stops[j].Location.Lon,
						input.Stops[j].Location.Lat,
					},
				)
			}
			for j := 0; j < 2*len(input.Vehicles); j += 2 {
				floats[i][j+len(input.Stops)] = byPoint.Cost(
					measure.Point{
						input.Stops[i].Location.Lon,
						input.Stops[i].Location.Lat,
					},
					measure.Point{
						input.Vehicles[j/2].StartLocation.Lon,
						input.Vehicles[j/2].StartLocation.Lat,
					},
				)
				floats[i][j+1+len(input.Stops)] = byPoint.Cost(
					measure.Point{
						input.Stops[i].Location.Lon,
						input.Stops[i].Location.Lat,
					},
					measure.Point{
						input.Vehicles[j/2].EndLocation.Lon,
						input.Vehicles[j/2].EndLocation.Lat,
					},
				)
			}
		}
		for i := 0; i < 2*len(input.Vehicles); i += 2 {
			for x := 0; x < 2; x++ {
				if x == 0 {
					floats[i+len(input.Stops)] = make([]float64, size)
					for j := 0; j < len(input.Stops); j++ {
						floats[i+len(input.Stops)][j] = byPoint.Cost(
							measure.Point{
								input.Vehicles[i/2].StartLocation.Lon,
								input.Vehicles[i/2].StartLocation.Lat,
							},
							measure.Point{
								input.Stops[j].Location.Lon,
								input.Stops[j].Location.Lat,
							},
						)
					}
					for j := 0; j < len(input.Vehicles); j++ {
						floats[i+len(input.Stops)][j+len(input.Stops)] = byPoint.Cost(
							measure.Point{
								input.Vehicles[i/2].StartLocation.Lon,
								input.Vehicles[i/2].StartLocation.Lat,
							},
							measure.Point{
								input.Vehicles[j].StartLocation.Lon,
								input.Vehicles[j].StartLocation.Lat,
							},
						)
						floats[i+len(input.Stops)][j+1+len(input.Stops)] = byPoint.Cost(
							measure.Point{
								input.Vehicles[i/2].StartLocation.Lon,
								input.Vehicles[i/2].StartLocation.Lat,
							},
							measure.Point{
								input.Vehicles[j].EndLocation.Lon,
								input.Vehicles[j].EndLocation.Lat,
							},
						)
					}
				} else {
					floats[i+1+len(input.Stops)] = make([]float64, size)
					for j := 0; j < len(input.Stops); j++ {
						floats[i+1+len(input.Stops)][j] = byPoint.Cost(
							measure.Point{
								input.Vehicles[i/2].EndLocation.Lon,
								input.Vehicles[i/2].EndLocation.Lat,
							},
							measure.Point{
								input.Stops[j].Location.Lon,
								input.Stops[j].Location.Lat,
							},
						)
					}
					for j := 0; j < len(input.Vehicles); j++ {
						floats[i+1+len(input.Stops)][j+len(input.Stops)] = byPoint.Cost(
							measure.Point{
								input.Vehicles[i/2].EndLocation.Lon,
								input.Vehicles[i/2].EndLocation.Lat,
							},
							measure.Point{
								input.Vehicles[j].StartLocation.Lon,
								input.Vehicles[j].StartLocation.Lat,
							},
						)
						floats[i+1+len(input.Stops)][j+1+len(input.Stops)] = byPoint.Cost(
							measure.Point{
								input.Vehicles[i/2].EndLocation.Lon,
								input.Vehicles[i/2].EndLocation.Lat,
							},
							measure.Point{
								input.Vehicles[j].EndLocation.Lon,
								input.Vehicles[j].EndLocation.Lat,
							},
						)
					}
				}
			}
		}

		// scale the coordinates to the range -180, 180 and -90, 90
		for i := 0; i < len(input.Stops); i++ {
			stop := input.Stops[i]
			stop.Location.Lon = input.Stops[i].Location.Lon / maxLon * 180
			stop.Location.Lat = input.Stops[i].Location.Lat / maxLat * 90
			input.Stops[i] = stop
		}
		for i := 0; i < len(input.Vehicles); i++ {
			vehicle := input.Vehicles[i]
			vehicle.StartLocation.Lon = input.Vehicles[i].StartLocation.Lon / maxLon * 180
			vehicle.StartLocation.Lat = input.Vehicles[i].StartLocation.Lat / maxLat * 90
			vehicle.EndLocation.Lon = input.Vehicles[i].EndLocation.Lon / maxLon * 180
			vehicle.EndLocation.Lat = input.Vehicles[i].EndLocation.Lat / maxLat * 90
			input.Vehicles[i] = vehicle
		}

		input.DurationMatrix = &floats
		input.DistanceMatrix = &floats

		return nil
	}

	depotRowSeen := false
	depotIndexInMatrix := depotNumber - 1

	for rowIndex := 0; rowIndex <= size; rowIndex++ {
		if rowIndex == depotIndexInMatrix {
			depotRowSeen = true
			continue
		}
		effectiveRowIndex := rowIndex
		if depotRowSeen {
			effectiveRowIndex--
		}
		var row []float64
		if rowIndex < len(*matrix) {
			row = (*matrix)[rowIndex]
		} else {
			row = (*matrix)[depotIndexInMatrix]
		}
		floats[effectiveRowIndex] = make([]float64, size)
		depotColumnSeen := false
		for i := 0; i <= size; i++ {
			if i == depotIndexInMatrix {
				depotColumnSeen = true
				continue
			}
			effectiveI := i
			if depotColumnSeen {
				effectiveI--
			}
			if i < len(row) {
				floats[effectiveRowIndex][effectiveI] = row[i]
			} else {
				floats[effectiveRowIndex][effectiveI] = row[depotIndexInMatrix]
			}
		}
	}

	input.DurationMatrix = &floats
	input.DistanceMatrix = &floats

	return nil
}
