// Package sdk provides API declarations of the Nextmv software development
// kit. Please visit the associated packages (directories) for complete
// documentation.
package sdk

import (
	_ "embed"
	"runtime/debug"
	"strings"
)

// This will be needed for examples and tests within this repo only.
//
//go:embed VERSION
var versionFallback string

// VERSION of Nextmv SDK.
var VERSION = getVersion()

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return versionFallback
	}

	for _, dep := range bi.Deps {
		// Get only care about this repo being used as a dependency.
		if !strings.HasPrefix(dep.Path, "github.com/nextmv-io/sdk") {
			continue
		}

		// If reference to this module was replaced, use fallback.
		if dep.Replace != nil {
			return versionFallback
		}

		return dep.Version
	}

	return versionFallback
}
