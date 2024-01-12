import argparse
import datetime
import json
import sys
from typing import Any

from pyomo.environ import ConcreteModel, Var, Objective, Constraint, SolverFactory, Integers

# Duration parameter for the solver.
SUPPORTED_PROVIDER_DURATIONS = {
    "cbc": "sec",
    "glpk": "tmlim",
}

# Status of the solver after optimizing.
STATUS = {
    "FEASIBLE": "feasible",
    "INFEASIBLE": "infeasible",
    "OPTIMAL": "optimal",
    "UNBOUNDED": "unbounded",
}
ANY_SOLUTION = ["FEASIBLE", "OPTIMAL"]

def main() -> None:
    """Entry point for the template."""

    parser = argparse.ArgumentParser(
        description="Solve shift-creation with Pyomo MIP."
    )
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
    parser.add_argument(
        "-provider",
        default="cbc",
        help="Solver provider. Default is cbc.",
    )
    args = parser.parse_args()

    # Read input data, solve the problem and write the solution.
    input_data = read_input(args.input)
    log("Solving shift-creation:")
    log(f"  - shifts-templates: {len(input_data.get('shifts', []))}")
    log(f"  - demands: {len(input_data.get('demands', []))}")
    log(f"  - max duration: {args.duration} seconds")
    solution = solve(input_data, args.duration, args.provider)
    write_output(args.output, solution)

def solve(input_data: dict[str, Any], duration: int, provider: str) -> dict[str, Any]:
    """Solves the given problem and returns the solution."""

    # Creates the solver.
    solver = SolverFactory(provider)  # Use an appropriate solver name
    solver.options[SUPPORTED_PROVIDER_DURATIONS[provider]] = duration * 1000  # Pyomo time limit is in milliseconds

    # Prepare data
    shifts, demands = convert_data(input_data)

    # Generate concrete shifts from shift templates.
    concrete_shifts = get_concrete_shifts(shifts)

    # Determine all unique time periods in which demands occur and the shifts covering them.
    periods = get_coverage(concrete_shifts, demands)

    # Create integer variables indicating how many times a shift is planned.
    model = ConcreteModel()
    model.x_assign = Var([s["id"] for s in concrete_shifts], within=Integers)

    # Objective function: minimize the cost of the planned shifts
    model.set_objective(
        sum(model.x_assign[s["id"]] * s["cost"] for s in concrete_shifts)
    )

    # >> Constraints

    # We need to make sure that all demands are covered
    for p in periods:
        model.add_constraint(
            sum(model.x_assign[s["id"]] for s in p.covering_shifts) >= sum(d["count"] for d in p.demands),
            f"DemandCover_{p.start_time}_{p.end_time}_{p.qualification}",
        )

    # Solves the problem.
    results = solver.solve(model)

    # Convert to solution format.
    schedule = {
        "planned_shifts": [
            {
                "id": s["id"],
                "shift_id": s["shift_id"],
                "time_id": s["time_id"],
                "start_time": s["start_time"],
                "end_time": s["end_time"],
                "qualification": s["qualification"],
                "count": int(round(model.x_assign[s["id"]].value)),
            }
            for s in concrete_shifts
            if model.x_assign[s["id"]].value > 0.5
        ]
        if results.solver.termination_condition in ANY_SOLUTION
        else [],
    }

    # Creates the statistics.
    statistics = {
        "result": {
            "custom": {
                "constraints": len(model.constraints),
                "provider": "Pyomo",
                "status": STATUS.get(results.solver.termination_condition, "unknown"),
                "variables": len(model.x_assign),
            },
            "duration": results.solver.wallclock_time,
            "value": model.objective() if results.solver.termination_condition in ANY_SOLUTION else None,
        },
        "run": {
            "duration": results.solver.wallclock_time,
        },
        "schema": "v1",
    }

    log(f"  - status: {statistics['result']['custom']['status']}")
    log(f"  - value: {statistics['result']['value']}")

    return {
        "solutions": [schedule],
        "statistics": statistics,
    }


def convert_input(input_data: dict[str, Any]) -> tuple[list, list, dict]:
    """Converts the input data to the format expected by the model."""
    workers = input_data["workers"]
    shifts = input_data["shifts"]

    # In-place convert timestamps to datetime objects
    for s in shifts:
        s["start_time"] = datetime.datetime.fromisoformat(s["start_time"])
        s["end_time"] = datetime.datetime.fromisoformat(s["end_time"])
    for e in workers:
        for a in e["availability"]:
            a["start_time"] = datetime.datetime.fromisoformat(a["start_time"])
            a["end_time"] = datetime.datetime.fromisoformat(a["end_time"])

    # Add default values for rules
    for r in input_data["rules"]:
        r["min_shifts"] = r.get("min_shifts", 0)
        r["max_shifts"] = r.get("max_shifts", 1000)

    # Merge availabilities of workers that start right where another one ends
    for e in workers:
        e["availability"] = sorted(e["availability"], key=lambda x: x["start_time"])
        i = 0
        while i < len(e["availability"]) - 1:
            if (
                e["availability"][i]["end_time"]
                == e["availability"][i + 1]["start_time"]
            ):
                e["availability"][i]["end_time"] = e["availability"][i + 1]["end_time"]
                del e["availability"][i + 1]
            else:
                i += 1

    # Convert rules to dict
    rules_per_worker = {}
    for e in workers:
        rule = [r for r in input_data.get("rules", {}) if r["id"] == e["rules"]]
        if len(rule) != 1:
            raise ValueError(f"Invalid rule for worker {e['id']}")
        rules_per_worker[e["id"]] = rule[0]

    return workers, shifts, rules_per_worker


def custom_serial(obj):
    """JSON serializer for objects not serializable by default serializer."""
    if isinstance(obj, (datetime.datetime | datetime.date)):
        return obj.isoformat()
    raise TypeError("Type %s not serializable" % type(obj))


def log(message: str) -> None:
    """Logs a message. We need to use stderr since stdout is used for the solution."""
    print(message, file=sys.stderr)


def read_input(input_path) -> dict[str, Any]:
    """Reads the input from stdin or a given input file."""
    input_file = {}
    if input_path:
        with open(input_path, encoding="utf-8") as file:
            input_file = json.load(file)
    else:
        input_file = json.load(sys.stdin)
    return input_file


def write_output(output_path, output) -> None:
    """Writes the output to stdout or a given output file."""
    content = json.dumps(output, indent=2, default=custom_serial)
    if output_path:
        with open(output_path, "w", encoding="utf-8") as file:
            file.write(content + "\n")
    else:
        print(content)



if __name__ == "__main__":
    main()
