// Package decode implements a parser for TSPLIB instances.
package decode

import (
	"bufio"
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
type TSPLIBDecoder struct{}

const (
	inUndefined  = iota
	inCoords     = iota
	inEdgeWeight = iota
	inDemand     = iota
	inDepot      = iota
)

// Parse a tsplib instance from the given reader to the nextroute input format.
func (j TSPLIBDecoder) Decode(reader io.Reader, anyInput any) error {
	input, ok := anyInput.(*schema.Input)
	if !ok {
		return fmt.Errorf("input is not of type schema.Input")
	}
	// Prepare
	scanner := bufio.NewScanner(reader)
	var byPoint measure.ByPoint
	var matrix *[][]float64

	// Parse file
	var capacity int
	numbers := []int{}
	stopsByNumber := map[int]schema.Stop{}
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
		case "CAPACITY":
			capacity, _ = strconv.Atoi(s[2])
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
			number, _ := strconv.Atoi(s[0])
			// Create a request for the coordinates
			x, _ := strconv.ParseFloat(s[1], 64)
			y, _ := strconv.ParseFloat(s[2], 64)
			stopsByNumber[number] = schema.Stop{
				Location: schema.Location{
					Lat: x,
					Lon: y,
				},
			}
			numbers = append(numbers, number)
		case inEdgeWeight:
			// Create requests from edges
			row := make([]float64, len(s))
			for index, v := range s {
				row[index], _ = strconv.ParseFloat(v, 64)
			}
			*matrix = append(*matrix, row)
		case inDemand:
			// Get the number of the associated request
			number, _ := strconv.Atoi(s[0])
			// Set demand for request
			l := stopsByNumber[number]
			quantity, _ := strconv.Atoi(s[1])
			l.Quantity = -quantity // We need negative quantities as these are pickups
			stopsByNumber[number] = l
			// If demand is zero, it must be a depot - mark it
			if l.Quantity == 0 {
				depots[number] = struct{}{}
			}
		case inDepot:
			// Get the number of the associated request
			number, _ := strconv.Atoi(s[0])
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
			input.Stops = append(input.Stops, schema.Stop{})
		}
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
	depot := stopsByNumber[depotNumber]

	// Create estimated / sufficient vehicles
	for i := 0; i < 10; i++ {
		input.Vehicles = append(input.Vehicles, schema.Vehicle{
			StartLocation: &depot.Location,
			EndLocation:   &depot.Location,
			Capacity:      capacity,
		})
	}

	input.DurationMatrix = matrix

	if byPoint != nil {
		// Create duration matrix
		return fmt.Errorf("implement me")
	}

	return nil
}
