package run

import (
	"os"
	"sync"

	"github.com/nextmv-io/sdk/plugin"
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

	slug := os.Getenv("NEXTMV_RUNNER")
	if slug == "" {
		slug = "cli"
	}
	slug = "run-" + slug

	plugin.Connect(slug, "RunNew", &newFunc)
}
