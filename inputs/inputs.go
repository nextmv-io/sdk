// Package inputs contains variables holding embed example input files
// for use with the cloud routing API.
package inputs

import (
	// Embed ...
	_ "embed"
)

var (
	// RoutingReadme is the README.md file for the cloud routing API inputs.
	//go:embed routing/README.md
	RoutingReadme string
	// RoutingTinyFleet is the tiny fleet input for the cloud routing API.
	//go:embed routing/fleet-tiny.json
	RoutingTinyFleet string
	// RoutingFleet is the fleet with precedence input for the
	// cloud routing API.
	//go:embed routing/fleet-pd.json
	RoutingFleet string
	// RoutingTinyDelivery is the tiny delivery input for the
	// cloud routing API.
	//go:embed routing/delivery-tiny.json
	RoutingTinyDelivery string
	// RoutingDelivery is the delivery input for the cloud routing API.
	//go:embed routing/delivery-advanced.json
	RoutingDelivery string
	// RoutingTinyDistro is the tiny distribution input for the
	// cloud routing API.
	//go:embed routing/distribution-tiny.json
	RoutingTinyDistro string
	// RoutingDistribution is the distribution input for the
	// cloud routing API.
	//go:embed routing/distribution-route-limit.json
	RoutingDistribution string
	// RoutingTinySourcing is the tiny sourcing input for the
	// cloud routing API.
	//go:embed routing/sourcing-tiny.json
	RoutingTinySourcing string
	// OnfleetReadme is the README.md file for the cloud
	// routing API with onfleet integration.
	//go:embed onfleet/README.md
	OnfleetReadme string
	// OnfleetBase is the sample onfleet input for the
	// cloud routing API with onfleet integration.
	//go:embed onfleet/onfleet.json
	OnfleetBase string
)
