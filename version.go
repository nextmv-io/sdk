// Package sdk provides API declarations of the Nextmv software development
// kit. Please visit the associated packages (directories) for complete
// documentation.
package sdk

import (
	"runtime/debug"
	"strings"
)

// This will be needed for examples and tests within this repo only.
var versionFallback string = "v0.16.0-dev.0-5"

// VERSION of Nextmv SDK.
var VERSION string = getVersion()

func getVersion() string {
	bi, _ := debug.ReadBuildInfo()
	for _, dep := range bi.Deps {
		if strings.HasPrefix(dep.Path, "github.com/nextmv-io/sdk") {
			return dep.Version
		}
	}
	return versionFallback
}
