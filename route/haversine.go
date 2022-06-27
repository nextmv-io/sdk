package route

// HaversineByPoint estimates meters connecting two points along the surface
// of the earth.
func HaversineByPoint() ByPoint {
	connect()
	return haversineByPointFunc()
}

var haversineByPointFunc func() ByPoint
