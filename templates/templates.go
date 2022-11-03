// Package templates holds the embeddings of the different files used in
// templates.
package templates

import (
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
)
