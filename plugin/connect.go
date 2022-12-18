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

// Token is used to permanently inject a Token into an application.
var Token string

// Connect a symbol in a plugin to a func target.
//
//	var fooFunc func()
//	plugin.Connect("sdk", "Foo", &func)
func Connect[T any](slug string, name string, target *T) {
	// If a token is injected, use it
	if Token != "" {
		err := os.Setenv("NEXTMV_BAKED_IN_TOKEN", Token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err setting token: %v", err)
		}
	}

	// the two locations plugins can be found in are the current working
	// directory and the nextmv library path
	paths, err := potentialPluginPaths(slug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err getting plugin paths: %v", err)
	}
	pluginPath := ""
	for _, path := range paths {
		if _, err = os.Stat(path); !errors.Is(err, os.ErrNotExist) {
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

func potentialPluginPaths(slug string) ([]string, error) {
	// Get Nextmv library path
	libraryPath := os.Getenv("NEXTMV_LIBRARY_PATH")
	if libraryPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not fetch user home dir: %v", err)
		}
		libraryPath = filepath.Join(homeDir, ".nextmv", "lib")
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not fetch current working directory: %v", err)
	}

	// Get binary directory
	binaryPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not fetch binary path: %v", err)
	}
	binaryDir := filepath.Dir(binaryPath)

	// Assemble potential plugin paths
	filename := fmt.Sprintf(
		"nextmv-%s-%s-%s-%s-%s%s.so",
		slug,
		sdk.VERSION,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		debug,
	)
	paths := []string{
		filepath.Join(cwd, filename),
		filepath.Join(binaryDir, filename),
		filepath.Join(libraryPath, filename),
	}
	return paths, nil
}
