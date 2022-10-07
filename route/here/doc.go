// © 2019-2022 nextmv.io inc. All rights reserved.
// nextmv.io, inc. CONFIDENTIAL
//
// This file includes unpublished proprietary source code of nextmv.io, inc.
// The copyright notice above does not evidence any actual or intended
// publication of such source code. Disclosure of this source code or any
// related proprietary information is strictly prohibited without the express
// written permission of nextmv.io, inc.

// Package here provides a client for measuring distances and durations.
//
// A HERE client requests distance and duration data using HERE Maps API. It
// makes requests to construct a matrix measure.
//
//	client := here.NewClient("<API_KEY>")
//
// The client can construct a distance matrix, a duration matrix, or both.
//
//	distances, err := client.DistanceMatrix(ctx, points)
//	durations, err := client.DurationMatrix(ctx, points)
//	distances, durations, err := client.DistanceDurationMatrices(ctx, points)
//
// Each of these functions will use a synchronous request flow if the number
// of points requested is below HERE's size limit for synchronous API calls -
// otherwise, an asynchronous flow will be used. The functions all take a
// context which can be used to cancel the request flow while it is in progress.
//
// These measures implement route.ByIndex.
//
// These matrix-generating functions can also take one or more options
// that allow you to configure the routes that will be included in the matrices.
// For example, you can set a specific departure time to use when factoring
// in traffic time to the route durations:
//
//	durations, err := client.DurationMatrix(
//	    ctx,
//	    points,
//	    here.WithDepartureTime(time.Date(2021, 12, 10, 8, 30, 0, 0, loc)),
//	)
//
// Or, you can configure a truck profile:
//
//	distances, err := client.DistanceMatrix(
//	    ctx,
//	    points,
//	    here.WithTransportMode(here.Truck),
//	    here.WithTruck(
//	        Type: here.TruckTypeTractor,
//	        TrailerCount: 2,
//	        ShippedHazardousGoods: []here.HazardousGood{here.Poison},
//	    ),
//	)
package here
