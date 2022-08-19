package mip

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

	plugin.Connect(slug, "MIPNewSolveOptions", &newSolveOptions)
	plugin.Connect(slug, "MIPNewDefinition", &newDefinition)
	plugin.Connect(slug, "MIPNewSolver", &newSolver)
}
