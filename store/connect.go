package store

import (
	"sync"

	"github.com/nextmv-io/sdk/plugin"
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

	plugin.Connect(slug, "StoreNew", &newFunc)

	// Declare variables
	plugin.Connect(slug, "StoreNewVar", &newVarFunc)

	// Conditions
	plugin.Connect(slug, "StoreAnd", &andFunc)
	plugin.Connect(slug, "StoreFalse", &falseFunc)
	plugin.Connect(slug, "StoreNot", &notFunc)
	plugin.Connect(slug, "StoreOr", &orFunc)
	plugin.Connect(slug, "StoreTrue", &trueFunc)
	plugin.Connect(slug, "StoreXor", &xorFunc)

	// State generation
	plugin.Connect(slug, "StoreIf", &ifFunc)
	plugin.Connect(slug, "StoreScope", &scopeFunc)

	// Collections
	plugin.Connect(slug, "StoreMapInt", &mapIntFunc)
	plugin.Connect(slug, "StoreMapString", &mapStringFunc)
	plugin.Connect(slug, "StoreNewSlice", &newSliceFunc)

	// Solver
	plugin.Connect(slug, "StoreDefaultOptions", &defaultOptionsFunc)
}
