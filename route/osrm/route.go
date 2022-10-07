// © 2019-2022 nextmv.io inc. All rights reserved.
// nextmv.io, inc. CONFIDENTIAL
//
// This file includes unpublished proprietary source code of nextmv.io, inc.
// The copyright notice above does not evidence any actual or intended
// publication of such source code. Disclosure of this source code or any
// related proprietary information is strictly prohibited without the express
// written permission of nextmv.io, inc.

package osrm

import "github.com/nextmv-io/sdk/route"

// Polyline requests polylines for the given points. The first parameter returns
// a polyline from start to end and the second parameter returns a list of
// polylines, one per leg.
func Polyline(
	c Client, points []route.Point,
) (string, []string, error) {
	polyline, legLines, err := c.Polyline(points)
	if err != nil {
		return "", []string{}, err
	}

	return polyline, legLines, nil
}
