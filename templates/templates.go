// Package templates contains variables holding embedded template files.
package templates

import (
	// This package is required to embed files using the `//go:embed` directive
	// albeit it is not used directly.
	_ "embed"
)

var (
	// MipMain is the main.go file in the mip template.
	//go:embed mip/main.go
	MipMain string
	// MipInput is the input.json file in the mip template.
	//go:embed mip/input.json
	MipInput string
	// MipReadme is the README.md file in the mip template.
	//go:embed mip/README.md
	MipReadme string
	// MipManifest is the app.yaml file in the mip template.
	//go:embed mip/app.yaml
	MipManifest string

	// NextrouteMain is the main.go file in the nextroute template.
	//go:embed nextroute/main.go
	NextrouteMain string
	// NextrouteInput is the input.json file in the nextroute template.
	//go:embed nextroute/input.json
	NextrouteInput string
	// NextrouteReadme is the README.md file in the nextroute template.
	//go:embed nextroute/README.md
	NextrouteReadme string
	// NextrouteManifest is the app.yaml file in the nextroute template.
	//go:embed nextroute/app.yaml
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
