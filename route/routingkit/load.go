package routingkit

import (
	r "github.com/nextmv-io/sdk/measure/routingkit"
)

// ByPointLoader can be embedded in schema structs and unmarshals a ByPoint JSON
// object into the appropriate implementation, including a routingkit.ByPoint.
type ByPointLoader = r.ByPointLoader

// ByIndexLoader can be embedded in schema structs and unmarshals a ByIndex JSON
// object into the appropriate implementation, including a routingkit.ByIndex.
type ByIndexLoader = r.ByIndexLoader

// ProfileLoader can be embedded in schema structs and unmarshals a
// routingkit.Profile JSON object into the appropriate implementation.
type ProfileLoader = r.ProfileLoader
