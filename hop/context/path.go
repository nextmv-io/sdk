package context

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/nextmv-io/sdk"
)

type operatingSystem string

const (
	supportedOperatingSystemDarwin = operatingSystem("darwin")
	supportedOperatingSystemLinux  = operatingSystem("linux")
)

type architecture string

const (
	supportedArchitecture386   = architecture("386")
	supportedArchitectureAmd64 = architecture("amd64")
	supportedArchitectureArm64 = architecture("arm64")
)

type targets map[operatingSystem]map[architecture][]string

var supportedTargets = targets{
	supportedOperatingSystemDarwin: {
		supportedArchitectureAmd64: {},
		supportedArchitectureArm64: {},
	},
	supportedOperatingSystemLinux: {
		supportedArchitecture386:   {},
		supportedArchitectureAmd64: {},
		supportedArchitectureArm64: {},
	},
}

const (
	defaultNextmvLibraryPath       = ".nextmv/lib"
	defaultNextmvLibraryNamePrefix = "nextmv-sdk"
)

func getPath() (string, error) {
	nextmvLibraryNamePrefix := defaultNextmvLibraryNamePrefix

	nextmvLibraryPath := os.Getenv("NEXTMV_LIBRARY_PATH")

	pathPrefix := ""

	if nextmvLibraryPath != "" {
		nextmvLibraryCleanPath := path.Clean(nextmvLibraryPath)

		if _, err := os.Stat(nextmvLibraryCleanPath); os.IsNotExist(err) {
			return "",
				fmt.Errorf("nextmv library path '%s' does not exist",
					nextmvLibraryCleanPath)
		}

		pathPrefix = nextmvLibraryCleanPath
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		pathPrefix = filepath.Join(home, defaultNextmvLibraryPath)
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH

	_, exists := supportedTargets[operatingSystem(goos)][architecture(goarch)]

	var err error

	if !exists {
		err = fmt.Errorf("unsupported target, os %s, architecture %s",
			goos,
			goarch)
	}

	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s-%s-%s-%s.so",
		nextmvLibraryNamePrefix,
		goos,
		goarch,
		sdk.VERSION)

	return filepath.Join(pathPrefix, fileName), nil
}
