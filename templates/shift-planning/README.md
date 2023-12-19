# Shift creation with OR-Tools

This is an example of how to use OR-Tools to solve a shift creation problem. The
goal is to select/create a number of shifts according to a given demand that
will later be filled by employees.

## Usage

```bash
python main.py -input data/input.json -output output.json -duration 30
python plot.py -input output.json
```

## Model formulation

TODO: add model formulation

### Variables

- $x_{s,t}$: number of times

### Parameters

- $p_{i,j}$: preference of employee $i$ for shift $j$
- $d_{j}$: number of employees needed for shift $j$
- $m_{i}$: max number of shifts employee $i$ can work

- Each employee has a max

This sentence uses `$` delimiters to show math inline:  $\sqrt{3x-1}+(1+x)^2$
