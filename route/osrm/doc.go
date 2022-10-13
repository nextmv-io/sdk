// Package osrm provides a client for measuring distances and durations.
//
// An OSRM client requests distance and duration data from an OSRM server. It
// makes requests to construct a matrix measurer.
//
//	client := osrm.DefaultClient("http://localhost:5000", true)
//
// The client can construct a distance matrix, a duration matrix, or both.
//
//	dist := osrm.DistanceMatrix(client, points, 0)
//	dur := osrm.DurationMatrix(client, points, 0)
//	dist, dur, err := osrm.DistanceDurationMatrices(client, points, 0)
//
// These measures implement route.ByIndex.
package osrm
