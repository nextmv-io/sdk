// Package plugin provides functions for loading Hop plugins.
package plugin

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime"
	"sync"
)

//go:embed version.txt
var version string

// Connect a symbol in a plugin to a func target.
//
//    var fooFunc func()
//    plugin.Connect("sdk", "Foo", &func)
func Connect[T any](slug string, name string, target *T) {
	p, err := loadPlugin(slug)
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

var loaded = map[string]*plugin.Plugin{}

var mtx sync.Mutex

func loadPlugin(slug string) (*plugin.Plugin, error) {
	// Only load the plugin once. Then reuse the plugin pointer.
	if p, ok := loaded[slug]; ok {
		return p, nil
	}

	mtx.Lock()
	defer mtx.Unlock()

	if p, ok := loaded[slug]; ok {
		return p, nil
	}

	p, err := plugin.Open(pluginPath(slug))
	if err != nil {
		return nil, err
	}
	loaded[slug] = p

	return p, nil
}

func pluginPath(slug string) string {
	libraryPath := os.Getenv("NEXTMV_LIBRARY_PATH")
	if libraryPath == "" {
		libraryPath = "."
	}

	filename := fmt.Sprintf(
		"nextmv-%s-%s-%s-%s.so",
		slug,
		runtime.GOOS,
		runtime.GOARCH,
		version,
	)
	return filepath.Join(libraryPath, filename)
}
