package schema

import "time"

type JsonSolution struct {
	Epoch      time.Time       `json:"epoch"`
	Statistics Statistics      `json:"statistics"`
	Unplanned  []JsonModelStop `json:"unplanned"`
	Vehicles   []JsonVehicle   `json:"vehicles"`
	Objective  JsonObjective   `json:"objective"`
	Summary    OutputSummary   `json:"value_summary"`
}

type JsonObjective struct {
	Name       string          `json:"name"`
	Objectives []JsonObjective `json:"objectives"`
	Value      float64         `json:"value"`
}

type JsonVehicle struct {
	Start time.Time  `json:"start"`
	End   time.Time  `json:"end"`
	Name  string     `json:"name"`
	Stops []JsonStop `json:"stops"`
	OutputVehicle
	Index int `json:"index"`
}

type JsonModelStop struct {
	Name     string       `json:"name"`
	Location JsonLocation `json:"location"`
	Index    int          `json:"index"`
}

type JsonStop struct {
	OutputStop
	Arrival          time.Time             `json:"arrival"`
	EarliestStart    time.Time             `json:"earliest_start"`
	Start            time.Time             `json:"start"`
	End              time.Time             `json:"end"`
	ExpressionValues []JsonExpressionValue `json:"expression_values"`
	ConstraintValues []JsonConstraintValue `json:"constraint_values"`
	JsonModelStop
	Position       int           `json:"position"`
	TravelDuration time.Duration `json:"travel_duration"`
}
type JsonExpressionValue struct {
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Cumulative float64 `json:"cumulative"`
}

type JsonConstraintValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type JsonLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// OutputSummary contains a breakdown of the objective value.
type OutputSummary struct {
	TotalEarlinessPenalty           *int    `json:"total_earliness_penalty,omitempty"`
	TotalLatenessPenalty            *int    `json:"total_lateness_penalty,omitempty"`
	Value                           float64 `json:"value"`
	TotalTravelDistance             float64 `json:"total_travel_distance"`
	TotalTravelTime                 int     `json:"total_travel_time"`
	TotalUnassignedPenalty          float64 `json:"total_unassigned_penalty"`
	TotalVehicleInitializationCosts int     `json:"total_vehicle_initialization_costs"`
}

// OutputVehicle contains the route for one vehicle.
type OutputVehicle struct {
	EarlinessPenalty           *int         `json:"earliness_penalty,omitempty"`
	LatenessPenalty            *int         `json:"lateness_penalty,omitempty"`
	ID                         string       `json:"id"`
	Polyline                   string       `json:"polyline,omitempty"`
	Route                      []OutputStop `json:"route,omitempty"`
	Value                      int          `json:"value"`
	TravelDistance             float64      `json:"travel_distance"`
	TravelTime                 int          `json:"travel_time"`
	VehicleInitializationCosts int          `json:"vehicle_initialization_costs"`
}

// OutputStop contains all data of a stop.
type OutputStop struct {
	Distance         *float64   `json:"distance,omitempty"`
	ETA              *time.Time `json:"eta,omitempty"`
	ETS              *time.Time `json:"ets,omitempty"`
	ETD              *time.Time `json:"etd,omitempty"`
	EarlinessPenalty *int       `json:"earliness_penalty,omitempty"`
	LatenessPenalty  *int       `json:"lateness_penalty,omitempty"`
	ID               string     `json:"id"`
	Polyline         string     `json:"polyline,omitempty"`
	Lon              float64    `json:"lon"`
	Lat              float64    `json:"lat"`
}

type Statistics struct {
	Value *float64 `json:"value"`
	Time  struct {
		ElapsedSeconds *float64 `json:"elapsed_seconds"`
	} `json:"time"`
}
