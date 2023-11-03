"""
Template for working with Mosek.
"""

import argparse
import json
import sys
from typing import Any, Dict

import mosek as mk


# Status of the solver after optimizing.
STATUS = {
    mk.solsta.prim_feas: "suboptimal",
    mk.prosta.prim_infeas: "infeasible",
    mk.solsta.integer_optimal: "optimal",
    mk.prosta.prim_infeas_or_unbounded: "unbounded",
}


def main() -> None:
    """Entry point for the template."""

    parser = argparse.ArgumentParser(description="Solve problems with Mosek.")
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

    # Creates the environment and task.
    env = mk.Env()
    num_vars = len(input_data["items"])
    num_constraints = 1
    task = env.Task(
        numcon=num_constraints,
        numvar=num_vars,
    )
    task.putdouparam(mk.dparam.optimizer_max_time, duration)

    # This constraint ensures the weight capacity of the knapsack will not be
    # exceeded.
    task.appendcons(num_constraints)
    task.putconbound(0, mk.boundkey.up, 0.0, input_data["weight_capacity"])

    # Creates the decision variables and sets the coefficients on the objective
    # function and the constraint matrix.
    task.appendvars(num_vars)
    for i, item in enumerate(input_data["items"]):
        task.putvarbound(i, mk.boundkey.ra, 0.0, 1.0)
        task.putvartype(i, mk.variabletype.type_int)
        task.putcj(i, item["value"])
        task.putaij(0, i, item["weight"])

    # Sets the objective function: maximize the value of the chosen items.
    task.putobjsense(mk.objsense.maximize)

    # Solves the problem.
    task.optimize()

    # Determines which items were chosen.
    chosen_items = []
    for i, item in enumerate(input_data["items"]):
        if task.getxx(mk.soltype.itg)[i] > 0.9:
            chosen_items.append(item)

    # Creates the statistics.
    statistics = {
        "result": {
            "custom": {
                "constraints": num_constraints,
                "provider": "mosek",
                "status": STATUS.get(
                    task.getsolsta(mk.soltype.itg),
                    STATUS.get(task.getprosta(mk.soltype.itg), "unknown"),
                ),
                "variables": num_vars,
            },
            "duration": task.getdouinf(mk.dinfitem.optimizer_time),
            "value": task.getprimalobj(mk.soltype.itg),
        },
        "run": {
            "duration": task.getdouinf(mk.dinfitem.optimizer_time),
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
