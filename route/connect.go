package route

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/nextmv-io/sdk/plugin"
)

const slug = "sdk"

var connected = map[any]struct{}{}

var mtx sync.Mutex

func connect[T any](target *T) {
	if _, ok := connected[target]; ok {
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
	plugin.Connect(slug, fmt.Sprintf("Route%s", name), target)
	connected[target] = struct{}{}
}
