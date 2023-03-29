// Package templates contains variables holding embded template files.
package templates

import (
	// This package is required to embed files using the `//go:embed` directive
	// albeit it is not used directly.
	_ "embed"
)

var (
	// RoutingMain is the main.go file in the routing template.
	//go:embed routing/main.go
	RoutingMain string
	// RoutingInput is the input.json file in the routing template.
	//go:embed routing/input.json
	RoutingInput string
	// RoutingReadme is the README.md file in the routing template.
	//go:embed routing/README.md
	RoutingReadme string

	// RoutingInforms is the main.go file in the routing template.
	//go:embed routing-informs/main.go
	RoutingInforms string
	// RoutingInformsInput is the input.json file in the routing template.
	//go:embed routing-informs/input.json
	RoutingInformsInput string
	// RoutingInformsReadme is the README.md file in the routing template.
	//go:embed routing-informs/README.md
	RoutingInformsReadme string

	// RoutingMatrixMain is the main.go file in the routing-matrix template.
	//go:embed routing-matrix-input/main.go
	RoutingMatrixMain string
	// RoutingMatrixInput is the input.json file in the routing-matrix template.
	//go:embed routing-matrix-input/input.json
	RoutingMatrixInput string
	// RoutingMatrixReadme is the README.md file in the routing-matrix template.
	//go:embed routing-matrix-input/README.md
	RoutingMatrixReadme string

	// MeasureMatrixMain is the main.go file in the measure-matrix template.
	//go:embed measure-matrix/main.go
	MeasureMatrixMain string
	// MeasureMatrixReadme is the README.md file in the measure-matrix template.
	//go:embed measure-matrix/README.md
	MeasureMatrixReadme string

	// KnapsackMain is the main.go file in the knapsack template.
	//go:embed knapsack/main.go
	KnapsackMain string
	// KnapsackInput is the input.json file in the knapsack template.
	//go:embed knapsack/input.json
	KnapsackInput string
	// KnapsackReadme is the README.md file in the knapsack template.
	//go:embed knapsack/README.md
	KnapsackReadme string

	// KnapsackMIPMain is the main.go file in the knapsack-mip template.
	//go:embed mip-knapsack/main.go
	KnapsackMIPMain string
	// KnapsackMIPInput is the input.json file in the knapsack-mip template.
	//go:embed mip-knapsack/input.json
	KnapsackMIPInput string
	// KnapsackMIPReadme is the README.md file in the knapsack-mip template.
	//go:embed mip-knapsack/README.md
	KnapsackMIPReadme string

	// ShiftSchedulingMain is the main.go file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/main.go
	ShiftSchedulingMain string
	// ShiftSchedulingInput is the input.json file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/input.json
	ShiftSchedulingInput string
	// ShiftSchedulingReadme is the README.md file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/README.md
	ShiftSchedulingReadme string

	// SudokuMain  is the main.go file in the sudoku template.
	//go:embed sudoku/main.go
	SudokuMain string
	// SudokuInput is the input.json file in the sudoku template.
	//go:embed sudoku/input.json
	SudokuInput string
	// SudokuReadme is the README.md file in the sudoku template.
	//go:embed sudoku/README.md
	SudokuReadme string

	// MipMealAllocationMain is the main.go file in the
	// mip-meal-allocation template.
	//go:embed mip-meal-allocation/main.go
	MipMealAllocationMain string
	// MipMealAllocationInput is the input.json file in the
	// mip-meal-allocation template.
	//go:embed mip-meal-allocation/input.json
	MipMealAllocationInput string
	// MipMealAllocationReadme is the README.md file in the
	// mip-meal-allocation template.
	//go:embed mip-meal-allocation/README.md
	MipMealAllocationReadme string

	// MipIncentiveAllocationMain is the main.go file in the
	// mip-incentive-allocation template.
	//go:embed mip-incentive-allocation/main.go
	MipIncentiveAllocationMain string
	// MipIncentiveAllocationInput is the input.json file in the
	// mip-incentive-allocation template.
	//go:embed mip-incentive-allocation/input.json
	MipIncentiveAllocationInput string
	// MipIncentiveAllocationReadme is the README.md file in the
	// mip-incentive-allocation template.
	//go:embed mip-incentive-allocation/README.md
	MipIncentiveAllocationReadme string

	// CloudRoutingMain is the main.go file in the cloud-routing template.
	//go:embed cloud-routing/main.go
	CloudRoutingMain string
	// CloudRoutingSolver is the solver.go file in the cloud-routing template.
	//go:embed cloud-routing/solver.go
	CloudRoutingSolver string
	// CloudRoutingSchema is the schema.go file in the cloud-routing template.
	//go:embed cloud-routing/schema.go
	CloudRoutingSchema string
	// CloudRoutingHelper is the helper.go file in the cloud-routing template.
	//go:embed cloud-routing/helper.go
	CloudRoutingHelper string
	// CloudRoutingReadme is the README.md file in the cloud-routing template.
	//go:embed cloud-routing/README.md
	CloudRoutingReadme string
	// CloudDeliveryAdvancedInput contains advanced delivery input for
	// the cloud-routing template.
	//go:embed cloud-routing/data/delivery-advanced.json
	CloudDeliveryAdvancedInput string
	// CloudDeliveryBaseInput contains base delivery input for
	// the cloud-routing template.
	//go:embed cloud-routing/data/delivery-base.json
	CloudDeliveryBaseInput string
	// CloudDeliveryRouteLimitInput contains delivery route limit input for
	// the cloud-routing template.
	//go:embed cloud-routing/data/delivery-route-limit.json
	CloudDeliveryRouteLimitInput string
	// CloudDeliveryTinyInput contains tiny delivery input for
	// the cloud-routing template.
	//go:embed cloud-routing/data/delivery-tiny.json
	CloudDeliveryTinyInput string
	// CloudDistributionAdvancedInput contains advanced distribution input for
	// the cloud-routing template.
	//go:embed cloud-routing/data/distribution-advanced.json
	CloudDistributionAdvancedInput string
	// CloudDistributionBaseInput contains base distribution input for
	// the cloud-routing template.
	//go:embed cloud-routing/data/distribution-base.json
	CloudDistributionBaseInput string
	// CloudDistributionRouteLimitInput contains distribution route limit
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/distribution-route-limit.json
	CloudDistributionRouteLimitInput string
	// CloudDistributionTinyInput contains distribution tiny
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/distribution-tiny.json
	CloudDistributionTinyInput string
	// CloudFleetPDInput contains fleet with precedence
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/fleet-pd.json
	CloudFleetPDInput string
	// CloudFleetBaseInput contains fleet base
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/fleet-base.json
	CloudFleetBaseInput string
	// CloudFleetPDTWInput contains fleet with precedence and time windows
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/fleet-pdtw.json
	CloudFleetPDTWInput string
	// CloudFleetTinyInput contains tiny fleet
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/fleet-tiny.json
	CloudFleetTinyInput string
	// CloudSourcingBaseInput contains base sourcing
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/sourcing-base.json
	CloudSourcingBaseInput string
	// CloudSourcingRouteLimitInput contains route limit sourcing
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/sourcing-route-limit.json
	CloudSourcingRouteLimitInput string
	// CloudSourcingTinyInput contains tiny sourcing
	// input for the cloud-routing template.
	//go:embed cloud-routing/data/sourcing-tiny.json
	CloudSourcingTinyInput string

	// NewAppMain is the main.go file in the new-app template.
	//go:embed new-app/main.go
	NewAppMain string
	// NewAppInput is the input.json file in the new-app template.
	//go:embed new-app/input.json
	NewAppInput string
	// NewAppReadme is the README.md file in the new-app template.
	//go:embed new-app/README.md
	NewAppReadme string

	// PagerDutyMain is the main.go file in the pager-duty template.
	//go:embed pager-duty/main.go
	PagerDutyMain string
	// PagerDutyInput is the input.json file in the pager-duty template.
	//go:embed pager-duty/input.json
	PagerDutyInput string
	// PagerDutyReadme is the README.md file in the pager-duty template.
	//go:embed pager-duty/README.md
	PagerDutyReadme string

	// ClusterMain is the main.go file in the cluster template.
	//go:embed cluster/main.go
	ClusterMain string
	// ClusterInput is the input.json file in the cluster template.
	//go:embed cluster/input.json
	ClusterInput string
	// ClusterReadme is the README.md file in the cluster template.
	//go:embed cluster/README.md
	ClusterReadme string

	// ClusterTspMain is the main.go file in the cluster-tsp template.
	//go:embed cluster-tsp/main.go
	ClusterTspMain string
	// ClusterTspReadme is the README.md file in the cluster-tsp template.
	//go:embed cluster-tsp/README.md
	ClusterTspReadme string
	// ClusterTspInput contains input for the cluster-tsp template.
	//go:embed cluster-tsp/input.json
	ClusterTspInput string

	// TimeDependentMain is the main.go file in the routing template.
	//go:embed time-dependent-measure/main.go
	TimeDependentMain string
	// TimeDependentInput is the input.json file in the routing template.
	//go:embed time-dependent-measure/input.json
	TimeDependentInput string
	// TimeDependentReadme is the README.md file in the routing template.
	//go:embed time-dependent-measure/README.md
	TimeDependentReadme string

	// UniqueMatrixMain is the main.go file in the routing template.
	//go:embed unique-matrix-measure/main.go
	UniqueMatrixMain string
	// UniqueMatrixInput is the input.json file in the routing template.
	//go:embed unique-matrix-measure/input.json
	UniqueMatrixInput string
	// UniqueMatrixReadme is the README.md file in the routing template.
	//go:embed unique-matrix-measure/README.md
	UniqueMatrixReadme string
)
