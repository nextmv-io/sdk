package store

import "reflect"

// A Map stores key-value pairs in a Store.
type Map[K Key, V any] interface {
	// Delete a key from the map.
	Delete(K) Change
	// Get an index of a vector.
	Get(Store, K) (V, bool)
	// Len returns the number of keys in a map,
	Len(Store) int
	// Map representation that is mutable.
	Map(Store) map[K]V
	// Set a key to a value.
	Set(K, V) Change
}

// A Key for a Map.
type Key interface{ int | string }

// NewMap returns a new Map and stores it in a Store.
func NewMap[K Key, V any](s Store) Map[K, V] {
	p := mapProxy[K, V]{}

	var k K
	if isInt(k) {
		p.mapInt = newMapIntFunc(s)
	} else {
		p.mapString = newMapStringFunc(s)
	}

	return p
}

// Since type constraints cannot cross the plugin boundary bidirectionally, we
// simulate a union type.
type mapProxy[K Key, V any] struct {
	mapInt    Map[int, any]
	mapString Map[string, any]
}

// Implements Map

func (m mapProxy[K, V]) Delete(key K) Change {
	k := any(key)
	if m.mapInt != nil {
		return m.mapInt.Delete(k.(int))
	}
	return m.mapString.Delete(k.(string))
}

func (m mapProxy[K, V]) Get(s Store, key K) (V, bool) {
	k := any(key)

	if m.mapInt != nil {
		value, ok := m.mapInt.Get(s, k.(int))
		return value.(V), ok
	}

	value, ok := m.mapString.Get(s, k.(string))
	return value.(V), ok
}

func (m mapProxy[K, V]) Len(s Store) int {
	if m.mapInt != nil {
		return m.mapInt.Len(s)
	}
	return m.mapString.Len(s)
}

func (m mapProxy[K, V]) Map(s Store) map[K]V {
	mapKV := map[K]V{}

	if m.mapInt != nil {
		for key, value := range m.mapInt.Map(s) {
			mapKV[any(key).(K)] = value.(V)
		}
		return mapKV
	}

	for key, value := range m.mapString.Map(s) {
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

var (
	newMapIntFunc    func(Store) Map[int, any]
	newMapStringFunc func(Store) Map[string, any]
)
