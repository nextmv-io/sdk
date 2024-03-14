// Package sdk provides API declarations of the Nextmv software development
// kit. Please visit the associated packages (directories) for complete
// documentation.
package sdk

import (
	_ "embed"
	"runtime/debug"
	"strings"
)

// VERSION of Nextmv SDK.
//
//go:embed VERSION
var versionFallback string

// VERSION of Nextmv SDK.
var VERSION = getVersion()

func getVersion() string {
	// Get version from module dependency.
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		// If this happens, we're not running in a module context. We default to
		// returning the fallback version.
		return versionFallback
	}

	for _, dep := range bi.Deps {
		// We only care about this repo being used as a dependency.
		if !strings.HasPrefix(dep.Path, "github.com/nextmv-io/sdk") {
			continue
		}

		return dep.Version
	}

	// If this happens, we're running in a module in which sdk is not a
	// dependency.
	return versionFallback
}
