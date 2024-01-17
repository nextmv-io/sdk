import argparse
import datetime
import json
import sys
from typing import Any

import pyomo.environ as pyo

# Duration parameter for the solver.
SUPPORTED_PROVIDER_DURATIONS = {
    "cbc": "sec",
    "glpk": "tmlim",
}

# Status of the solver after optimizing.
STATUS = {
    pyo.TerminationCondition.feasible: "suboptimal",
    pyo.TerminationCondition.infeasible: "infeasible",
    pyo.TerminationCondition.optimal: "optimal",
    pyo.TerminationCondition.unbounded: "unbounded",
}

def main() -> None:
    parser = argparse.ArgumentParser(
        description="Solve shift-planning with Pyomo."
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

    input_data = read_input(args.input)
    log("Solving shift-creation:")
    log(f"  - shifts-templates: {len(input_data.get('shifts', []))}")
    log(f"  - demands: {len(input_data.get('demands', []))}")
    log(f"  - max duration: {args.duration} seconds")
    solution = solve(input_data, args.duration, args.provider)
    write_output(args.output, solution)

def solve(input_data: dict[str, Any], duration: int, provider: str) -> dict[str, Any]:
    """Solves the given problem and returns the solution."""

    # Make sure the provider is supported.
    if provider not in SUPPORTED_PROVIDER_DURATIONS:
        raise ValueError(
            f"Unsupported provider: {provider}. The supported providers are: "
            f"{', '.join(SUPPORTED_PROVIDER_DURATIONS.keys())}"
        )

    # Creates the solver.
    solver = pyo.SolverFactory(provider)  # Use an appropriate solver name
    solver.options[SUPPORTED_PROVIDER_DURATIONS[provider]] = duration * 1000  # Pyomo time limit is in milliseconds


    shifts, demands = convert_data(input_data)

    concrete_shifts = get_concrete_shifts(shifts)
    periods = get_coverage(concrete_shifts, demands)

    model = pyo.ConcreteModel()

    # Create a set of shifts.
    model.shifts = pyo.Set(initialize=[s["id"] for s in concrete_shifts])
    # Create a dict for lower and upper bounds of shifts.
    lb = {s["id"]: s["min_workers"] for s in concrete_shifts}
    ub = {s["id"]: s["max_workers"] for s in concrete_shifts}

    def shift_bounds_rule(model, s):
        return (lb[s], ub[s])

    model.x_assign = pyo.Var(model.shifts, within=pyo.NonNegativeIntegers, bounds=shift_bounds_rule)

    model.objective = pyo.Objective(
        expr=sum(model.x_assign[s["id"]] * s["cost"] for s in concrete_shifts),
        sense=pyo.minimize
    )

    # Constraints
    for p in periods:
        model.add_component(
            f"DemandCover_{p.start_time}_{p.end_time}_{p.qualification}",
            pyo.Constraint(
                expr=sum(model.x_assign[s["id"]] for s in p.covering_shifts) >= sum(d["count"] for d in p.demands)
            )
        )


    # Solve the model.
    results = solver.solve(model)

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
    }

    statistics = {
        "result": {
            "custom": {
                "constraints": model.nconstraints(),
                "provider": provider,
                "status": STATUS.get(results.solver.termination_condition, "unknown"),
                "variables": model.nvariables(),
            },
            "duration": results.solver.time,
            "value": pyo.value(model.objective),
        },
        "run": {
            "duration": results.solver.time,
        },
        "schema": "v1",
    }

    log(f"  - status: {statistics['result']['custom']['status']}")
    log(f"  - value: {statistics['result']['value']}")

    return {
        "solutions": [schedule],
        "statistics": statistics,
    }

class UniqueTimeQualificationPeriod:
    def __init__(
        self,
        start_time: datetime.datetime,
        end_time: datetime.datetime,
        qualification: str,
        covering_shifts: list[str],
        demands: list[str],
    ):
        self.start_time = start_time
        self.end_time = end_time
        self.qualification = qualification
        self.covering_shifts = covering_shifts
        self.demands = demands

def get_concrete_shifts(shifts: list[dict[str, Any]]) -> list[dict[str, Any]]:
    concrete_shifts = [
        {
            "id": f"{shift['id']}_{time['id']}",
            "shift_id": shift["id"],
            "time_id": time["id"],
            "start_time": time["start_time"],
            "end_time": time["end_time"],
            "min_workers": time["min_workers"]
            if "min_workers" in time
            else shift["min_workers"]
            if "min_workers" in shift
            else 0,
            "max_workers": time["max_workers"]
            if "max_workers" in time
            else shift["max_workers"]
            if "max_workers" in shift
            else -1,
            "cost": time["cost"] if "cost" in time else shift["cost"],
            "qualification": shift["qualification"] if "qualification" in shift else "",
        }
        for shift in shifts
        for time in shift["times"]
    ]
    return concrete_shifts

def get_coverage(
    concrete_shifts: list[dict[str, Any]], demands: list[dict[str, Any]]
) -> list[UniqueTimeQualificationPeriod]:
    times = set()
    for d in demands:
        times.add(d["start_time"])
        times.add(d["end_time"])
    for s in concrete_shifts:
        times.add(s["start_time"])
        times.add(s["end_time"])
    times = sorted(times)

    periods = []
    for i in range(len(times) - 1):
        start, end = times[i], times[i + 1]
        associated_demands = [
            d for d in demands if d["start_time"] <= start and d["end_time"] >= end
        ]
        qualifications = set()
        for d in associated_demands:
            qualifications.add(d["qualification"] if "qualification" in d else "")
        for q in qualifications:
            periods.append(
                UniqueTimeQualificationPeriod(
                    start,
                    end,
                    q,
                    [
                        s
                        for s in concrete_shifts
                        if s["start_time"] <= start
                        and s["end_time"] >= end
                        and (q == "" or q == s["qualification"])
                    ],
                    [
                        d
                        for d in associated_demands
                        if q == "" or q == d["qualification"]
                    ],
                )
            )

    return periods

def convert_data(
    input_data: dict[str, Any]
) -> tuple[
    list[dict[str, Any]],
    list[dict[str, Any]],
]:
    shifts = input_data["shifts"]
    demands = input_data["demands"]
    for s in shifts:
        for t in s["times"]:
            t["start_time"] = datetime.datetime.fromisoformat(t["start_time"])
            t["end_time"] = datetime.datetime.fromisoformat(t["end_time"])
    for d in demands:
        d["start_time"] = datetime.datetime.fromisoformat(d["start_time"])
        d["end_time"] = datetime.datetime.fromisoformat(d["end_time"])
        d["qualification"] = d["qualification"] if "qualification" in d else ""
    return shifts, demands

def log(message: str) -> None:
    print(message, file=sys.stderr)

def read_input(input_path) -> dict[str, Any]:
    input_file = {}
    if input_path:
        with open(input_path, encoding="utf-8") as file:
            input_file = json.load(file)
    else:
        input_file = json.load(sys.stdin)

    return input_file

def write_output(output_path, output) -> None:
    content = json.dumps(output, indent=2, default=custom_serial)
    if output_path:
        with open(output_path, "w", encoding="utf-8") as file:
            file.write(content + "\n")
    else:
        print(content)

def custom_serial(obj):
    if isinstance(obj, (datetime.datetime | datetime.date)):
        return obj.isoformat()
    raise TypeError("Type %s not serializable" % type(obj))

if __name__ == "__main__":
    main()
