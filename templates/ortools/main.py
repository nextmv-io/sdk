"""
Template for working with Google OR-Tools.
"""

import json
import sys
import time
from typing import Any, Dict

from ortools.linear_solver import pywraplp


def main():
    """Entry point for the template."""
    input_data = read_input()
    solution = solve(input_data)

    # Writes to stdout.
    print(json.dumps(solution, indent=2))


def solve(input_data: Dict[str, Any]) -> Dict[str, Any]:
    """Solves the given problem and returns the solution."""

    # Creates the solver.
    provider = "SCIP"
    solver = pywraplp.Solver.CreateSolver(provider)

    # Initializes the linear sums.
    weights = 0.0
    values = 0.0

    # Creates the decision variables and adds them to the linear sums.
    items = []
    for item in input_data["items"]:
        item_variable = solver.IntVar(0, 1, item["item_id"])
        items.append({"item": item, "variable": item_variable})
        weights += item_variable * item["weight"]
        values += item_variable * item["value"]

    # This constraint ensures the weight capacity of the knapsack will not be
    # exceeded.
    solver.Add(weights <= input_data["weight_capacity"])

    # Sets the objective function: maximize the value of the chosen items.
    solver.Maximize(values)

    # Solves the problem.
    start_time = time.time()
    status = solver.Solve()
    end_time = time.time()

    # Determines which items were chosen.
    chosen_items = []
    for item in items:
        if item["variable"].solution_value() == 1.0:
            chosen_items.append(item["item"])

    # Creates the statistics.
    statistics = {
        "result": {
            "custom": {
                "constraints": solver.NumConstraints(),
                "provider": provider,
                "status": status,
                "variables": solver.NumVariables(),
            },
            "duration": end_time - start_time,
            "value": solver.Objective().Value(),
        },
        "run": {
            "duration": end_time - start_time,
        },
        "schema": "v1",
    }

    return {
        "solutions": [{"items": chosen_items}],
        "statistics": statistics,
    }


def read_input() -> Dict[str, Any]:
    """Reads the input from stdin."""
    input_file = {}
    if len(sys.argv) == 2:
        with open(sys.argv[1], "r", encoding="utf-8") as file:
            input_file = json.load(file)

    elif len(sys.argv) == 1:
        input_file = json.load(sys.stdin)

    else:
        print("Usage: main.py filename")
        sys.exit(0)

    return input_file


if __name__ == "__main__":
    main()
