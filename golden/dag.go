package golden

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

// DagTestCase represents a test case in a directed acyclic graph (DAG) test.
type DagTestCase struct {
	Name   string
	Needs  []string
	Config *ScriptConfig
	Path   string
}

// DagTest runs a set of test cases in topological order.
// Each test case is a ScriptTest, and the test cases are connected by their
// dependencies. If a test case has dependencies, it will only be run after all
// of its dependencies have been run.
//
// Sample usage:
//
//	cases := []golden.DagTestCase{
//	  {
//	    name:   "app-create",
//	    needs:  []string{},
//	    config: ScriptConfig{ /**/ },
//	    path:   "app-create",
//	  },
//	  {
//	    name:   "app-push",
//	    needs:  []string{"app-create"},
//	    config: ScriptConfig{ /**/ },
//	    path:   "app-push",
//	  },
//	}
//	golden.DagTest(t, cases)
func DagTest(t *testing.T, cases []DagTestCase) {
	t.Parallel()
	err := validate(cases)
	if err != nil {
		t.Fatal(err)
	}

	open := cases
	done := make(map[string]bool)

	for len(open) > 0 {
		// Pick the first case from the open list that has all its needs met.
		next := make([]DagTestCase, 0)
		for _, c := range open {
			ready := true
			for _, need := range c.Needs {
				if !done[need] {
					ready = false
					break
				}
			}
			if ready {
				next = append(next, c)
			}
		}

		// If we didn't find a case to run, we have a cycle.
		if len(next) == 0 {
			t.Fatal("cycle detected")
		}

		// Run the test cases in a goroutine.
		var wg sync.WaitGroup
		for _, nextCase := range next {
			wg.Add(1)
			config := NewScriptConfig()
			if nextCase.Config != nil {
				config = *nextCase.Config
			}
			if len(config.ScriptExtensions) == 0 {
				// Default script extension if none is provided.
				config.ScriptExtensions = []ScriptExtension{{Extension: ".sh", Command: "bash"}}
			}

			// Get the script extension for the test case.
			ext, err := dagGetScriptExtension(nextCase.Path, config)
			if err != nil {
				t.Fatal(err)
			}

			nextCase := nextCase // Capture the variable for the goroutine.
			go func() {
				defer wg.Done()
				// Run the test case.
				ScriptTestFile(t, ext.Command, nextCase.Path, config)
			}()
		}

		wg.Wait()

		// Mark the case as done.
		for _, nextCase := range next {
			done[nextCase.Name] = true
			// Remove the case from the open list.
			for i, c := range open {
				if c.Name == nextCase.Name {
					open = append(open[:i], open[i+1:]...)
					break
				}
			}
		}
	}
}

func dagGetScriptExtension(path string, config ScriptConfig) (ScriptExtension, error) {
	// Get extension from the path.
	ext := filepath.Ext(path)
	// Search for fitting script definition among config.ScriptExtensions.
	for _, def := range config.ScriptExtensions {
		if def.Extension == ext {
			return def, nil
		}
	}
	return ScriptExtension{}, fmt.Errorf("no script definition found for path %s with extension %s", path, ext)
}

func validate(cases []DagTestCase) error {
	// Ensure that all cases have unique names.
	names := make(map[string]bool)
	for _, c := range cases {
		if names[c.Name] {
			return fmt.Errorf("duplicate test case name: %s", c.Name)
		}
		names[c.Name] = true
	}

	// Ensure that all dependencies are valid.
	// and that there is only one case with no dependencies.
	hasNoNeeds := false
	for _, c := range cases {
		if len(c.Needs) == 0 {
			if hasNoNeeds {
				return fmt.Errorf("multiple cases with no dependencies")
			}
			hasNoNeeds = true
		}
		for _, need := range c.Needs {
			if !names[need] {
				return fmt.Errorf("unknown dependency: %s", need)
			}
		}
	}

	return nil
}
