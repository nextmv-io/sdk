package context

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type SupportedOperatingSystem string

const (
	SupportedOperatingSystemDarwin = SupportedOperatingSystem("darwin")
	SupportedOperatingSystemLinux  = SupportedOperatingSystem("linux")
)

type SupportedArchitecture string

const (
	SupportedArchitecture386   = SupportedArchitecture("386")
	SupportedArchitectureAmd64 = SupportedArchitecture("amd64")
	SupportedArchitectureArm64 = SupportedArchitecture("arm64")
)

type SupportedTargets map[SupportedOperatingSystem]map[SupportedArchitecture][]string

var supportedTargets = SupportedTargets{
	SupportedOperatingSystemDarwin: {
		SupportedArchitectureAmd64: {},
		SupportedArchitectureArm64: {},
	},
	SupportedOperatingSystemLinux: {
		SupportedArchitecture386:   {},
		SupportedArchitectureAmd64: {},
		SupportedArchitectureArm64: {},
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

	versionDeliminator := ""

	if len(version) > 0 {
		versionDeliminator = "-"
	}

	fileName := fmt.Sprintf("%s-%s-%s%s%s.so", nextmvLibraryNamePrefix, target.os, target.arch, versionDeliminator, version)

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

	_, exists := supportedTargets[SupportedOperatingSystem(goos)][SupportedArchitecture(goarch)]

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
