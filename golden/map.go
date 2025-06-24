package golden

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

// filterMatchingJPaths returns paths from candidates that match the wildcard
// filter pattern. It also returns the candidates that matched the wildcard.
func filterMatchingJPaths(candidates []string, wildcard string) (result []string, matches []string) {
	wildParts := parseJPath(wildcard)
	prefixLen := len(wildParts)

	unique := make(map[string]struct{})

	for _, cand := range candidates {
		candParts := parseJPath(cand)
		if isJPathsMatch(wildParts, candParts) {
			// If candidate matches the wildcard, add it to matches.
			matches = append(matches, cand)

			// Also, we want to store the unique prefixes.
			if len(candParts) >= prefixLen {
				// Collapse the candidate path to its prefix.
				prefixParts := candParts[:prefixLen]
				collapsed := collapseJPath(prefixParts)
				unique[collapsed] = struct{}{}
			} else {
				// If candidate is shorter than the wildcard, store it as is.
				unique[cand] = struct{}{}
			}
		}
	}

	// Extract keys as result
	for path := range unique {
		result = append(result, path)
	}
	return result, matches
}

// parseJPath splits a JPath like
// ".grouped_summaries[0].foo[1].bar.custom.duration" into parts like
// ["grouped_summaries", "[0]", "foo", "[1]", "bar", "custom", "duration"].
func parseJPath(jpath string) []string {
	var parts []string
	jpath = strings.TrimPrefix(jpath, "$.")
	var sb strings.Builder

	inBracket := false
	for _, r := range jpath {
		if r == '.' && !inBracket {
			if sb.Len() > 0 {
				parts = append(parts, sb.String())
				sb.Reset()
			}
		} else {
			if r == '[' {
				if sb.Len() > 0 {
					parts = append(parts, sb.String())
					sb.Reset()
				}
				inBracket = true
			}
			if r == ']' {
				inBracket = false
			}
			sb.WriteRune(r)
		}
	}
	if sb.Len() > 0 {
		parts = append(parts, sb.String())
	}
	return parts
}

// collapseJPath combines parts of a JPath into a single string.
func collapseJPath(parts []string) string {
	var sb strings.Builder
	sb.WriteString("$")
	for _, part := range parts {
		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			sb.WriteString(part) // No dot prefix for brackets
		} else {
			sb.WriteString(".")
			sb.WriteString(part)
		}
	}
	return sb.String()
}

// isJPathsMatch checks whether a candidate JPath matches the wildcard
// path pattern.
func isJPathsMatch(pattern, candidate []string) bool {
	pLen := len(pattern)
	cLen := len(candidate)

	if pLen > cLen {
		return false
	}

	for i := 0; i < pLen; i++ {
		pp := pattern[i]
		cp := candidate[i]

		if pp == "[]" {
			if !strings.HasPrefix(cp, "[") || !strings.HasSuffix(cp, "]") {
				return false
			}
		} else if pp != cp {
			return false
		}
	}
	return true
}

// replaceTransient replaces all the values in a map whose key is contained in
// the variadic list of transientFields. The replaced value has a stable value
// according to the data type.
func replaceTransient(
	path string,
	original map[string]any,
	transientFields ...TransientField,
) (map[string]any, error) {
	// Make a copy of the original map to avoid modifying it directly.
	replaced := make(map[string]any, len(original))
	for key, value := range original {
		replaced[key] = value
	}

	// Apply the replacements one by one.
	for _, field := range transientFields {
		// See whether we should skip this replacement for the given path.
		skip, err := skipFile(path, field.FileRegex, field.FileRegexFullPath)
		if err != nil {
			return nil, fmt.Errorf("error checking file regex: %w", err)
		}
		if skip {
			continue
		}

		// Find all keys that match the field's key.
		originalKeys := make([]string, 0, len(replaced))
		for key := range replaced {
			originalKeys = append(originalKeys, key)
		}
		replacementKeys, replacedKeys := filterMatchingJPaths(
			originalKeys,
			field.Key,
		)
		if len(replacedKeys) == 0 {
			// No keys matched the field's key, so we skip this replacement.
			continue
		}

		// Keep the first value found for potential type based fallback
		// replacement.
		var firstValue any = replaced[replacedKeys[0]]

		// Remove the matched keys from the original map.
		for _, key := range replacedKeys {
			delete(replaced, key)
		}

		// Insert the replacements.
		for _, key := range replacementKeys {
			if field.Replacement != nil {
				replaced[key] = field.Replacement
				continue
			}

			// If the replacement is nil, we fall back to default stable values
			// (based on the type of the last value found).
			if stringValue, isString := firstValue.(string); isString {
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
			if _, isFloat := firstValue.(float64); isFloat {
				replaced[key] = StableFloat
				continue
			}
			if _, isInt := firstValue.(int); isInt {
				replaced[key] = StableInt
				continue
			}
			if _, isBool := firstValue.(bool); isBool {
				replaced[key] = StableBool
				continue
			}
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
		// See whether we should skip this rounding for the given path.
		skip, err := skipFile(path, field.FileRegex, field.FileRegexFullPath)
		if err != nil {
			return nil, fmt.Errorf("error checking file regex: %w", err)
		}
		if skip {
			continue
		}
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
