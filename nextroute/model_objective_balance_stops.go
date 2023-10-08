package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewBalanceStopsPerVehicleObjective creates a new objective that balances
// stops across vehicles.
func NewBalanceStopsPerVehicleObjective(mode BalanceObjectiveMode) BalanceStopsPerVehicleObjective {
	connect.Connect(con, &newBalanceStopsObjective)
	return newBalanceStopsObjective(mode)
}

// BalanceStopsPerVehicleObjective is an objective that balances stops across
// vehicles.
type BalanceStopsPerVehicleObjective interface {
	ModelObjective
}
