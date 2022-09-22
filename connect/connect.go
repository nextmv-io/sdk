// Package connect provides a Connector which allows to connect method
// definition with their implementations in plugins
package connect

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/nextmv-io/sdk/plugin"
)

// Connector connects methods with their implementations in plugins.
type Connector struct {
	connected    map[any]struct{}
	mtx          sync.Mutex
	slug, prefix string
}

// NewConnector creates a new Connector.
func NewConnector(slug, prefix string) *Connector {
	return &Connector{
		connected: make(map[any]struct{}),
		mtx:       sync.Mutex{},
		slug:      slug,
		prefix:    prefix,
	}
}

// Connect connects a method with its implementation.
func Connect[T any](c *Connector, target *T) {
	if _, ok := c.connected[target]; ok {
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if _, ok := connected[target]; ok {
		return
	}

	pc, _, _, ok := runtime.Caller(1)
	_ = ok
	fullName := runtime.FuncForPC(pc).Name()
	split := strings.Split(fullName, ".")
	name := split[len(split)-1]
	plugin.Connect(c.slug, fmt.Sprintf("%s%s", c.prefix, name), target)
	c.connected[target] = struct{}{}
}
