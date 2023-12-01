package here

import (
	"encoding/json"
)

type validationError struct {
	message string
}

func (e validationError) Error() string {
	return e.message
}

type responseStatus string

const (
	responseStatusComplete   responseStatus = "completed"
	responseStatusAccepted   responseStatus = "accepted"
	responseStatusInProgress responseStatus = "inProgress"
)

func isKnownStatusResponse(status responseStatus) bool {
	return status == responseStatusComplete ||
		status == responseStatusAccepted ||
		status == responseStatusInProgress
}

type statusResponse struct {
	MatrixID  string          `json:"matrixId"` //nolint:tagliatelle
	Status    responseStatus  `json:"status"`
	StatusURL string          `json:"statusUrl"` //nolint:tagliatelle
	ResultURL string          `json:"resultUrl"` //nolint:tagliatelle
	Error     json.RawMessage `json:"error"`
}

type matrixResponse struct {
	Matrix           matrix           `json:"matrix"`
	RegionDefinition regionDefinition `json:"regionDefinition"` //nolint:tagliatelle,lll
}

type matrix struct {
	NumOrigins      int   `json:"numOrigins"`      //nolint:tagliatelle
	NumDestinations int   `json:"numDestinations"` //nolint:tagliatelle
	TravelTimes     []int `json:"travelTimes"`     //nolint:tagliatelle
	Distances       []int `json:"distances"`
}

type regionDefinition struct {
	Type string `json:"type"`
}

type point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lng"` //nolint:tagliatelle
}

type matrixRequest struct {
	Origins          []point          `json:"origins"`
	RegionDefinition regionDefinition `json:"regionDefinition,omitempty"` //nolint:tagliatelle,lll
	// This is either an RFC 3339 timestamp or the string "any"
	DepartureTime    string        `json:"departureTime,omitempty"` //nolint:tagliatelle,lll
	MatrixAttributes []string      `json:"matrixAttributes"`        //nolint:tagliatelle,lll
	TransportMode    TransportMode `json:"transportMode,omitempty"` //nolint:tagliatelle,lll
	Avoid            *avoid        `json:"avoid,omitempty"`
	Truck            *Truck        `json:"truck,omitempty"`
	Scooter          *Scooter      `json:"scooter,omitempty"`
	Taxi             *Taxi         `json:"taxi,omitempty"`
}

type avoid struct {
	// Features string `json:"features,omitempty"`
	Features []Feature `json:"features,omitempty"`
	Areas    []area    `json:"areas,omitempty"`
}

type area struct {
	Type  string  `json:"type"`
	North float64 `json:"north"`
	South float64 `json:"south"`
	East  float64 `json:"east"`
	West  float64 `json:"west"`
}

// BoundingBox represents a region using four cooordinates corresponding
// to the furthest points in each of the cardinal directions within that region.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type BoundingBox struct {
	North float64
	South float64
	East  float64
	West  float64
}

// TransportMode represents the type of vehicle that will be used for
// the calculated routes.#//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type TransportMode string

// TransportModeCar causes routes to be calculated for car travel.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TransportModeCar TransportMode = "car"

// TransportModeTruck causes routes to be calculated for truck travel.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TransportModeTruck TransportMode = "truck"

// TransportModePedestrian causes routes to be calculated for pedestrian travel.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TransportModePedestrian TransportMode = "pedestrian"

// TransportModeBicycle causes routes to be calculated for bicycle travel.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TransportModeBicycle TransportMode = "bicycle"

// TransportModeTaxi causes routes to be calculated for taxi travel.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TransportModeTaxi TransportMode = "taxi"

// TransportModeScooter causes routes to be calculated for scooter travel.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TransportModeScooter TransportMode = "scooter"

// Feature represents a geographical feature.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type Feature string

// TollRoad designates a toll road feature.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TollRoad Feature = "tollRoad"

// ControlledAccessHighway designates a controlled access highway.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const ControlledAccessHighway Feature = "controlledAccessHighway"

// Ferry designates a ferry route.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Ferry Feature = "ferry"

// Tunnel designates a tunnel.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Tunnel Feature = "tunnel"

// DirtRoad designates a dirt road.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const DirtRoad Feature = "dirtRoad"

// SeasonalClosure designates a route that is closed for the season.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const SeasonalClosure Feature = "seasonalClosure"

// CarShuttleTrain designates a train that can transport cars.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const CarShuttleTrain Feature = "carShuttleTrain"

// DifficultTurns represents u-turns, difficult turns, and sharp turns.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const DifficultTurns Feature = "difficultTurns"

// UTurns designates u-turns.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const UTurns Feature = "uTurns"

// Truck captures truck-specific routing parameters.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type Truck struct {
	ShippedHazardousGoods []HazardousGood `json:"shippedHazardousGoods,omitempty"` //nolint:tagliatelle,lll
	// in kilograms
	GrossWeight int32 `json:"grossWeight,omitempty"` //nolint:tagliatelle
	// in kilograms
	WeightPerAxle int32 `json:"weightPerAxle,omitempty"` //nolint:tagliatelle
	// in centimeters
	Height int32 `json:"height,omitempty"`
	// in centimeters
	Width int32 `json:"width,omitempty"`
	// in centimeters
	Length             int32               `json:"length,omitempty"`
	TunnelCategory     TunnelCategory      `json:"tunnelCategory,omitempty"` //nolint:tagliatelle,lll
	AxleCount          int32               `json:"axleCount,omitempty"`      //nolint:tagliatelle,lll
	Type               TruckType           `json:"type,omitempty"`
	TrailerCount       int32               `json:"trailerCount,omitempty"`       //nolint:tagliatelle,lll
	WeightPerAxleGroup *WeightPerAxleGroup `json:"weightPerAxleGroup,omitempty"` //nolint:tagliatelle,lll
}

// WeightPerAxleGroup captures the weights of different axle groups.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type WeightPerAxleGroup struct {
	Single int32 `json:"single"`
	Tandem int32 `json:"tandem"`
	Triple int32 `json:"triple"`
}

// TunnelCategory is a tunnel category used to restrict the transport of
// certain goods.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type TunnelCategory string

// TunnelCategoryB represents tunnels with B category restrictions.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TunnelCategoryB TunnelCategory = "B"

// TunnelCategoryC represents tunnels with C category restrictions.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TunnelCategoryC TunnelCategory = "C"

// TunnelCategoryD represents a tunnel with D category restrictions.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TunnelCategoryD TunnelCategory = "D"

// TunnelCategoryE represents a tunnel with E category restrictions.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TunnelCategoryE TunnelCategory = "E"

// TunnelCategoryNone represents a tunnel with no category restrictions.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TunnelCategoryNone TunnelCategory = "None"

// TruckType specifies the type of truck.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type TruckType string

// TruckTypeStraight refers to trucks with a permanently attached cargo area.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TruckTypeStraight TruckType = "straight"

// TruckTypeTractor refers to vehicles that can tow one or more semi-trailers.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const TruckTypeTractor TruckType = "tractor"

// HazardousGood indicates a hazardous good that trucks can transport.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type HazardousGood string

// Explosive represents explosive materials.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Explosive HazardousGood = "explosive"

// Gas designates gas.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Gas HazardousGood = "gas"

// Flammable designates flammable materials.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Flammable HazardousGood = "flammable"

// Combustible designates combustible materials.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Combustible HazardousGood = "combustible"

// Organic designates organical materials.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Organic HazardousGood = "organic"

// Poison designates poison.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Poison HazardousGood = "poison"

// Radioactive indicates radioactive materials.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Radioactive HazardousGood = "radioactive"

// Corrosive indicates corrosive materials.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const Corrosive HazardousGood = "corrosive"

// PoisonousInhalation refers to materials that are poisonous to inhale.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const PoisonousInhalation HazardousGood = "poisonousInhalation"

// HarmfulToWater indicates materials that are harmful to water.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const HarmfulToWater HazardousGood = "harmfulToWater"

// OtherHazardousGood refers to other types of hazardous materials.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
const OtherHazardousGood HazardousGood = "other"

// Scooter captures routing parameters that can be set on scooters.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type Scooter struct {
	AllowHighway bool `json:"allowHighway"` //nolint:tagliatelle
}

// Taxi captures routing parameters that can be set on taxis.
//
// Deprecated: This package is deprecated and will be removed in a future.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/here].
type Taxi struct {
	AllowDriveThroughTaxiRoads bool `json:"allowDriveThroughTaxiRoads"` //nolint:tagliatelle,lll
}
