# Nextmv OR-Tools Python template

This template demonstrates how to solve a Mixed Integer Programming problem
using the open source software suite [OR-Tools][or-tools].

To solve a Mixed Integer Problem (MIP) is to optimize a linear objective
function of many variables, subject to linear constraints. We demonstrate this
by solving the knapsack problem using the [integer
optimzation][integer-optimization] interface.

Knapsack is a classic combinatorial optimization problem. Given a collection of
items with a value and weight, our objective is to maximize the total value
without exceeding the weight capacity of the knapsack.

The input defines a number of items which have an id to identify the item, a
weight and a value. Additionally there is a weight capacity.

The most important files created are `main.py` and `input.json`.

* `main.py` implements a MIP knapsack solver.
* `input.json` is a sample input file.

Run the command below to check that everything works as expected:

```bash
python3 main.py input.json
```

The optimal knapsack solution should be printed to stdout.

## Next steps

* Open `main.py` and read through the comments to understand the model.
* Further documentation, guides, and API references about custom modelling and
  deployment can also be found on our [blog](https://www.nextmv.io/blog) and on
  our [documentation site](https://docs.nextmv.io).
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!

[or-tools]: https://developers.google.com/optimization
[integer-optimization]: https://developers.google.com/optimization/mip
