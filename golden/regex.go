package golden

import (
	"fmt"
	"regexp"
)

// VolatileRegexReplacement defines a regex replacement to be applied to the
// output before comparison.
type VolatileRegexReplacement struct {
	// Regex to match.
	Regex string
	// Replacement to apply.
	Replacement string
}

// regexReplaceAllDefault applies all default regex replacements to the given
// text.
func regexReplaceAllDefault(text string) string {
	text = regexReplaceEndpoints(text)
	text = regexReplaceGoSemverLike(text, StableVersion)
	text = regexReplaceGUID(text, "00000000-0000-0000-0000-000000000000")
	text = regexReplaceElapsed(text, StableDuration)
	text = regexReplaceElapsedSeconds(text, StableFloat)
	text = regexReplaceStart(text, StableTime)
	return text
}

// regexReplaceEndpoints replaces all endpoints with the default endpoint.
func regexReplaceEndpoints(text string) string {
	text = regexp.
		MustCompile(`us1.api.staging.nxmv.xyz`).
		ReplaceAllString(text, "api.cloud.nextmv.io")
	text = regexp.
		MustCompile(`us1.api.development.nxmv.xyz`).
		ReplaceAllString(text, "api.cloud.nextmv.io")
	return text
}

// regexReplaceGoSemverLike replaces any occurrence of a go semver version with
// optional suffix with 'placeholder'. This is to avoid having to update the
// golden files every time the version changes.
func regexReplaceGoSemverLike(text, placeholder string) string {
	return regexp.
		MustCompile(`v\d+\.\d+\.\d+(\-[0-9,a-z,\-,\.]+)?`).
		ReplaceAllString(text, placeholder)
}

// regexReplaceGUID replaces any occurrence of a GUID with a fixed one. This is
// necessary, since the GUIDs are generated randomly and thus vary between runs.
func regexReplaceGUID(text, placeholder string) string {
	return regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}`+
		`-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`).
		ReplaceAllString(text, placeholder)
}

// regexReplaceElapsed replaces the "elapsed" value in a JSON with
// 'placeholder', since it varies among runs. Only the value is replaced, not
// the key.
//
// Example match:
//
//	"elapsed": "2.1325ms"
func regexReplaceElapsed(text, placeholder string) string {
	// A duration may be formatted in different ways. The following regex makes
	// sure we cover them all.
	//
	// Examples of valid durations:
	//   552h59m59.999999999s
	//   999.999999ms
	//   999.999µs
	//   999ns
	//   552h59m59s
	return regexp.
		MustCompile(`"elapsed":\s*"(\d+h)?(\d+m)?\d+(\.\d+)?(s|ms|µs|ns)"`).
		ReplaceAllString(text, fmt.Sprintf(`"elapsed": "%s"`, placeholder))
}

// regexReplaceElapsedSeconds replaces the "elapsed_seconds" value in a JSON
// with 'placeholder', since it varies among runs. Only the value is replaced,
// not the key.
//
// Example match:
//
//	"elapsed_seconds": 0.0021325
//	"elapsed_seconds": 2
func regexReplaceElapsedSeconds(text string, placeholder float64) string {
	return regexp.
		MustCompile(`"elapsed_seconds":\s*\d+(\.\d+)?`).
		ReplaceAllString(text, fmt.Sprintf(`"elapsed_seconds": %v`, placeholder))
}

// regexReplaceStart replaces the "start" value in a JSON with 'placeholder',
// since it varies among runs. Only the value is replaced, not the key.
//
// Example matches:
//
//	"start": "2023-01-12T21:22:57.581596+01:00"
//	"start": "2023-01-12T21:22:57.581596Z"
func regexReplaceStart(text, placeholder string) string {
	return regexp.
		MustCompile(`"start":\s*"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d+`+
			`((\+|-)\d{2}:\d{2}|Z)"`).
		ReplaceAllString(text, fmt.Sprintf(`"start": "%s"`, placeholder))
}

// regexReplaceCustom replaces any occurrence of a regex with a placeholder.
func regexReplaceCustom(text, placeholder, regex string) string {
	return regexp.
		MustCompile(regex).
		ReplaceAllString(text, placeholder)
}
