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
	// If an override version is set, use it.
	if v := os.Getenv("NEXTMV_SDK_OVERRIDE_VERSION"); v != "" {
		return v
	}

	// Get version from module dependency.
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		// If this happens, we're not running in a module context. We default to
		// returning the fallback version.
		return versionFallback
	}

	for _, dep := range bi.Deps {
		// We only care about this repo being used as a dependency.
		if !strings.HasPrefix(dep.Path, "github.com/nextmv-io/nextroute") {
			continue
		}

		return dep.Version
	}

	// If this happens, we're running in a module in which sdk is not a
	// dependency. In this case, we expect the NEXTMV_SDK_OVERRIDE_VERSION to be
	// set. Thus, this is unexpected. So, return a string that helps us find the
	// way back here.
	return "no-overridden-version-found"
}
