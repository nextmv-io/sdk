package golden

import "testing"

type DagTestCase struct {
	name   string
	needs  []string
	config BashConfig
	path   string
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
	open := cases
	done := make(map[string]bool)

	for len(open) > 0 {
		// Pick the first case from the open list that has all its needs met.
		var next DagTestCase
		for _, c := range open {
			ready := true
			for _, need := range c.needs {
				if !done[need] {
					ready = false
					break
				}
			}
			if ready {
				next = c
				break
			}
		}

		// If we didn't find a case to run, we have a cycle.
		if next.name == "" {
			t.Fatal("cycle detected")
		}

		// Run the case and mark it as done.
		t.Run(next.name, func(t *testing.T) {
			// Run the test case.
			BashTest(t, next.path, next.config)
		})
		done[next.name] = true

		// Remove the case from the open list.
		for i, c := range open {
			if c.name == next.name {
				open = append(open[:i], open[i+1:]...)
				break
			}
		}
	}
}
