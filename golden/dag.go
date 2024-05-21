package golden

import (
	"fmt"
	"sync"
	"testing"
)

// DagTestCase represents a test case in a directed acyclic graph (DAG) test.
type DagTestCase struct {
	Name   string
	Needs  []string
	Config *BashConfig
	Path   string
}

// DagTest runs a set of test cases in topological order.
// Each test case is a BashTest, and the test cases are connected by their
// dependencies. If a test case has dependencies, it will only be run after all
// of its dependencies have been run.
//
// Sample usage:
//
//	cases := []golden.DagTestCase{
//	  {
//	    name:   "app-create",
//	    needs:  []string{},
//	    config: BashConfig{ /**/ },
//	    path:   "app-create",
//	  },
//	  {
//	    name:   "app-push",
//	    needs:  []string{"app-create"},
//	    config: BashConfig{ /**/ },
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
			config := BashConfig{}
			if nextCase.Config != nil {
				config = *nextCase.Config
			}

			nextCase := nextCase
			go func() {
				// Run the test case.
				BashTestFile(t, nextCase.Path, config)
				wg.Done()
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
