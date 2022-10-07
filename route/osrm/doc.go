// © 2019-2022 nextmv.io inc. All rights reserved.
// nextmv.io, inc. CONFIDENTIAL
//
// This file includes unpublished proprietary source code of nextmv.io, inc.
// The copyright notice above does not evidence any actual or intended
// publication of such source code. Disclosure of this source code or any
// related proprietary information is strictly prohibited without the express
// written permission of nextmv.io, inc.

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
