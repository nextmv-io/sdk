package golden

import (
	"fmt"
	"math"
	"regexp"
	"time"
)

// replaceTransient replaces all the values in a map whose key is contained in
// the variadic list of transientFields. The replaced value has a stable value
// according to the data type.
func replaceTransient(
	path string,
	original map[string]any,
	transientFields ...TransientField,
) (map[string]any, error) {
	transientLookup := map[string]any{}
	for _, field := range transientFields {
		// See whether we should skip this replacement for the given path.
		skip, err := skipFile(path, field.FileRegex, field.FileRegexFullPath)
		if err != nil {
			return nil, fmt.Errorf("error checking file regex: %w", err)
		}
		if skip {
			continue
		}
		transientLookup[field.Key] = field.Replacement
	}

	replaced := map[string]any{}
	for key, value := range original {
		// Check whether the // Keep the original value.
		replaced[key] = value

		// Check if the field is meant to be replaced. If not, continue.
		// We also check for wildcard replacements that are meant to replace all
		// fields in a slice.
		replacement, isTransient := transientLookup[key]
		cleanedKey := replaceIndicesInKeys(key)
		replacementCleaned, isTransientCleaned := transientLookup[cleanedKey]
		if !isTransient && !isTransientCleaned {
			// No replacement defined, continue and keep the original value.
			continue
		}
		if isTransientCleaned {
			replacement = replacementCleaned
		}

		// Replace the value with the replacement value.
		if replacement != nil {
			replaced[key] = replacement
			continue
		}

		// No replacement defined, we fall back to default stable values here
		// (based on type).

		if stringValue, isString := value.(string); isString {
			if _, err := time.Parse(time.RFC3339, stringValue); err == nil {
				replaced[key] = StableTime
				continue
			}

			if _, err := time.ParseDuration(stringValue); err == nil {
				replaced[key] = StableDuration
				continue
			}

			replaced[key] = StableText
			continue
		}

		if _, isFloat := value.(float64); isFloat {
			replaced[key] = StableFloat
			continue
		}

		if _, IsInt := value.(int); IsInt {
			replaced[key] = StableInt
			continue
		}

		if _, isBool := value.(bool); isBool {
			replaced[key] = StableBool
			continue
		}
	}

	return replaced, nil
}

// roundFields rounds all the values in a map whose key is contained in the
// variadic list of roundingConfigs. The rounded value has a stable value
// according to the data type.
func roundFields(
	path string,
	original map[string]any,
	roundedFields ...RoundingConfig,
) (map[string]any, error) {
	roundingLookup := map[string]int{}
	for _, field := range roundedFields {
		roundingLookup[field.Key] = field.Precision
	}

	replaced := map[string]any{}
	for key, value := range original {
		// Keep the original value.
		replaced[key] = value

		// Check if the field is meant to be rounded. If not, continue.
		// We also check for wildcard replacements that are meant to replace all
		// fields in a slice.
		cleanedKey := replaceIndicesInKeys(key)
		replacement, isRounded := roundingLookup[key]
		replacementCleaned, isRoundedCleaned := roundingLookup[cleanedKey]
		if !isRounded && !isRoundedCleaned {
			// No rounding defined, continue and keep the original value.
			continue
		}
		if isRoundedCleaned {
			replacement = replacementCleaned
		}

		// We don't deal with negative precision values.
		if replacement < 0 {
			continue
		}

		// Replace the value with the rounded value.
		if _, isFloat := value.(float64); isFloat {
			replaced[key] = round(value.(float64), replacement)
			continue
		}

		// If the value was not a float, return an error.
		return nil, fmt.Errorf("field %s is not a float", key)
	}

	return replaced, nil
}

// round rounds a float64 value to a given precision.
func round(value float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(value*shift) / shift
}

var keyIndexMatcher = regexp.MustCompile(`\[\d+\]`)

// replaceIndicesInKeys replaces all the indices in a key with "[]" to make it
// easier to match them with configuration defined in jq-style.
func replaceIndicesInKeys(key string) string {
	return keyIndexMatcher.ReplaceAllString(key, "[]")
}
