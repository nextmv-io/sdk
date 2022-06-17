// Package plugin provides functions for loading Hop plugins.
package plugin

import (
	"errors"
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
//    plugin.Connect("sdk", "Foo", &func)
func Connect[T any](slug string, name string, target *T) {
	path := pluginPath(slug)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "plugin file %q does not exist\n\n", path)
		os.Exit(1)
	}

	p, err := loadPlugin(slug, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading plugin %q\n\n", path)
		panic(err)
	}

	sym, err := p.Lookup(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting symbol %q\n\n", name)
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

func loadPlugin(slug, path string) (*plugin.Plugin, error) {
	// Only load the plugin once. Then reuse the plugin pointer.
	if p, ok := loaded[slug]; ok {
		return p, nil
	}

	mtx.Lock()
	defer mtx.Unlock()

	if p, ok := loaded[slug]; ok {
		return p, nil
	}

	p, err := plugin.Open(path)
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
		"nextmv-%s-%s-%s-%s-%s.so",
		slug,
		sdk.VERSION,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
	return filepath.Join(libraryPath, filename)
}
