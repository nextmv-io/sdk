package context

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

const header = "This early release software is provided for evaluation and " +
	"feedback only.\n© 2019-2022 nextmv.io inc. All rights reserved."

var loaded bool

var mtx sync.Mutex

func load() {
	if loaded {
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if loaded {
		return
	}
	loaded = true

	path := pluginPath()
	p, err := plugin.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to load plugin file %q\n\n%w", path, err))
	}

	if _, err := fmt.Fprintln(os.Stderr, header); err != nil {
		panic(err)
	}

	connect(p, "NewContext", &newContextFunc)

	// Declare variables
	connect(p, "Declare", &declareFunc)
	connect(p, "DeclaredGet", &declaredGetFunc)
	connect(p, "DeclaredSet", &declaredSetFunc)

	// Conditions
	connect(p, "And", &andFunc)
	connect(p, "False", &falseFunc)
	connect(p, "Not", &notFunc)
	connect(p, "Or", &orFunc)
	connect(p, "True", &trueFunc)
	connect(p, "Xor", &xorFunc)
}

func connect[T any](p *plugin.Plugin, name string, target *T) {
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

func pluginPath() string {
	libraryPath := os.Getenv("NEXTMV_LIBRARY_PATH")
	if libraryPath == "" {
		libraryPath = "."
	}

	fileName := fmt.Sprintf(
		"nextmv-sdk-%s-%s-%s.so",
		runtime.GOOS,
		runtime.GOARCH,
		sdk.VERSION,
	)
	return filepath.Join(libraryPath, fileName)
}
