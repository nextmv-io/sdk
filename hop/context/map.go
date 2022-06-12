package context

// // A Map stores key-value pairs.
// type Map[K comparable, V any] interface {
// 	// Data underlying a map.
// 	Data(Context) map[K]V
// 	// Delete a key from the map.
// 	Delete(K) Change
// 	// Get an index of a vector.
// 	Get(Context, K) (V, bool)
// 	// Len returns the number of keys in a map,
// 	Len(Context) int
// 	// Set a key to a value.
// 	Set(K, V) Change
// }

// // NewMap returns a new Map.
// func NewMap[K comparable, V any](ctx Context) Map[K, V] {
// 	return mapProxy[K, V]{_map: newMapFunc[K](ctx)}
// }

// type mapProxy[K comparable, V any] struct {
// 	_map Map[K, any] // "map" was taken
// }

// func (m mapProxy[K, V]) Data(ctx Context) map[K]V {
// 	mapAny := m._map.Data(ctx)
// 	mapKV := make(map[K]V, len(mapAny))
// 	for k, v := range mapAny {
// 		mapKV[k] = v.(V)
// 	}
// 	return mapKV
// }

// func (m mapProxy[K, V]) Delete(key K) Change {
// 	return m._map.Delete(key)
// }

// func (m mapProxy[K, V]) Get(ctx Context, key K) (V, bool) {
// 	v, ok := m._map.Get(ctx, key)
// 	return v.(V), ok
// }

// func (m mapProxy[K, V]) Len(ctx Context) int {
// 	return m._map.Len(ctx)
// }

// func (m mapProxy[K, V]) Set(key K, value V) Change {
// 	return m._map.Set(key, value)
// }

// var newMapFunc func(Context) Map[comparable, any]
