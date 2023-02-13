// Package sdk provides API declarations of the Nextmv software development
// kit. Please visit the associated packages (directories) for complete
// documentation.
package sdk

import (
	_ "embed"
	"runtime/debug"
	"strings"

	"github.com/nextmv-io/sdk/util"
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
		// We only care about this repo being used as a dependency.
		if !strings.HasPrefix(dep.Path, "github.com/nextmv-io/sdk") {
			continue
		}

		// Clean pseudo versions by falling back to the base version.
		if util.IsPseudoVersion(dep.Version) {
			base, err := util.GetBaseOfPseudoVersion(dep.Version)
			if err != nil {
				panic(err)
			}

			return base
		}

		return dep.Version
	}

	// This should never happen. If it does, it means that the SDK is not being
	// used as a dependency. However, we are apparently being invoked.
	return versionFallback
}
