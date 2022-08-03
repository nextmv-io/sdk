// Package sdk provides API declarations of the Nextmv software development
// kit. Please visit the associated packages (directories) for complete
// documentation.
package sdk

import (
	"runtime/debug"
	"strings"
)

// VERSION of Nextmv SDK.
var VERSION string = getVersion()

func getVersion() string {
	bi, _ := debug.ReadBuildInfo()
	for _, dep := range bi.Deps {
		if strings.HasPrefix(dep.Path, "github.com/nextmv-io/sdk") {
			return dep.Version
		}
	}
	return "version_not_found"
}
