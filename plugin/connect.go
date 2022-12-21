// Package plugin provides functions for connecting plugins built from a
// private source.
package plugin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime"

	"github.com/nextmv-io/sdk"
)

// Connect a symbol in a plugin to a func target.
//
//	var fooFunc func()
//	plugin.Connect("sdk", "Foo", &func)
func Connect[T any](slug string, name string, target *T) {
	// the two locations plugins can be found in are the current working
	// directory and the nextmv library path
	paths := potentialPluginPaths(slug)
	pluginPath := ""
	for _, path := range paths {
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			pluginPath = path
			break
		}
	}
	if pluginPath == "" {
		fmt.Fprintf(os.Stderr,
			"could not find plugin %q in any of the paths %q\n\n", slug, paths)
		os.Exit(1)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading plugin %q\n\n", pluginPath)
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

func potentialPluginPaths(slug string) []string {
	// Get plugin filename we are looking for
	filename := fmt.Sprintf(
		"nextmv-%s-%s-%s-%s-%s%s.so",
		slug,
		sdk.VERSION,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		debug,
	)

	// Collect potential plugin paths
	paths := []string{}

	// Get Nextmv library path (use default if not set)
	libraryPath := os.Getenv("NEXTMV_LIBRARY_PATH")
	if libraryPath != "" {
		paths = append(paths, filepath.Join(libraryPath, filename))
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			paths = append(paths, filepath.Join(homeDir, ".nextmv", "lib", filename))
		}
	}

	// Get current working directory
	cwdPath, err := os.Getwd()
	if err == nil {
		paths = append(paths, filepath.Join(cwdPath, filename))
	}

	// Get binary directory (only use if not equal to cwd)
	binaryPath, err := os.Executable()
	if err == nil {
		binaryPath = filepath.Join(filepath.Dir(binaryPath), filename)
		if binaryPath != paths[len(paths)-1] {
			paths = append(paths, binaryPath)
		}
	}

	return paths
}
