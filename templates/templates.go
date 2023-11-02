// Package templates contains variables holding embedded template files.
package templates

import (
	// This package is required to embed files using the `//go:embed` directive
	// albeit it is not used directly.
	_ "embed"
)

var (
	// CoptManifest is the app manifest file in the copt template.
	//go:embed copt/app.yaml
	CoptManifest string
	// CoptMain is the main.py file in the copt template.
	//go:embed copt/main.py
	CoptMain string
	// CoptInput is the input.json file in the copt template.
	//go:embed copt/input.json
	CoptInput string
	// CoptReadme is the README.md file in the copt template.
	//go:embed copt/README.md
	CoptReadme string
	// CoptRequirements is the requirements.txt file in the copt template.
	//go:embed copt/requirements.txt
	CoptRequirements string

	// GamsManifest is the app manifest file in the gams template.
	//go:embed gams/app.yaml
	GamsManifest string
	// GamsMain is the main.py file in the gams template.
	//go:embed gams/main.py
	GamsMain string
	// GamsInput is the input.json file in the gams template.
	//go:embed gams/input.json
	GamsInput string
	// GamsReadme is the README.md file in the gams template.
	//go:embed gams/README.md
	GamsReadme string
	// GamsRequirements is the requirements.txt file in the gams template.
	//go:embed gams/requirements.txt
	GamsRequirements string

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
