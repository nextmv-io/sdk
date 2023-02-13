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
		// Just return the fallback version.
		return versionFallback
	}

	for _, dep := range bi.Deps {
		// We only care about this repo being used as a dependency.
		if !strings.HasPrefix(dep.Path, "github.com/nextmv-io/sdk") {
			continue
		}

		return dep.Version
	}

	// This should never happen. If it does, it means that the SDK was used on
	// its own and not as a dependency. In that case, we fallback to the version
	// in the VERSION file.
	return versionFallback
}
