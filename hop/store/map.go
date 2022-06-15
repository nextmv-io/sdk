package store

import (
	"reflect"

	"github.com/nextmv-io/sdk/hop/store/types"
)

// Map returns a new Map and stores it in a Store.
func Map[K types.Key, V any](s types.Store) types.Map[K, V] {
	p := mapProxy[K, V]{}

	var k K
	if isInt(k) {
		p.mapInt = mapIntFunc(s)
	} else {
		p.mapString = mapStringFunc(s)
	}

	return p
}

// Since type constraints cannot cross the plugin boundary bidirectionally, we
// simulate a union type.
type mapProxy[K types.Key, V any] struct {
	mapInt    types.Map[int, any]
	mapString types.Map[string, any]
}

// Implements types.Map.

func (m mapProxy[K, V]) Delete(key K) types.Change {
	k := any(key)
	if m.mapInt != nil {
		return m.mapInt.Delete(k.(int))
	}
	return m.mapString.Delete(k.(string))
}

func (m mapProxy[K, V]) Get(s types.Store, key K) (V, bool) {
	k := any(key)

	if m.mapInt != nil {
		value, ok := m.mapInt.Get(s, k.(int))
		return value.(V), ok
	}

	value, ok := m.mapString.Get(s, k.(string))
	return value.(V), ok
}

func (m mapProxy[K, V]) Len(s types.Store) int {
	if m.mapInt != nil {
		return m.mapInt.Len(s)
	}
	return m.mapString.Len(s)
}

func (m mapProxy[K, V]) Map(s types.Store) map[K]V {
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

func (m mapProxy[K, V]) Set(key K, value V) types.Change {
	k := any(key)
	if m.mapInt != nil {
		return m.mapInt.Set(k.(int), value)
	}
	return m.mapString.Set(k.(string), value)
}

func isInt[K types.Key](k K) bool {
	// We can't type switch on a type constraint (e.g. "switch k.(type)"), but
	// we can use reflection.
	return reflect.TypeOf(k) == reflect.TypeOf(0)
}

var (
	mapIntFunc    func(types.Store) types.Map[int, any]
	mapStringFunc func(types.Store) types.Map[string, any]
)
