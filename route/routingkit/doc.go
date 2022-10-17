// Package routingkit provides measures that calculate the cost of travelling
// between points on road networks.
//
// The `routingkit.ByPoint` constructor allow you to create a route.ByPoint
// that finds the shortest path between two points in terms of distance
// travelled:
//
//	byPoint, err := routingkit.ByPoint(
//		osmFile, 	// path to an .osm.pbf file
//		1000, 	 	// maximum distance to snap points to
//		1<<30, 	 	// use max 1GB to cache point distances
//		routingkit.Car, // limit to roads accessible by car
//		fallbackMeasure, // used when no route is found between points
//	)
//
// `routingkit.DurationByPoint` constructs a `route.ByPoint` that finds the
// shortest path in terms of travel time. (Only car travel times are supported)
//
//	byPoint, err := routingkit.DurationByPoint(
//		osmFile,
//		1000,
//		1<<30,
//		fallbackMeasure
//	)
//
// Finally, `routingkit.Matrix` and `routingkit.DurationMatrix` construct a
// `route.ByIndex` containing pre-built matrices of all point-to-point costs,
// using the relevant units.
//
//	distanceByIndex, err := routingkit.Matrix(
//		osmFile,
//		1000,
//		srcs,
//		dests,
//		routingkit.Car,
//		fallbackMeasure
//	)
//	durationByIndex, err := routingkit.DurationMatrix(
//		osmFile,
//		1000,
//		srcs,
//		dests,
//		fallbackMeasure
//	)
package routingkit
