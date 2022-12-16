// Package sudoku holds the implementation of the sudoku template.
package main

import (
	"fmt"
	"hash/maphash"
	"time"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	run.Run(solver)
}

func solver(input [9][9]int, opts store.Options) (store.Solver, error) {
	root := store.New()

	// Create 81 domains, each domain has 9 possible values for each cell.
	x := store.Repeat(root, 9*9, model.NewDomain(model.NewRange(1, 9)))

	// Create 27 constraints, one for each of the 9 rows, one for each of
	// the 9 columns and one for each of the 9 regions. Each constraint
	// propagates values from the domains which can no longer be part of a
	// solution assuming the domains will only have values removed during the
	// search.
	constraints := make([]constraint, 0)

	for i := 0; i < 9; i++ {
		r := make([]int, 9)
		c := make([]int, 9)

		for j := 0; j < 9; j++ {
			r[j] = index(i, j)
			c[j] = index(j, i)
		}

		constraint, err := uniqueConstraint(root, x, r)
		if err != nil {
			return nil, err
		}

		constraints = append(constraints, constraint)

		constraint, err = uniqueConstraint(root, x, c)

		if err != nil {
			return nil, err
		}

		constraints = append(constraints, constraint)
	}

	for rr := 0; rr < 3; rr++ {
		for rc := 0; rc < 3; rc++ {
			indices := make([]int, 9)

			i := 0

			for r := rr * 3; r < (rr+1)*3; r++ {
				for c := rc * 3; c < (rc+1)*3; c++ {
					indices[i] = index(r, c)
					i++
				}
			}

			constraint, err := uniqueConstraint(root, x, indices)
			if err != nil {
				return nil, err
			}

			constraints = append(constraints, constraint)
		}
	}

	// Assign the partially completed grid provided as input.
	for i, row := range input {
		for j, cell := range row {
			if cell >= 1 && cell <= 9 {
				root = root.Apply(x.Assign(index(i, j), cell))

				for _, constraint := range constraints {
					root = root.Propagate(constraint.propagate)
				}
			}
		}
	}

	// Check if it is not proven infeasible.
	for _, constraint := range constraints {
		if constraint.provenInfeasible(root) {
			return nil,
				fmt.Errorf("input puzzle is proven infeasible")
		}
	}

	// Now the most important part: given a store, how can we create new
	// stores that bring us closer to a solution.
	root = root.Generate(func(s store.Store) store.Generator {
		// We select the cell with the smallest number of possible values.
		i, ok := x.Smallest(s)

		values := x.Domain(s, i).Slice()

		return store.Lazy(
			func() bool {
				// We can keep assigning values to cells if it is not proven
				// infeasible and there are still values to be tried.
				for _, constraint := range constraints {
					if constraint.provenInfeasible(s) {
						return false
					}
				}
				return ok && len(values) > 0
			},
			func() store.Store {
				// We will try to assign the values in the selected smallest
				// domain one by one, starting with the first one at position 0.
				next := values[0]
				// Remove the next value from the options for the next.
				// invocation
				values = values[1:]

				s2 := s.Apply([]store.Change{
					x.AtLeast(i, next),
					x.AtMost(i, next),
				}...)

				for _, constraint := range constraints {
					s2 = s2.Propagate(constraint.propagate)
				}

				for _, constraint := range constraints {
					if constraint.provenInfeasible(s2) {
						return nil
					}
				}
				return s2
			},
		)
	}).Validate(func(s store.Store) bool {
		// We have a valid solution if all cells have a single value.
		return x.Singleton(s)
	}).Format(format(x))

	// If the duration limit is unset, we set it to 10s. You can configure
	// longer solver run times here. For local runs there is no time limitation.
	// If you want to make cloud runs for longer than 5 minutes, please contact:
	// support@nextmv.io
	if opts.Limits.Duration == 0 {
		opts.Limits.Duration = 10 * time.Second
	}

	return root.Satisfier(opts), nil
}

// Now we define some helper structs and functions.
// We define a constraint interface which allows us to propagate and
// query if the constraint is proven to be infeasible.
type constraint interface {
	propagate(store.Store) []store.Change
	provenInfeasible(store.Store) bool
}

type unique struct {
	domains store.Domains
	indices []int
}

// uniqueConstraint creates a unique constraint for all the domains at indices
// in domains.
func uniqueConstraint(
	s store.Store,
	domains store.Domains,
	indices []int,
) (constraint, error) {
	for _, i := range indices {
		if i < 0 || i >= domains.Len(s) {
			return nil,
				fmt.Errorf(
					"all different, no domain for index %v",
					i,
				)
		}
	}

	return &unique{
		domains: domains,
		indices: indices,
	}, nil
}

var hasher = maphash.Hash{}

// Hash function used in propagate.
func hash(values []int) (uint64, error) {
	hasher.Reset()
	for _, value := range values {
		if _, err := hasher.WriteString(fmt.Sprint(value)); err != nil {
			return uint64(0), err
		}
	}
	h := hasher.Sum64()

	return h, nil
}

// As soon as we find n cells with the same n values we can
// remove the n values from the other cells that have to be unique.
//
// For example n=1, if we find one cell with one value we can
// remove that value from the other cells.
// For example n=2, if we find two cells with each only two possible values,
// and they are the same values we can remove them from the other cells.
func (a *unique) propagate(s store.Store) []store.Change {
	changes := make([]store.Change, 0)

	type tupleData struct {
		indices map[int]bool
		values  []int
	}
	// For each domain having n values we keep track what values, and what
	// indices.
	// For example a cell at index 6 having values [1,2,4]
	// 		tuples[3][hash([1,2,4]){indices: [6], values: [1,2,4]}
	// whenever there is another cell for example at index 5 with the same
	// values we add it:
	// 		tuples[3][hash([1,2,4]){indices: [6, 5], values: [1,2,4]}
	// as soon as len(indices) is 3 we can remove it from the other cells not
	// in tuples[3][hash([1,2,4]].indices but in a.indices.
	tuples := make(map[int]map[uint64]tupleData)

	for _, i := range a.indices {
		domain := a.domains.Domain(s, i)

		if domain.Empty() {
			return changes
		}

		if !domain.Empty() && domain.Len() < len(a.indices) {
			values := domain.Slice()

			valuesHash, err := hash(values)
			if err != nil {
				panic(err)
			}

			if _, present := tuples[domain.Len()]; !present {
				tuples[domain.Len()] = make(map[uint64]tupleData)
			}

			if _, present := tuples[domain.Len()][valuesHash]; !present {
				tuples[domain.Len()][valuesHash] = tupleData{
					indices: make(map[int]bool),
					values:  values,
				}
			}

			tuple := tuples[domain.Len()][valuesHash]

			if len(tuple.indices) < len(tuple.values) {
				tuple.indices[i] = true
				tuples[domain.Len()][valuesHash] = tuple
			}
		}
	}

	for size, tuple := range tuples {
		for _, h := range tuple {
			if len(h.indices) >= size {
				for _, i := range a.indices {
					if _, present := h.indices[i]; !present {
						for _, v := range h.values {
							if a.domains.Domain(s, i).Contains(v) {
								changes = append(
									changes,
									a.domains.Remove(i, []int{v}),
								)
							}
						}
					}
				}
			}
		}
	}

	return changes
}

// provenInfeasible returns true if the values in the domains which have to
// be unique can no longer be unique under the assumption values can only
// be removed.
func (a *unique) provenInfeasible(s store.Store) bool {
	singles := make(map[int]bool)

	for _, i := range a.indices {
		domain := a.domains.Domain(s, i)

		if domain.Empty() {
			return true
		}

		if v, singleton := domain.Value(); singleton {
			if _, present := singles[v]; present {
				return true
			}
			singles[v] = true
		}
	}

	return false
}

// Index returns the index in the array of domains for the cell [col, row].
func index(row, col int) int {
	return (row * 9) + col
}

// format returns a function to format the solution output.
func format(
	x store.Domains,
) func(s store.Store) any {
	return func(s store.Store) any {
		grid := [9][9]model.Domain{}
		for i := 0; i < 9*9; i++ {
			grid[i/9][i%9] = x.Domain(s, i)
		}
		return grid
	}
}
