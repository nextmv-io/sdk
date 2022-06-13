package context

import "reflect"

// A Map stores key-value pairs.
type Map[K Key, V any] interface {
	// Delete a key from the map.
	Delete(K) Change
	// Get an index of a vector.
	Get(Context, K) (V, bool)
	// Len returns the number of keys in a map,
	Len(Context) int
	// Map representation that is mutable.
	Map(Context) map[K]V
	// Set a key to a value.
	Set(K, V) Change
}

// A Key for a map.
type Key interface{ int | string }

// NewMap returns a new Map.
func NewMap[K Key, V any](ctx Context) Map[K, V] {
	p := mapProxy[K, V]{}

	var k K
	if isInt(k) {
		p.mapInt = newMapIntFunc(ctx)
	} else {
		p.mapString = newMapStringFunc(ctx)
	}

	return p
}

// Since type constraints cannot cross the plugin boundary bidirectionally, we
// simulate a union type.
type mapProxy[K Key, V any] struct {
	mapInt    Map[int, any]
	mapString Map[string, any]
}

func (m mapProxy[K, V]) Delete(key K) Change {
	k := any(key)
	if m.mapInt != nil {
		return m.mapInt.Delete(k.(int))
	}
	return m.mapString.Delete(k.(string))
}

func (m mapProxy[K, V]) Get(ctx Context, key K) (V, bool) {
	k := any(key)

	if m.mapInt != nil {
		value, ok := m.mapInt.Get(ctx, k.(int))
		return value.(V), ok
	}

	value, ok := m.mapString.Get(ctx, k.(string))
	return value.(V), ok
}

func (m mapProxy[K, V]) Len(ctx Context) int {
	if m.mapInt != nil {
		return m.mapInt.Len(ctx)
	}
	return m.mapString.Len(ctx)
}

func (m mapProxy[K, V]) Map(ctx Context) map[K]V {
	mapKV := map[K]V{}

	if m.mapInt != nil {
		for key, value := range m.mapInt.Map(ctx) {
			mapKV[any(key).(K)] = value.(V)
		}
		return mapKV
	}

	for key, value := range m.mapString.Map(ctx) {
		mapKV[any(key).(K)] = value.(V)
	}
	return mapKV
}

func (m mapProxy[K, V]) Set(key K, value V) Change {
	k := any(key)
	if m.mapInt != nil {
		return m.mapInt.Set(k.(int), value)
	}
	return m.mapString.Set(k.(string), value)
}

func isInt[K Key](k K) bool {
	// We can't type switch on a type constraint (e.g. "switch k.(type)"), but
	// we can use reflection.
	return reflect.TypeOf(k) == reflect.TypeOf(0)
}

var newMapIntFunc func(Context) Map[int, any]

var newMapStringFunc func(Context) Map[string, any]
