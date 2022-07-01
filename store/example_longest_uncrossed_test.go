package store_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nextmv-io/sdk/store"
)

// intersects returns true if the segment p0->p1 intersects with the segment
// p2->p3. It is based on this link:
// https://stackoverflow.com/a/14795484/15559724
func intersects(p0, p1, p2, p3 position) bool {
	s10x := p1.row - p0.row
	s10y := p1.col - p0.col
	s32x := p3.row - p2.row
	s32y := p3.col - p2.col

	denom := s10x*s32y - s32x*s10y
	if denom == 0 {
		return false
	}
	denomPositive := denom > 0

	s02x := p0.row - p2.row
	s02y := p0.col - p2.col
	sNumer := s10x*s02y - s10y*s02x
	if (sNumer < 0) == denomPositive {
		return false
	}

	tNumer := s32x*s02y - s32y*s02x
	if (tNumer < 0) == denomPositive {
		return false
	}

	if ((sNumer > denom) == denomPositive) || ((tNumer > denom) == denomPositive) {
		return false
	}

	return true
}

// intersection asserts that the option does not intersect the tour.
func intersection(opt position, tour []position, visited bool) bool {
	intersection := false
	if len(tour) >= 3 && !visited {
		for i := 0; i < len(tour)-1; i++ {
			if intersects(tour[i], tour[i+1], tour[len(tour)-1], opt) {
				intersection = true
				break
			}
		}
	}
	return intersection
}

// unintersected returns the positions that a knight can move to from the given
// row and column, asserting that they are inside the board, have not been
// visited yet and do not intersect the knight's path.
func unintersected(n, row, col int, tour []position) []position {
	// Create all the possible movement options.
	options := make([]position, len(offsets))
	for i := range options {
		options[i] = position{
			row: row + offsets[i].row,
			col: col + offsets[i].col,
		}
	}

	// Create positions that are feasible candidates.
	var positions []position
	for _, opt := range options {
		// Assert the option is inside the board, has not been visited and does
		// not intersect the path.
		if opt.row >= 0 && opt.row < n && opt.col >= 0 && opt.col < n {
			visited := visited(opt, tour)
			if !visited && !intersection(opt, tour, visited) {
				positions = append(positions, opt)
			}
		}
	}

	return positions
}

// The longest uncrossed (or nonintersecting) knight's path is a mathematical
// problem involving a knight on a square n×n chess board. The problem is to
// find the longest path the knight can take on the given board, such that the
// path does not intersect itself. Definitions reused from the knight's tour
// example: position, offsets, format, visited.
func Example_longestUncrossedKnightsPath() {
	// Board size and initial position.
	n := 5
	p := position{row: 0, col: 0}

	// Create the knight's tour model.
	knight := store.New()

	// Track the sequence of moves.
	tour := store.NewSlice(knight, p)

	// Define the output format.
	knight = knight.Format(format(tour, n))

	// Define the value to maximize: the number of jumps made.
	knight = knight.Value(func(s store.Store) int { return tour.Len(s) - 1 })

	// Define the generation of the tour.
	knight = knight.Generate(
		store.Scope(func(s store.Store) store.Generator {
			// Gets the last move made and all the candidate positions from
			// there.
			lastMove := tour.Get(s, tour.Len(s)-1)
			candidates := unintersected(
				n,
				lastMove.row,
				lastMove.col,
				tour.Slice(s),
			)

			// Create new stores by adding each candidate to the tour.
			stores := make([]store.Store, len(candidates))
			for i, candidate := range candidates {
				stores[i] = s.Apply(tour.Append(candidate))
			}

			return store.
				// Generate new stores as long as there are elements left.
				If(func(s store.Store) bool { return len(stores) > 0 }).
				Then(func(s store.Store) store.Store {
					// Get the first element of the stores' queue and pop it.
					generated := stores[0]
					stores = stores[1:]
					return generated
				}) // The store is always operationally valid.
		}),
	)

	// The solver type is a maximizer because the store is searching for the
	// highest number of moves.
	solver := knight.Maximizer(store.DefaultOptions())

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Output:
	// {
	//   "0": "00 -- 02 -- -- ",
	//   "1": "-- -- -- -- 03 ",
	//   "2": "06 01 -- -- -- ",
	//   "3": "-- -- 07 04 -- ",
	//   "4": "-- 05 -- -- -- "
	// }
}
