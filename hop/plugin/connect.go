package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime"
	"sync"

	"github.com/nextmv-io/sdk"
)

// Connect a symbol in a plugin to a func target.
//
//    var fooFunc func()
//    plugin.Connect("Foo", &func)
func Connect[T any](name string, target *T) {
	p, err := loadPlugin()
	if err != nil {
		panic(err)
	}

	sym, err := p.Lookup(name)
	if err != nil {
		panic(err)
	}

	// Names in the plugin are associated with pointers to functions.
	// Thus we cannot: *target = sym(T)
	*target = reflect.ValueOf(sym). // *func(...) as reflect.Value
					Elem().         // dereferences to func(...)
					Interface().(T) // any.(func(...))
}

var loadedPlugin *plugin.Plugin

var mtx sync.Mutex

func loadPlugin() (*plugin.Plugin, error) {
	// Only load the plugin once. Then reuse the plugin pointer.
	if loadedPlugin != nil {
		return loadedPlugin, nil
	}

	mtx.Lock()
	defer mtx.Unlock()

	if loadedPlugin != nil {
		return loadedPlugin, nil
	}

	p, err := plugin.Open(pluginPath())
	if err != nil {
		return nil, err
	}
	loadedPlugin = p

	return loadedPlugin, nil
}

func pluginPath() string {
	libraryPath := os.Getenv("NEXTMV_LIBRARY_PATH")
	if libraryPath == "" {
		libraryPath = "."
	}

	filename := fmt.Sprintf(
		"nextmv-sdk-%s-%s-%s.so",
		runtime.GOOS,
		runtime.GOARCH,
		sdk.VERSION,
	)
	return filepath.Join(libraryPath, filename)
}
