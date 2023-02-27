package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewTravelDurationObjective returns a new TravelDurationObjective that
// uses the travel duration as an objective.
func NewTravelDurationObjective(
	factor float64,
) TravelDurationObjective {
	connect.Connect(con, &newTravelDurationObjective)
	return newTravelDurationObjective(factor)
}

// TravelDurationObjective is an objective that uses the travel duration as an
// objective.
type TravelDurationObjective interface {
	ModelObjective
}
