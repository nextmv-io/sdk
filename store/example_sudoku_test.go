package store_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/store"
)

func assign(
	s store.Store,
	puzzle store.Domains,
	row, col, value int,
) store.Store {
	ind := index(row, col)
	changes := []store.Change{puzzle.Assign(ind, value)}

	for j := 0; j < 9; j++ {
		if j != col {
			changes = append(changes, puzzle.Remove(index(row, j), value))
		}
		if j != row {
			changes = append(changes, puzzle.Remove(index(j, col), value))
		}
	}

	i, j := 3*(row/3), 3*(col/3)
	for m := i; m < i+3; m++ {
		for n := j; n < j+3; n++ {
			if k := index(m, n); k != ind {
				changes = append(changes, puzzle.Remove(k, value))
			}
		}
	}

	return s.Apply(changes...)
}

func index(row, col int) int {
	return (row * 9) + col
}

// Solve a classic sudoku puzzle.
func Example_sudoku() {
	// Initial puzzle.
	initial := [9][9]int{
		{0, 0, 0, 0, 0, 6, 0, 3, 0},
		{1, 5, 0, 0, 0, 0, 0, 0, 4},
		{0, 0, 6, 0, 0, 0, 0, 0, 0},
		{0, 0, 3, 0, 0, 8, 4, 0, 0},
		{2, 0, 0, 0, 0, 0, 0, 6, 0},
		{0, 0, 0, 0, 4, 9, 0, 2, 5},
		{0, 0, 0, 7, 0, 5, 0, 0, 0},
		{0, 7, 0, 0, 9, 0, 1, 0, 3},
		{0, 6, 0, 0, 0, 0, 8, 0, 0},
	}

	// Create the sudoku model.
	sudoku := store.New()

	// Create a variable to track the puzzle.
	puzzle := store.Repeat(sudoku, 9*9, model.NewDomain(model.NewRange(1, 9)))

	// Assign the initial numbers.
	for i, row := range initial {
		for j, cell := range row {
			if cell >= 1 && cell <= 9 {
				sudoku = assign(sudoku, puzzle, i, j, cell)
			}
		}
	}

	// Define how the solution is generated.
	sudoku = sudoku.Generate(
		store.Scope(func(s store.Store) store.Generator {
			i, ok := puzzle.Smallest(s)
			values := puzzle.Domain(s, i).Slice()

			for i := 0; i < 9*9; i++ {
				if puzzle.Domain(s, i).Empty() {
					ok = false
					break
				}
			}

			return store.
				If(func(store.Store) bool { return ok && len(values) > 0 }).
				Then(func(store.Store) store.Store {
					// Get the first element of the stores' queue and pop it.
					next := values[0]
					values = values[1:]
					return assign(s, puzzle, i/9, i%9, next)
				}).
				With(puzzle.Singleton)
		}),
	).Format(func(s store.Store) any {
		board := map[string]string{}
		cells := puzzle.Domains(s).Slices()
		for i := 0; i < 9; i++ {
			row := cells[i*9 : i*9+9]
			board[strconv.Itoa(i)] = strings.ReplaceAll(
				strings.ReplaceAll(
					strings.Join(strings.Fields(fmt.Sprint(row)), " "),
					"[", "",
				),
				"]", "",
			)
		}

		return board
	})

	// The solver type is a satisfier because only operationally valid grids
	// are needed, there is no value associated.
	opt := store.DefaultOptions()
	opt.Limits.Solutions = 1
	solver := sudoku.Satisfier(opt)

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Output:
	// {
	//   "0": "8 2 9 4 5 6 7 3 1",
	//   "1": "1 5 7 3 8 2 6 9 4",
	//   "2": "4 3 6 9 1 7 5 8 2",
	//   "3": "6 9 3 5 2 8 4 1 7",
	//   "4": "2 4 5 1 7 3 9 6 8",
	//   "5": "7 8 1 6 4 9 3 2 5",
	//   "6": "3 1 8 7 6 5 2 4 9",
	//   "7": "5 7 2 8 9 4 1 4 3",
	//   "8": "9 6 4 2 3 1 8 5 7"
	// }
}
