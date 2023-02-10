// Package util contains utility functions.
package util

import (
	"fmt"
	"regexp"
	"strconv"
)

// IsPseudoVersion returns true if the Go version string is a pseudo version.
func IsPseudoVersion(version string) bool {
	regex := regexp.MustCompile(`^v\d+\.\d+\.\d+-0\.[0-9a-f]{14}-[0-9a-f]{12}$`)
	return regex.MatchString(version)
}

// GetBaseOfPseudoVersion gets the version a pseudo version is based on.
func GetBaseOfPseudoVersion(version string) (string, error) {
	if !IsPseudoVersion(version) {
		return "", fmt.Errorf("version %s is not a pseudo version", version)
	}
	crop := version[:len(version)-30]
	// Decrement the patch part by one.
	regex := regexp.MustCompile(`\d+$`)
	patch := regex.FindString(crop)
	patchInt, err := strconv.Atoi(patch)
	if err != nil {
		return "", err
	}
	patchInt--
	return regex.ReplaceAllString(crop, strconv.Itoa(patchInt)), nil
}
