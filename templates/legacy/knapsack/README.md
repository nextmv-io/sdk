# Nextmv knapsack template

Knapsack is a classic combinatorial optimization problem. Given a collection of
items with a value and weight, our objective is to maximize the total value
without exceeding the weight capacity of the knapsack.

This template uses our *custom modeling* framework `store` to model and solve
such a packing problem.

The most important files created are `main.go` and `input.json`.

* `main.go` implements a Knapsack solver.
* `input.json` is a sample input file that follows the input definition in
`main.go`.

Before you start customizing, run the command below to see if everything works
as expected:

```bash
nextmv sdk run . -- -runner.input.path input.json \
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with the optimal Knapsack solution

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
