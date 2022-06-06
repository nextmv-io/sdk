package context

import (
	"sync"

	"github.com/nextmv-io/sdk/hop/plugin"
)

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

	plugin.Connect("HopContextNewContext", &newContextFunc)

	// Declare variables
	plugin.Connect("HopContextDeclare", &declareFunc)
	plugin.Connect("HopContextDeclaredGet", &declaredGetFunc)
	plugin.Connect("HopContextDeclaredSet", &declaredSetFunc)

	// Conditions
	plugin.Connect("HopContextAnd", &andFunc)
	plugin.Connect("HopContextFalse", &falseFunc)
	plugin.Connect("HopContextNot", &notFunc)
	plugin.Connect("HopContextOr", &orFunc)
	plugin.Connect("HopContextTrue", &trueFunc)
	plugin.Connect("HopContextXor", &xorFunc)
}
