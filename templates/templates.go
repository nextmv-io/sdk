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

	// OrtoolsManifest is the app manifest file in the ortools template.
	//go:embed ortools/app.yaml
	OrtoolsManifest string
	// OrtoolsMain is the main.py file in the ortools template.
	//go:embed ortools/main.py
	OrtoolsMain string
	// OrtoolsInput is the input.json file in the ortools template.
	//go:embed ortools/input.json
	OrtoolsInput string
	// OrtoolsReadme is the README.md file in the ortools template.
	//go:embed ortools/README.md
	OrtoolsReadme string
	// OrtoolsRequirements is the requirements.txt file in the ortools template.
	//go:embed ortools/requirements.txt
	OrtoolsRequirements string

	// PyomoManifest is the app manifest file in the pyomo template.
	//go:embed pyomo/app.yaml
	PyomoManifest string
	// PyomoMain is the main.py file in the pyomo template.
	//go:embed pyomo/main.py
	PyomoMain string
	// PyomoInput is the input.json file in the pyomo template.
	//go:embed pyomo/input.json
	PyomoInput string
	// PyomoReadme is the README.md file in the pyomo template.
	//go:embed pyomo/README.md
	PyomoReadme string
	// PyomoRequirements is the requirements.txt file in the pyomo template.
	//go:embed pyomo/requirements.txt
	PyomoRequirements string

	// ShiftCreationManifest is the app manifest file in the shift-creation template.
	//go:embed shift-creation/app.yaml
	ShiftCreationManifest string
	// ShiftCreationMain is the main.py file in the shift-creation template.
	//go:embed shift-creation/main.py
	ShiftCreationMain string
	// ShiftCreationInput is the input.json file in the shift-creation template.
	//go:embed shift-creation/input.json
	ShiftCreationInput string
	// ShiftCreationReadme is the README.md file in the shift-creation template.
	//go:embed shift-creation/README.md
	ShiftCreationReadme string
	// ShiftCreationRequirements is the requirements.txt file in the shift-creation template.
	//go:embed shift-creation/requirements.txt
	ShiftCreationRequirements string

	// ShiftAssignmentManifest is the app manifest file in the shift-assignment template.
	//go:embed shift-assignment/app.yaml
	ShiftAssignmentManifest string
	// ShiftAssignmentMain is the main.py file in the shift-assignment template.
	//go:embed shift-assignment/main.py
	ShiftAssignmentMain string
	// ShiftAssignmentInput is the input.json file in the shift-assignment template.
	//go:embed shift-assignment/input.json
	ShiftAssignmentInput string
	// ShiftAssignmentReadme is the README.md file in the shift-assignment template.
	//go:embed shift-assignment/README.md
	ShiftAssignmentReadme string
	// ShiftAssignmentRequirements is the requirements.txt file in the shift-assignment template.
	//go:embed shift-assignment/requirements.txt
	ShiftAssignmentRequirements string

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
