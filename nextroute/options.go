package nextroute

import (
	"github.com/nextmv-io/sdk/nextroute/schema"
)

// An Option configures a model.
type Option func(Model, schema.Input) error
