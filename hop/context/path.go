package context

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
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

type Target struct {
	os, arch string
}

func getPath() (string, error) {
	nextmvLibraryNamePrefix := defaultNextmvLibraryNamePrefix

	nextmvLibraryPath := os.Getenv("NEXTMV_LIBRARY_PATH")

	pathPrefix := ""

	if len(nextmvLibraryPath) > 0 {
		nextmvLibraryCleanPath := path.Clean(nextmvLibraryPath)

		if _, err := os.Stat(nextmvLibraryCleanPath); os.IsNotExist(err) {
			return "", fmt.Errorf("nextmv library path '%s' does not exist", nextmvLibraryCleanPath)
		}

		pathPrefix = nextmvLibraryCleanPath
	} else {
		home, err := os.UserHomeDir()

		if err != nil {
			return "", err
		}

		pathPrefix = filepath.Join(home, defaultNextmvLibraryPath)
	}

	target, err := getTarget()

	if err != nil {
		return "", err
	}

	version, err := getVersion()

	if err != nil {
		return "", err
	}

	versionDelimiter := ""

	if len(version) > 0 {
		versionDelimiter = "-"
	}

	fileName := fmt.Sprintf("%s-%s-%s%s%s.so", nextmvLibraryNamePrefix, target.os, target.arch, versionDelimiter, version)

	return filepath.Join(pathPrefix, fileName), nil
}

func getVersion() (string, error) {
	version := os.Getenv("NEXTMV_SDK_VERSION")

	return version, nil
}

func getTarget() (Target, error) {
	goos := os.Getenv("NEXTMV_SDK_OS")

	if goos == "" {
		goos = runtime.GOOS
	}

	goarch := os.Getenv("NEXTMV_SDK_ARCH")

	if goarch == "" {
		goarch = runtime.GOARCH
	}

	_, exists := supportedTargets[operatingSystem(goos)][architecture(goarch)]

	if !exists {
		return Target{
			os:   goos,
			arch: goarch,
		}, fmt.Errorf("unsupported target, os %s, architecture %s", goos, goarch)
	}

	return Target{
		os:   goos,
		arch: goarch,
	}, nil
}
