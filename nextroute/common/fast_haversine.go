package common

import "math"

func NewFastHaversine(lat float64) FastHaversine {
	const RE = 6378.137          // equatorial radius
	const FE = 1 / 298.257223563 // flattening

	const E2 = FE * (2 - FE)
	const RAD = math.Pi / 180
	const m = RAD * RE * 1000
	var coslat = math.Cos(lat * RAD)
	var w2 = 1 / (1 - E2*(1-coslat*coslat))
	var w = math.Sqrt(w2)

	return FastHaversine{
		kx: m * w * coslat,        // based on normal radius of curvature
		ky: m * w * w2 * (1 - E2), // based on meridonal radius of curvature
	}
}

type FastHaversine struct {
	kx float64
	ky float64
}

func wrap(deg float64) float64 {
	for deg < -180 {
		deg += 360
	}
	for deg > 180 {
		deg -= 360
	}
	return deg
}

func (f FastHaversine) Distance(from, to Location) float64 {
	var dx = wrap(from.Longitude()-to.Longitude()) * f.kx
	var dy = (from.Latitude() - to.Latitude()) * f.ky
	return math.Sqrt(dx*dx + dy*dy)
}
