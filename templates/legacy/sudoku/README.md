# Nextmv sudoku template

The [Sudoku problem](https://en.wikipedia.org/wiki/Sudoku) asks to assign
values to 81 cells so that each column, each row, and each of the nine 3 × 3
sub-grids that compose the grid contain all of the digits from 1 to 9.

This templates uses our *custom modeling* framework `store` to model and solve
such a Sudoku problem.

The most important files created are `main.go` and `input.json`.

* `main.go` implements a Sudoku solver.
* `input.json` is a sample input file that contains a partially completed
grid. Cells with a value between 1 and 9 are the completed cells, cells with a
value 0 are the ones that need to be decided by the player.

Before you start customizing, run the command below to see if everything works
as expected:

```bash
nextmv sdk run . -- -runner.input.path input.json \
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with a solution to the Sudoku
puzzle.

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
