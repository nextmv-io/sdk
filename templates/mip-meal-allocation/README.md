# Nextmv mip-meal-allocation template

This template demonstrates how to solve a Mixed Integer Programming problem
using the open-source solver [HiGHS](https://github.com/ERGO-Code/HiGHS).

To solve a Mixed Integer Problem (MIP) is to optimize a linear objective
function of many variables, subject to linear constraints. We demonstrate this
by solving a made up problem we named MIP meal allocation.

MIP meal allocation is a demo program in which we maximize the number of
binkies our bunnies will execute by selecting their meal.

A binky is when a bunny jumps straight up and quickly twists its hind end,
head, or both. A bunny may binky because it is feeling happy or safe in its
environment.

The input defines a number of meals we can use to maximize binkies. Each
meal consists out of one or more items and in total we can only use the
number of items we have in stock.

The most important files created are `main.go` and `input.json`.

* `main.go` implements a binkies MIP solver.
* `input.json` is a sample input file that follows the input definition in
`main.go`.

Before you start customizing, run the command below to see if everything works
as expected:

```bash
nextmv run local . -- -runner.input.path input.json \
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with the optimal binkies solution

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
