package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewTravelDurationObjective returns a new TravelDurationObjective that
// uses the travel duration as an objective.
func NewTravelDurationObjective() TravelDurationObjective {
	connect.Connect(con, &newTravelDurationObjective)
	return newTravelDurationObjective()
}

// TravelDurationObjective is an objective that uses the travel duration as an
// objective.
type TravelDurationObjective interface {
	ModelObjective
}
