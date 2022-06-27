package route

// NewRouter returns a router interface. It receives a set of stops that must be
// serviced by a fleet of vehicles and a list of options. When an option is
// applied, an error is returned if there are validation issues. The router is
// composable, meaning that several options may be used or none at all. The
// options, unless otherwise noted, can be used independently of each other.
func NewRouter(stops []Stop,
	vehicles []string,
	opts ...Option,
) (Router, error) {
	connect()
	return newRouterFunc(stops, vehicles, opts...)
}

var newRouterFunc func(
	[]Stop,
	[]string, ...Option,
) (Router, error)
