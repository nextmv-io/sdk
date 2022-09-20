package store_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/nextmv-io/sdk/store"
)

// position tracks the knight on the board.
type position struct {
	row int
	col int
}

// offsets for obtaining the eight different positions a knight can reach from
// any given square (inside the board or not).
var offsets = []struct {
	row int
	col int
}{
	{2, 1}, {1, 2}, {-1, 2}, {-2, 1}, {-2, -1}, {-1, -2}, {1, -2}, {2, -1},
}

// onward implements sort.Interface for []Position based on the moves field.
type onward struct {
	moves      []int
	candidates []position
}

func (o onward) Len() int { return len(o.moves) }
func (o onward) Swap(i, j int) {
	o.moves[i], o.moves[j] = o.moves[j], o.moves[i]
	o.candidates[i], o.candidates[j] = o.candidates[j], o.candidates[i]
}
func (o onward) Less(i, j int) bool { return o.moves[i] < o.moves[j] }

// positions returns the positions that a knight can move to from the given row
// and column, asserting that they are inside the board and have not been
// visited yet.
func positions(n, row, col int, tour []position) []position {
	// Create all the possible movement options.
	options := make([]position, len(offsets))
	for i := range options {
		options[i] = position{
			row: row + offsets[i].row,
			col: col + offsets[i].col,
		}
	}

	// Create positions that are inside the board and have not been visited
	// yet.
	var positions []position
	for _, opt := range options {
		// Assert the option is inside the board.
		if opt.row >= 0 && opt.row < n && opt.col >= 0 && opt.col < n {
			if !visited(opt, tour) {
				positions = append(positions, opt)
			}
		}
	}

	return positions
}

// visited asserts that the option has not been visited in the tour.
func visited(opt position, tour []position) bool {
	visited := false
	for _, move := range tour {
		if reflect.DeepEqual(opt, move) {
			visited = true
			break
		}
	}
	return visited
}

/*
format defines the JSON formatting of the store as a board, e.g.

	{
		"0": "00 -- 02 -- -- ",
		"1": "-- -- -- -- 03 ",
		"2": "06 01 -- -- -- ",
		"3": "-- -- 07 04 -- ",
		"4": "-- 05 -- -- -- "
	}
*/
func format(tour store.Slice[position], n int) func(s store.Store) any {
	return func(s store.Store) any {
		// Empty board.
		board := map[string]string{}
		for i := 0; i < n; i++ {
			board[strconv.Itoa(i)] = strings.Repeat("-- ", n)
		}

		// Loop over the knight's tour to fill the board with positions.
		for i, p := range tour.Slice(s) {
			// Make every number a double digit.
			num := strconv.Itoa(i)
			if i < 10 {
				num = "0" + num
			}

			// Set the visited position on the board.
			cols := strings.Split(board[strconv.Itoa(p.row)], " ")
			cols[p.col] = num
			board[strconv.Itoa(p.row)] = strings.Join(cols, " ")
		}

		return board
	}
}

// A knight's tour is a sequence of moves of a knight on an nxn chessboard such
// that the knight visits every square exactly once. This example implements an
// open knight's tour, given that the last position will not necessarily be one
// move away from the first.
func Example_knightsTour() {
	// Board size and initial position.
	n := 5
	p := position{row: 0, col: 0}

	// Create the knight's tour model.
	knight := store.New()

	// Track the sequence of moves.
	tour := store.NewSlice(knight, p)

	// Define the output format.
	knight = knight.Format(format(tour, n))

	// The store is operationally valid if the tour is complete.
	knight = knight.Validate(func(s store.Store) bool {
		return tour.Len(s) == n*n
	})

	// Define the generation of the tour.
	knight = knight.Generate(func(s store.Store) store.Generator {
		// Gets the last move made and all the candidate positions from
		// there.
		lastMove := tour.Get(s, tour.Len(s)-1)
		candidates := positions(n, lastMove.row, lastMove.col, tour.Slice(s))

		// Obtain the number of onward moves per candidate, excluding
		// visited squares. Sort candidates increasingly by the number
		// of onward moves.
		moves := make([]int, len(candidates))
		for i, candidate := range candidates {
			moves[i] = len(positions(n, candidate.row, candidate.col, tour.Slice(s)))
		}
		onward := onward{moves: moves, candidates: candidates}
		sort.Sort(onward)

		// Starting from the most constrained candidate, create a store
		// queue by adding each candidate to the tour.
		stores := make([]store.Store, len(onward.candidates))
		for i, candidate := range onward.candidates {
			stores[i] = s.Apply(tour.Append(candidate))
		}

		return store.Eager(stores...)
	})

	// The solver type is a satisfier because only operationally valid tours
	// are needed, there is no value associated.
	opt := store.DefaultOptions()
	opt.Limits.Solutions = 1
	opt.Diagram.Expansion.Limit = 1
	solver := knight.Satisfier(opt)

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "0": "00 17 06 11 20 ",
	//   "1": "07 12 19 16 05 ",
	//   "2": "18 01 04 21 10 ",
	//   "3": "13 08 23 02 15 ",
	//   "4": "24 03 14 09 22 "
	// }
}
