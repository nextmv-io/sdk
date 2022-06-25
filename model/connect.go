package model

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

	// Domain
	plugin.Connect(slug, "ModelNewDomain", &newDomainFunc)
	plugin.Connect(slug, "ModelSingleton", &singletonFunc)
	plugin.Connect(slug, "ModelMultiple", &multipleFunc)

	// Domains
	plugin.Connect(slug, "ModelNewDomains", &newDomainsFunc)
	plugin.Connect(slug, "ModelRepeat", &repeatFunc)

	// Range
	plugin.Connect(slug, "ModelNewRange", &newRangeFunc)
}
