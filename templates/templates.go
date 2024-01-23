// Package templates contains variables holding embedded template files.
package templates

import (
	// This package is required to embed files using the `//go:embed` directive
	// albeit it is not used directly.
	_ "embed"
)

var (
	// KnapsackGoMain is the main.go file in the mip template.
	//go:embed knapsack-gosdk/main.go
	KnapsackGoMain string
	// KnapsackGoInput is the input.json file in the mip template.
	//go:embed knapsack-gosdk/input.json
	KnapsackGoInput string
	// KnapsackGoReadme is the README.md file in the mip template.
	//go:embed knapsack-gosdk/README.md
	KnapsackGoReadme string
	// KnapsackGoManifest is the app.yaml file in the mip template.
	//go:embed knapsack-gosdk/app.yaml
	KnapsackGoManifest string

	// NextrouteMain is the main.go file in the nextroute template.
	//go:embed nextroute-gosdk/main.go
	NextrouteMain string
	// NextrouteInput is the input.json file in the nextroute template.
	//go:embed nextroute-gosdk/input.json
	NextrouteInput string
	// NextrouteReadme is the README.md file in the nextroute template.
	//go:embed nextroute-gosdk/README.md
	NextrouteReadme string
	// NextrouteManifest is the app.yaml file in the nextroute template.
	//go:embed nextroute-gosdk/app.yaml
	NextrouteManifest string

	// KnapsackOrtoolsManifest is the app manifest file in the ortools template.
	//go:embed knapsack-ortools/app.yaml
	KnapsackOrtoolsManifest string
	// KnapsackOrtoolsMain is the main.py file in the ortools template.
	//go:embed knapsack-ortools/main.py
	KnapsackOrtoolsMain string
	// KnapsackOrtoolsInput is the input.json file in the ortools template.
	//go:embed knapsack-ortools/input.json
	KnapsackOrtoolsInput string
	// KnapsackOrtoolsReadme is the README.md file in the ortools template.
	//go:embed knapsack-ortools/README.md
	KnapsackOrtoolsReadme string
	// KnapsackOrtoolsRequirements is the requirements.txt file in the ortools template.
	//go:embed knapsack-ortools/requirements.txt
	KnapsackOrtoolsRequirements string

	// KnapsackPyomoManifest is the app manifest file in the pyomo template.
	//go:embed knapsack-pyomo/app.yaml
	KnapsackPyomoManifest string
	// KnapsackPyomoMain is the main.py file in the pyomo template.
	//go:embed knapsack-pyomo/main.py
	KnapsackPyomoMain string
	// KnapsackPyomoInput is the input.json file in the pyomo template.
	//go:embed knapsack-pyomo/input.json
	KnapsackPyomoInput string
	// KnapsackPyomoReadme is the README.md file in the pyomo template.
	//go:embed knapsack-pyomo/README.md
	KnapsackPyomoReadme string
	// KnapsackPyomoRequirements is the requirements.txt file in the pyomo template.
	//go:embed knapsack-pyomo/requirements.txt
	KnapsackPyomoRequirements string

	// RoutingOrtoolsManifest is the app manifest file in the pyomo template.
	//go:embed routing-ortools/app.yaml
	RoutingOrtoolsManifest string
	// RoutingOrtoolsMain is the main.py file in the pyomo template.
	//go:embed routing-ortools/main.py
	RoutingOrtoolsMain string
	// RoutingOrtoolsInput is the input.json file in the pyomo template.
	//go:embed routing-ortools/input.json
	RoutingOrtoolsInput string
	// RoutingOrtoolsReadme is the README.md file in the pyomo template.
	//go:embed routing-ortools/README.md
	RoutingOrtoolsReadme string
	// RoutingOrtoolsRequirements is the requirements.txt file in the pyomo template.
	//go:embed routing-ortools/requirements.txt
	RoutingOrtoolsRequirements string

	// ShiftPlanningOrtoolsManifest is the app manifest file in the shift-planning template.
	//go:embed shift-planning-ortools/app.yaml
	ShiftPlanningOrtoolsManifest string
	// ShiftPlanningOrtoolsMain is the main.py file in the shift-planning template.
	//go:embed shift-planning-ortools/main.py
	ShiftPlanningOrtoolsMain string
	// ShiftPlanningOrtoolsInput is the input.json file in the shift-planning template.
	//go:embed shift-planning-ortools/input.json
	ShiftPlanningOrtoolsInput string
	// ShiftPlanningOrtoolsReadme is the README.md file in the shift-planning template.
	//go:embed shift-planning-ortools/README.md
	ShiftPlanningOrtoolsReadme string
	// ShiftPlanningOrtoolsRequirements is the requirements.txt file in the shift-planning template.
	//go:embed shift-planning-ortools/requirements.txt
	ShiftPlanningOrtoolsRequirements string

	// ShiftPlanningPyomoManifest is the app manifest file in the shift-planning template.
	//go:embed shift-planning-pyomo/app.yaml
	ShiftPlanningPyomoManifest string
	// ShiftPlanningPyomoMain is the main.py file in the shift-planning template.
	//go:embed shift-planning-pyomo/main.py
	ShiftPlanningPyomoMain string
	// ShiftPlanningPyomoInput is the input.json file in the shift-planning template.
	//go:embed shift-planning-pyomo/input.json
	ShiftPlanningPyomoInput string
	// ShiftPlanningPyomoReadme is the README.md file in the shift-planning template.
	//go:embed shift-planning-pyomo/README.md
	ShiftPlanningPyomoReadme string
	// ShiftPlanningPyomoRequirements is the requirements.txt file in the shift-planning template.
	//go:embed shift-planning-pyomo/requirements.txt
	ShiftPlanningPyomoRequirements string

	// ShiftAssignmentOrtoolsManifest is the app manifest file in the shift-assignment template.
	//go:embed shift-assignment-ortools/app.yaml
	ShiftAssignmentOrtoolsManifest string
	// ShiftAssignmentOrtoolsMain is the main.py file in the shift-assignment template.
	//go:embed shift-assignment-ortools/main.py
	ShiftAssignmentOrtoolsMain string
	// ShiftAssignmentOrtoolsInput is the input.json file in the shift-assignment template.
	//go:embed shift-assignment-ortools/input.json
	ShiftAssignmentOrtoolsInput string
	// ShiftAssignmentOrtoolsReadme is the README.md file in the shift-assignment template.
	//go:embed shift-assignment-ortools/README.md
	ShiftAssignmentOrtoolsReadme string
	// ShiftAssignmentOrtoolsRequirements is the requirements.txt file in the shift-assignment template.
	//go:embed shift-assignment-ortools/requirements.txt
	ShiftAssignmentOrtoolsRequirements string

	// ShiftAssignmentPyomoManifest is the app manifest file in the shift-assignment template.
	//go:embed shift-assignment-pyomo/app.yaml
	ShiftAssignmentPyomoManifest string
	// ShiftAssignmentPyomoMain is the main.py file in the shift-assignment template.
	//go:embed shift-assignment-pyomo/main.py
	ShiftAssignmentPyomoMain string
	// ShiftAssignmentPyomoInput is the input.json file in the shift-assignment template.
	//go:embed shift-assignment-pyomo/input.json
	ShiftAssignmentPyomoInput string
	// ShiftAssignmentPyomoReadme is the README.md file in the shift-assignment template.
	//go:embed shift-assignment-pyomo/README.md
	ShiftAssignmentPyomoReadme string
	// ShiftAssignmentPyomoRequirements is the requirements.txt file in the shift-assignment template.
	//go:embed shift-assignment-pyomo/requirements.txt
	ShiftAssignmentPyomoRequirements string

	// ShiftSchedulingMain is the main.go file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/main.go
	ShiftSchedulingMain string
	// ShiftSchedulingSchema is the schema.go file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/schema.go
	ShiftSchedulingSchema string
	// ShiftSchedulingInput is the input.json file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/input.json
	ShiftSchedulingInput string
	// ShiftSchedulingReadme is the README.md file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/README.md
	ShiftSchedulingReadme string
	// ShiftSchedulingManifest is the app.yaml file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/app.yaml
	ShiftSchedulingManifest string

	// OrderFulfillmentMain is the main.go file in the
	// order-fulfillment template.
	//go:embed order-fulfillment-gosdk/main.go
	OrderFulfillmentMain string
	// OrderFulfillmentInput is the input.json file in the
	// order-fulfillment template.
	//go:embed order-fulfillment-gosdk/input.json
	OrderFulfillmentInput string
	// OrderFulfillmentReadme is the README.md file in the
	// order-fulfillment template.
	//go:embed order-fulfillment-gosdk/README.md
	OrderFulfillmentReadme string
	// OrderFulfillmentManifest is the app.yaml file in the
	// order-fulfillment template.
	//go:embed order-fulfillment-gosdk/app.yaml
	OrderFulfillmentManifest string

	// XpressManifest is the app manifest file in the xpress template.
	//go:embed xpress/app.yaml
	XpressManifest string
	// XpressMain is the main.py file in the xpress template.
	//go:embed xpress/main.py
	XpressMain string
	// XpressInput is the input.json file in the xpress template.
	//go:embed xpress/input.json
	XpressInput string
	// XpressReadme is the README.md file in the xpress template.
	//go:embed xpress/README.md
	XpressReadme string
	// XpressRequirements is the requirements.txt file in the xpress template.
	//go:embed xpress/requirements.txt
	XpressRequirements string
)
