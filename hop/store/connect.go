package store

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

	plugin.Connect(slug, "HopStoreNewStore", &newStoreFunc)

	// Declare variables
	plugin.Connect(slug, "HopStoreVar", &varFunc)

	// Conditions
	plugin.Connect(slug, "HopStoreAnd", &andFunc)
	plugin.Connect(slug, "HopStoreFalse", &falseFunc)
	plugin.Connect(slug, "HopStoreNot", &notFunc)
	plugin.Connect(slug, "HopStoreOr", &orFunc)
	plugin.Connect(slug, "HopStoreTrue", &trueFunc)
	plugin.Connect(slug, "HopStoreXor", &xorFunc)

	// State generation
	plugin.Connect(slug, "HopStoreIf", &ifFunc)
	plugin.Connect(slug, "HopStoreScope", &scopeFunc)

	// Collections
	plugin.Connect(slug, "HopStoreNewMapInt", &newMapIntFunc)
	plugin.Connect(slug, "HopStoreNewMapString", &newMapStringFunc)
	plugin.Connect(slug, "HopStoreNewSlice", &newSliceFunc)
}
