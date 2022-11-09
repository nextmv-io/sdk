# Nextmv MIP incentive allocation template

This template demonstrates how to solve a Mixed Integer Programming problem
using the open-source solver [HiGHS](https://github.com/ERGO-Code/HiGHS). To run
this beta template, please contact support@nextmv.io for access to the Nextmv
HiGHS interface.

To solve a Mixed Integer Problem (MIP) is to optimize a linear objective
function of many variables, subject to linear constraints. We demonstrate this
by solving a incentive allocation problem.

Given a collection of incentives per a user with a cost and effect, our
objective is to maximize the total effect without exceeding a given budget.

The input defines a number of users which have an id to identify them, and a set
of in incentives per user with a cost, an effect and a name. There is also an
available budget.

The most important files created are `main.go` and `input.json`.

* `main.go` implements a MIP incentive allocation solver.
* `input.json` is a sample input file that follows the input definition in
`main.go`.

Run the command below to see if everything works as expected:

```bash
nextmv run local main.go -- -hop.runner.input.path input.json \
  -hop.runner.output.path output.json -hop.solver.limits.duration 10s
```

A file `output.json` should have been created with the optimal knapsack
solution.

## Next steps

* Open `main.go` and read through the comments to understand the model.
* API documentation and examples can be found in the [package
  documentation](https://pkg.go.dev/github.com/nextmv-io/sdk/mip).
* Further documentation, guides, and API references about custom modelling and
deployment can also be found on our [blog](https://www.nextmv.io/blog) and on
our [documentation site](https://docs.nextmv.io).
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
