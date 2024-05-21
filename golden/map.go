package golden

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// replaceTransient replaces all the values in a map whose key is contained in
// the variadic list of transientFields. The replaced value has a stable value
// according to the data type.
func replaceTransient(
	original map[string]any,
	transientFields ...TransientField,
) map[string]any {
	transientLookup := map[string]any{}
	for _, field := range transientFields {
		transientLookup[field.Key] = field.Replacement
	}

	replaced := map[string]any{}
	for key, value := range original {
		replaced[key] = value
		replacement, isTransient := transientLookup[key]
		if !isTransient {
			continue
		}

		if replacement != nil {
			replaced[key] = replacement
			continue
		}

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

	return replaced
}

/*
flatten takes a nested map and flattens it into a single level map. The
flattening roughly follows the [JSONPath] standard. Please see test function to
understand how the flattened output looks like. Here is an example that may
fall out of date, so be careful:

If this is the nested input:

	map[string]any{
	    "a": "foo",
	    "b": []any{
	        map[string]any{
	            "c": "bar",
	            "d": []any{
	                map[string]any{
	                    "e": 2,
	                },
	                true,
	            },
	        },
	        map[string]any{
	            "c": "baz",
	            "d": []any{
	                map[string]any{
	                    "e": 3,
	                },
	                false,
	            },
	        },
	    },
	}

You can expect this flattened output:

	 map[string]any{
	    ".a":           "foo",
	    ".b[0].c":      "bar",
	    ".b[0].d[0].e": 2,
	    ".b[0].d[1]":   true,
	    ".b[1].c":      "baz",
	    ".b[1].d[0].e": 3,
	    ".b[1].d[1]":   false,
	}

[JSONPath]: https://goessner.net/articles/JsonPath/
*/
func flatten(nested map[string]any) map[string]any {
	flattened := map[string]any{}
	for childKey, childValue := range nested {
		setChildren(flattened, childKey, childValue)
	}

	return flattened
}

// setChildren is a helper function for flatten. It is invoked recursively on a
// child value. If the child is not a map or a slice, then the value is simply
// set on the flattened map. If the child is a map or a slice, then the
// function is invoked recursively on the child's values, until a
// non-map-non-slice value is hit.
func setChildren(flattened map[string]any, parentKey string, parentValue any) {
	newKey := fmt.Sprintf(".%s", parentKey)
	if reflect.TypeOf(parentValue) == nil {
		flattened[newKey] = parentValue
		return
	}

	if reflect.TypeOf(parentValue).Kind() == reflect.Map {
		children := parentValue.(map[string]any)
		for childKey, childValue := range children {
			newKey = fmt.Sprintf("%s.%s", parentKey, childKey)
			setChildren(flattened, newKey, childValue)
		}
		return
	}

	if reflect.TypeOf(parentValue).Kind() == reflect.Slice {
		children := parentValue.([]any)
		if len(children) == 0 {
			flattened[newKey] = children
			return
		}

		for childIndex, childValue := range children {
			newKey = fmt.Sprintf("%s[%v]", parentKey, childIndex)
			setChildren(flattened, newKey, childValue)
		}
		return
	}

	flattened[newKey] = parentValue
}

/*
nest takes a flattened map and nests it into a multi-level map. The flattened
map roughly follows the [JSONPath] standard. Please see test function to
understand how the nested output looks like. Here is an example that may fall
out of date, so be careful:

If this is the flattened input:

	map[string]any{
		".a":           "foo",
		".b[0].c":      "bar",
		".b[0].d[0].e": 2,
		".b[0].d[1]":   true,
		".b[1].c":      "baz",
		".b[1].d[0].e": 3,
		".b[1].d[1]":   false,
	}

You can expect this nested output:

	map[string]any{
		"a": "foo",
		"b": []any{
			map[string]any{
				"c": "bar",
				"d": []any{
					map[string]any{
						"e": 2,
					},
					true,
				},
			},
			map[string]any{
				"c": "baz",
				"d": []any{
					map[string]any{
						"e": 3,
					},
					false,
				},
			},
		},
	}

[JSONPath]: https://goessner.net/articles/JsonPath/
*/
func nest(flattened map[string]any) (map[string]any, error) {
	// First, convert the flat map to a nested map. Then reshape the map into a
	// slice where appropriate.
	const magicSliceKey = "isSlice"
	nested := make(map[string]any)
	for key, value := range flattened {
		p, err := pathFrom(key)
		if err != nil {
			return nil, err
		}

		current := nested
		for i, k := range p {
			key := k.Key()
			if k.IsSlice() {
				current[magicSliceKey] = true
			}

			isLast := i == len(p)-1
			if isLast {
				current[key] = value
				break
			}

			if current[key] == nil {
				current[key] = make(map[string]any)
			}

			current = current[key].(map[string]any)
		}
	}

	// Convert maps to slices where appropriate using non recursive breadth
	// first search.
	queue := []map[string]any{nested}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for k, v := range current {
			m, ok := v.(map[string]any)
			if !ok {
				// Not a map, we reached the end of the tree.
				continue
			}

			if m[magicSliceKey] == nil {
				// Just a normal map, enqueue.
				queue = append(queue, m)
				continue
			}

			// A map that needs to be converted to a slice.
			delete(m, magicSliceKey)
			slice, err := toSlice(m)
			if err != nil {
				return nil, err
			}

			for _, x := range slice {
				if _, ok := x.(map[string]any); ok {
					// Enqueue all maps in the slice.
					queue = append(queue, x.(map[string]any))
				}
			}
			current[k] = slice
		}
	}
	return nested, nil
}

func toSlice(x map[string]any) ([]any, error) {
	slice := make([]any, len(x))
	for k, v := range x {
		idx, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}

		if idx >= len(slice) || idx < 0 {
			return nil, fmt.Errorf("index %d out of bounds", idx)
		}

		slice[idx] = v
	}
	return slice, nil
}

type pathKey struct {
	name  string
	index int
}

func (p pathKey) IsSlice() bool {
	return p.index != -1
}

func (p pathKey) Key() string {
	if p.IsSlice() {
		return strconv.Itoa(p.index)
	}
	return p.name
}

type path []pathKey

func pathFrom(key string) (path, error) {
	split := strings.Split(key[1:], ".")
	p := make(path, 0, len(split))
	for _, s := range split {
		stops, err := pathKeysFrom(s)
		if err != nil {
			return path{}, err
		}

		p = append(p, stops...)
	}

	return p, nil
}

func pathKeysFrom(key string) ([]pathKey, error) {
	if strings.Contains(key, "[") {
		start := strings.Index(key, "[")
		end := strings.Index(key, "]")
		index, err := strconv.Atoi(key[start+1 : end])
		if err != nil {
			return []pathKey{}, err
		}

		return []pathKey{
			{name: key[:start], index: -1},
			{index: index},
		}, nil
	}

	return []pathKey{{name: key, index: -1}}, nil
}
