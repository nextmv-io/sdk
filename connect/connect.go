package connect

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/nextmv-io/sdk/plugin"
)

// Connector connects method definitions with plugins.
type Connector interface {
	Connect(any)
}

type connector struct {
	connected    map[any]struct{}
	mtx          *sync.Mutex
	slug, prefix string
}

// NewConnector creates a new Connector.
func NewConnector(slug, prefix string) Connector {
	return connector{
		connected: make(map[any]struct{}),
		mtx:       &sync.Mutex{},
		slug:      slug,
		prefix:    prefix,
	}
}

// Connect connects a method with its implementation.
func (c connector) Connect(target any) {
	if _, ok := c.connected[target]; ok {
		return
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if _, ok := c.connected[target]; ok {
		return
	}

	pc, _, _, ok := runtime.Caller(1)
	_ = ok
	fullName := runtime.FuncForPC(pc).Name()
	split := strings.Split(fullName, ".")
	name := split[len(split)-1]
	plugin.Connect(c.slug, fmt.Sprintf("%s%s", c.prefix, name), &target)
	c.connected[target] = struct{}{}
}
