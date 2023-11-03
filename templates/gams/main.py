"""
Template for working with GAMS.
"""

import argparse
import json
import sys
from typing import Any, Dict

import gamspy as gp


# Status of the solver after optimizing.
STATUS = {
    gp.ModelStatus.Feasible: "suboptimal",
    gp.ModelStatus.InfeasibleGlobal: "infeasible",
    gp.ModelStatus.OptimalGlobal: "optimal",
    gp.ModelStatus.Unbounded: "unbounded",
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

    # Creates the container.
    container = gp.Container()
    provider = "CPLEX"

    # Initializes the linear sums.
    weights = 0.0
    values = 0.0

    # Creates the decision variables and adds them to the linear sums.
    items = []
    for item in input_data["items"]:
        item_variable = gp.Variable(container, name=item["id"], type="binary")
        items.append({"item": item, "variable": item_variable})
        weights += item_variable * item["weight"]
        values += item_variable * item["value"]

    # This constraint ensures the weight capacity of the knapsack will not be
    # exceeded.
    gp.Equation(
        container,
        name="weight_capacity",
        definition=weights <= input_data["weight_capacity"],
    )

    # Creates the model and solves the problem.
    model = gp.Model(
        container,
        name="knapsack",
        equations=container.getEquations(),
        problem="MIP",
        # Sets the objective function: maximize the value of the chosen items.
        sense=gp.Sense.MAX,
        objective=values,
    )
    model.solve(
        solver=provider,
        solver_options={"timelimit": duration},
    )

    # Determines which items were chosen.
    chosen_items = []
    for item in items:
        if item["variable"].records.level[0] > 0.9:
            chosen_items.append(item["item"])

    # Creates the statistics.
    statistics = {
        "result": {
            "custom": {
                "constraints": len(container.getEquations()),
                "provider": provider,
                "status": STATUS.get(model.status, "unknown"),
                "variables": model.num_variables,
            },
            "duration": model.total_solve_time,
            "value": model.objective_value,
        },
        "run": {
            "duration": model.total_solve_time,
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
