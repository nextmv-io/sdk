package model

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

	// Domain
	plugin.Connect(slug, "HopModelDomain", &domainFunc)
	plugin.Connect(slug, "HopModelSingleton", &singletonFunc)
	plugin.Connect(slug, "HopModelMultiple", &multipleFunc)

	// Domains
	plugin.Connect(slug, "HopModelDomains", &domainsFunc)
	plugin.Connect(slug, "HopModelRepeat", &repeatFunc)

	// Range
	plugin.Connect(slug, "HopModelRange", &rangeFunc)
}
