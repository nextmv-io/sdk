package run

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

	plugin.Connect("HopRunRun", &runFunc)
}
