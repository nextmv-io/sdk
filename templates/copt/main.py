"""
Template for working with COPT.
"""

import argparse
import json
import sys
from typing import Any, Dict

import coptpy as cp
from coptpy import COPT


# Status of the solver after optimizing.
STATUS = {
    COPT.UNFINISHED: "suboptimal",
    COPT.INFEASIBLE: "infeasible",
    COPT.OPTIMAL: "optimal",
    COPT.UNBOUNDED: "unbounded",
}


def main() -> None:
    """Entry point for the template."""

    parser = argparse.ArgumentParser(description="Solve problems with COPT.")
    parser.add_argument(
        "-input",
        default="",
        help="Path to input file. Default is stdin.",
    )
    parser.add_argument(
        "-output",
        default="",
        help="Path to output file. Default is stdout.",
    )
    parser.add_argument(
        "-duration",
        default=30,
        help="Max runtime duration (in seconds). Default is 30.",
        type=int,
    )
    args = parser.parse_args()

    # Read input data, solve the problem and write the solution.
    input_data = read_input(args.input)
    log("Solving knapsack problem:")
    log(f"  - items: {len(input_data.get('items', []))}")
    log(f"  - capacity: {input_data.get('weight_capacity', 0)}")
    log(f"  - max duration: {args.duration} seconds")
    solution = solve(input_data, args.duration)
    write_output(args.output, solution)


def solve(input_data: Dict[str, Any], duration: int) -> Dict[str, Any]:
    """Solves the given problem and returns the solution."""

    # Creates the model.
    env = cp.Envr()
    model = env.createModel()
    model.setParam(COPT.Param.Logging, 0)  # Turns off verbosity.
    model.setParam(COPT.Param.TimeLimit, duration)

    # Initializes the linear sums.
    weights = 0.0
    values = 0.0

    # Creates the decision variables and adds them to the linear sums.
    items = []
    for item in input_data["items"]:
        item_variable = model.addVar(lb=0, ub=1, vtype=COPT.BINARY, name=item["id"])
        items.append({"item": item, "variable": item_variable})
        weights += item_variable * item["weight"]
        values += item_variable * item["value"]

    # This constraint ensures the weight capacity of the knapsack will not be
    # exceeded.
    model.addConstr(weights <= input_data["weight_capacity"])

    # Sets the objective function: maximize the value of the chosen items.
    model.setObjective(values, sense=COPT.MAXIMIZE)

    # Solves the problem.
    model.solve()

    # Determines which items were chosen.
    chosen_items = []
    for item in items:
        if model.getVarByName(item["item"]["id"]).x > 0.9:
            chosen_items.append(item["item"])

    # Creates the statistics.
    statistics = {
        "result": {
            "custom": {
                "constraints": model.rows,
                "provider": "COPT",
                "status": STATUS.get(model.status, "unknown"),
                "variables": model.cols,
            },
            "duration": model.solvingtime,
            "value": model.objval,
        },
        "run": {
            "duration": model.solvingtime,
        },
        "schema": "v1",
    }

    return {
        "solutions": [{"items": chosen_items}],
        "statistics": statistics,
    }


def log(message: str) -> None:
    """Logs a message. We need to use stderr since stdout is used for the solution."""

    print(message, file=sys.stderr)


def read_input(input_path) -> Dict[str, Any]:
    """Reads the input from stdin or a given input file."""

    input_file = {}
    if input_path:
        with open(input_path, "r", encoding="utf-8") as file:
            input_file = json.load(file)
    else:
        input_file = json.load(sys.stdin)

    return input_file


def write_output(output_path, output) -> None:
    """Writes the output to stdout or a given output file."""

    content = json.dumps(output, indent=2)
    if output_path:
        with open(output_path, "w", encoding="utf-8") as file:
            file.write(content + "\n")
    else:
        print(content)


if __name__ == "__main__":
    main()
