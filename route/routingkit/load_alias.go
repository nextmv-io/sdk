package routingkit

import (
	r "github.com/nextmv-io/sdk/measure/routingkit"
)

// ByPointLoader can be embedded in schema structs and unmarshals a ByPoint JSON
// object into the appropriate implementation, including a routingkit.ByPoint.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
type ByPointLoader = r.ByPointLoader

// ByIndexLoader can be embedded in schema structs and unmarshals a ByIndex JSON
// object into the appropriate implementation, including a routingkit.ByIndex.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
type ByIndexLoader = r.ByIndexLoader

// ProfileLoader can be embedded in schema structs and unmarshals a
// routingkit.Profile JSON object into the appropriate implementation.
//
// Deprecated: This package is deprecated and will be removed in a future.
// Use [github.com/nextmv-io/sdk/measure/routingkit] instead.
type ProfileLoader = r.ProfileLoader
