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
		slug:      slug,
		prefix:    prefix,
	}
}

// Connect connects a method with its implementation.
func Connect[T any](c *Connector, target *T, suffix ...string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if _, ok := c.connected[target]; ok {
		return
	}

	// get the calling function, get ok to make linter happy
	pc, _, _, ok := runtime.Caller(1)
	// we don't actually need ok, so noop
	_ = ok
	// get name of the calling function
	fullName := runtime.FuncForPC(pc).Name()
	trimmed := strings.TrimSuffix(fullName, "[...]")
	split := strings.Split(trimmed, ".")
	name := split[len(split)-1]

	// for methods that share more mapped names, such as NewMap, we use a suffix
	if len(suffix) == 1 {
		name += suffix[0]
	}

	// connect by convention
	plugin.Connect(c.slug, fmt.Sprintf("%s%s", c.prefix, name), target)
	c.connected[target] = struct{}{}
}
