// Package sdk provides API declarations of the Nextmv software development
// kit. Please visit the associated packages (directories) for complete
// documentation.
package sdk

import (
	_ "embed"
	"os"
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
	// internal test expectations use <<PRESENCE>> as the version so we do not
	// have to update expectations every time the version changes
	if ver, ok := os.LookupEnv("USE_PRESENCE"); ok && ver == "1" {
		return "<<PRESENCE>>"
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return versionFallback
	}

	for _, dep := range bi.Deps {
		// We only care about this repo being used as a dependency.
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
