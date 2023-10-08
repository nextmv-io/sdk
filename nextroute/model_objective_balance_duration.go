package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewBalanceDurationPerVehicleObjective creates a new objective that balances
// the total duration across vehicles.
func NewBalanceDurationPerVehicleObjective(mode BalanceObjectiveMode) BalanceDurationPerVehicleObjective {
	connect.Connect(con, &newBalanceDurationObjective)
	return newBalanceDurationObjective(mode)
}

// BalanceDurationPerVehicleObjective is an objective that balances the total
// duration across vehicles.
type BalanceDurationPerVehicleObjective interface {
	ModelObjective
}
