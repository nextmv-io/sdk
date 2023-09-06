// Package templates contains variables holding embded template files.
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
	// ShiftSchedulingInput is the input.json file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/input.json
	ShiftSchedulingInput string
	// ShiftSchedulingReadme is the README.md file in the
	// shift-scheduling template.
	//go:embed shift-scheduling/README.md
	ShiftSchedulingReadme string
)
