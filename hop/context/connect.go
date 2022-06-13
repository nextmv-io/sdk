package context

import (
	"sync"

	"github.com/nextmv-io/sdk/hop/plugin"
)

const slug = "sdk"

var connected bool

var mtx sync.Mutex

func connect() {
	if connected {
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if connected {
		return
	}
	connected = true

	// Create a new context
	plugin.Connect(slug, "HopContextNewContext", &newContextFunc)

	// Declare variables
	plugin.Connect(slug, "HopContextDeclare", &declareFunc)

	// Conditions
	plugin.Connect(slug, "HopContextAnd", &andFunc)
	plugin.Connect(slug, "HopContextFalse", &falseFunc)
	plugin.Connect(slug, "HopContextNot", &notFunc)
	plugin.Connect(slug, "HopContextOr", &orFunc)
	plugin.Connect(slug, "HopContextTrue", &trueFunc)
	plugin.Connect(slug, "HopContextXor", &xorFunc)

	// State generation
	plugin.Connect(slug, "HopContextIf", &ifFunc)
	plugin.Connect(slug, "HopContextScope", &scopeFunc)

	// Collections
	plugin.Connect(slug, "HopContextNewMapInt", &newMapIntFunc)
	plugin.Connect(slug, "HopContextNewMapString", &newMapStringFunc)
	plugin.Connect(slug, "HopContextNewSlice", &newSliceFunc)
}
