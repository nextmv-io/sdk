import argparse
import datetime
import json
import sys
from typing import Any

from pyomo.environ import ConcreteModel, Constraint, NonNegativeReals, Objective, SolverFactory, Var, minimize

# Duration parameter for the solver.
SUPPORTED_PROVIDER_DURATIONS = {
    "cbc": "sec",
    "glpk": "tmlim",
}

def main() -> None:
    """Entry point for the app."""

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

    # Read input data, solve the problem and write the solution.
    input_data = read_input(args.input)
    log("Solving shift-plannig:")
    log(f"  - shifts: {len(input_data.get('shifts', []))}")
    log(f"  - workers: {len(input_data.get('workers', []))}")
    log(f"  - rules: {len(input_data.get('rules', []))}")
    log(f"  - max duration: {args.duration} seconds")
    solution = solve(input_data, args.duration, args.provider)
    write_output(args.output, solution)


def solve(input_data: dict[str, Any], duration: int, provider: str) -> dict[str, Any]:
    shifts, demands = convert_data(input_data)
    options = input_data.get("options", {})

    model = ConcreteModel()

    # Generate concrete shifts from shift templates.
    concrete_shifts = get_concrete_shifts(shifts)

    # Determine all unique time periods in which demands occur and the shifts covering them.
    periods = get_demand_coverage_periods(concrete_shifts, demands)

    # Create variables
    x_assign = {}
    for s in concrete_shifts:
        x_assign[s["id"]] = Var(
            within=NonNegativeReals,
            bounds=(
                    s["min_workers"],
                    s["max_workers"] if s["max_workers"] >= 0 else None
                ),
            name=f'Planned_{s["id"]}'
        )
    if "under_supply_cost" in options:
        x_under = {}
        for p in periods:
            x_under[p] = Var(within=NonNegativeReals, name=f'UnderSupply_{p}')
            model.add_component(f'UnderSupply_{p}', x_under[p])
    if "over_supply_cost" in options:
        overSupply = Var(within=NonNegativeReals, name="OverSupply")
    shift_cost = Var(within=NonNegativeReals, name="ShiftCost")

    # Calculate required hours
    required_hours = sum((p.end_time - p.start_time).seconds / 3600 for p in periods)

    # Objective function
    obj_expr = sum(0 for _ in x_assign)
    obj_expr += sum(
        underSupply * options["under_supply_cost"]
        for underSupply in x_under.values()) if "under_supply_cost" in options else 0
    obj_expr += overSupply * options["over_supply_cost"] if "over_supply_cost" in options else 0
    obj_expr += shift_cost
    model.objective = Objective(expr=obj_expr, sense=minimize)

    # Constraints
    for p in periods:
        expression = sum(x_assign[s["id"]]
                         for s in p.covering_shifts) + x_under[p] if "under_supply_cost" in options else 0
        model.add_component(
            f'DemandCover_{p.start_time}_{p.end_time}_{p.qualification}',
            Constraint(expr=expression == sum(d["count"] for d in p.demands))
        )

    if "under_supply_cost" in options:
        for p in x_under:
            model.add_component(
                f'UnderSupply_{p}',
                Constraint(expr=x_under[p] * (p.end_time - p.start_time).seconds / 3600 == x_under[p]))

    if "over_supply_cost" in options:
        model.add_component(
            'OverSupply',
            Constraint(expr=overSupply == sum(
                x_assign[s["id"]] * (s["end_time"] - s["start_time"]).seconds / 3600
                for s in concrete_shifts) - required_hours)
            )

    model.add_component(
        'ShiftCost',
        Constraint(expr=shift_cost == sum(x_assign[s["id"]] * s["cost"] for s in concrete_shifts))
    )

    # Creates the solver.
    solver = SolverFactory(provider)
    solver.options[SUPPORTED_PROVIDER_DURATIONS[provider]] = duration

    # Solves the problem.
    results = solver.solve(model)

    # Convert to solution format
    has_solution = results.solver.termination_condition.is_success()
    schedule = {
        "planned_shifts": [
            {
                "id": s["id"],
                "shift_id": s["shift_id"],
                "time_id": s["time_id"],
                "start_time": s["start_time"],
                "end_time": s["end_time"],
                "qualification": s["qualification"],
                "count": round(x_assign[s["id"]].value),
            }
            for s in shifts
            if x_assign[s["id"]].value > 0.5
        ]
        if has_solution
        else [],
    }

    # Create the statistics
    statistics = {
        "result": {
            "custom": {
                "provider": "glpk",
                "status": "optimal" if has_solution else "unknown",
                "has_solution": has_solution,
                "constraints": len(model.component_objects(Constraint)),
                "variables": len(model.component_objects(Var)),
                "planned_shifts": len(schedule["planned_shifts"]),
                "planned_count": sum(s["count"] for s in schedule["planned_shifts"]),
                "shift_cost": shift_cost.value if has_solution else 0,
                "under_supply": sum(
                    x_under[p].value * (p.end_time - p.start_time).seconds / 3600
                    for p in x_under.keys()) if has_solution and "under_supply_cost" in options else 0.0,
                "over_supply": overSupply.value if has_solution and "over_supply_cost" in options else 0.0,
                "over_supply_cost": overSupply.value * options["over_supply_cost"]
                if has_solution and "over_supply_cost" in options else 0.0,
                "under_supply_cost": sum(
                    x_under[p].value
                    for p in x_under.keys()) * options["under_supply_cost"]
                    if has_solution and "under_supply_cost" in options else 0.0,
            },
            "duration": results.solver.wallclock_time / 1000,
            "value": model.objective.expr() if has_solution else None,
        },
        "run": {
            "duration": results.solver.wallclock_time / 1000,
        },
        "schema": "v1",
    }

    log(f"  - status: {statistics['result']['custom']['status']}")
    log(f"  - duration: {statistics['result']['duration']} seconds")
    log(f"  - value: {statistics['result']['value']}")
    log(f"  - planned shifts: {statistics['result']['custom']['planned_shifts']}")
    log(f"  - planned count: {statistics['result']['custom']['planned_count']}")
    log(f"  - under supply: {statistics['result']['custom']['under_supply']}")
    log(f"  - over supply: {statistics['result']['custom']['over_supply']}")
    log(f"  - shift cost: {statistics['result']['custom']['shift_cost']}")
    log(f"  - over supply cost: {statistics['result']['custom']['over_supply_cost']}")
    log(f"  - under supply cost: {statistics['result']['custom']['under_supply_cost']}")

    return {
        "solutions": [schedule],
        "statistics": statistics,
    }

class UniqueQualificationDemandPeriod:
    """
    Represents a unique time-period and qualification combination. It lists all demands
    causing the need for this qualification in this time period, as well as all shifts
    helping in covering them.
    """

    def __init__(
        self,
        start_time: datetime.datetime,
        end_time: datetime.datetime,
        qualification: str,
        covering_shifts: list[str],
        demands: list[str],
    ):
        """Creates a new unique time-period and qualification combination."""

        self.start_time = start_time
        self.end_time = end_time
        self.qualification = qualification
        self.covering_shifts = covering_shifts
        self.demands = demands

    def __str__(self) -> str:
        """Returns a string representation of this object."""

        return f"{self.start_time.isoformat()}_{self.end_time.isoformat()}_{self.qualification}"


def get_concrete_shifts(shifts: list[dict[str, Any]]) -> list[dict[str, Any]]:
    # Convert shift templates into concrete shifts. I.e., for every shift and every time
    # it can be planned, we create a concrete shift.
    # While most characteristics are given on the shift itself (except for the time), many
    # of them can be overwritten by the individual times a shift can be planned. E.g., the
    # maximum number of workers that can be assigned to a shift may be less during a night
    # shift than during a day shift.
    concrete_shifts = [
        {
            "id": f"{shift['id']}_{time['id']}",
            "shift_id": shift["id"],
            "time_id": time["id"],
            "start_time": time["start_time"],
            "end_time": time["end_time"],
            # Min workers is 0 at default. Furthermore, it can be overwritten by the individual time.
            "min_workers": time["min_workers"]
            if "min_workers" in time
            else shift["min_workers"]
            if "min_workers" in shift
            else 0,
            # Max workers is -1 at default (unbounded). Furthermore, it can be overwritten by the individual time.
            "max_workers": time["max_workers"]
            if "max_workers" in time
            else shift["max_workers"]
            if "max_workers" in shift
            else -1,
            # Cost is required. Furthermore, it can be overwritten by the individual time.
            "cost": time["cost"] if "cost" in time else shift["cost"],
            # Make sure that the qualification is present.
            "qualification": shift["qualification"] if "qualification" in shift else "",
        }
        for shift in shifts
        for time in shift["times"]
    ]
    return concrete_shifts


def get_demand_coverage_periods(
    concrete_shifts: list[dict[str, Any]], demands: list[dict[str, Any]]
) -> list[UniqueQualificationDemandPeriod]:
    """
    Determines all unique time-periods with demand for a qualification. It returns all
    demands contributing and all shifts potentially covering this time period.
    """

    # Group demands by their required qualification
    demands_per_qualification = {}
    for d in demands:
        qualification = d["qualification"] if "qualification" in d else ""
        if qualification not in demands_per_qualification:
            demands_per_qualification[qualification] = []
        demands_per_qualification[qualification].append(d)

    # Determine all concrete shifts covering a demand
    shifts_per_qualification = {}
    for q in demands_per_qualification:
        shifts_per_qualification[q] = [
            s for s in concrete_shifts if s.get("qualification") is not None and q == s["qualification"]
        ]

    # Determine all unique time periods
    periods = []
    for q in demands_per_qualification:
        # Determine all unique times for this qualification
        times = set()
        for d in demands_per_qualification[q]:
            times.add(d["start_time"])
            times.add(d["end_time"])
        for s in shifts_per_qualification[q]:
            times.add(s["start_time"])
            times.add(s["end_time"])
        times = sorted(times)

        # Create unique time periods
        for i in range(len(times) - 1):
            start, end = times[i], times[i + 1]
            # Collect all shifts covering this time period and demands contributing to it
            covering_shifts = [
                s
                for s in shifts_per_qualification[q]
                if s["start_time"] <= start and s["end_time"] >= end
            ]
            contributing_demands = [
                d
                for d in demands_per_qualification[q]
                if d["start_time"] <= start and d["end_time"] >= end
            ]
            if not any(contributing_demands):
                continue
            periods.append(
                UniqueQualificationDemandPeriod(
                    start,
                    end,
                    q,
                    covering_shifts,
                    contributing_demands,
                )
            )

    return periods


def convert_data(
    input_data: dict[str, Any],
) -> tuple[
    list[dict[str, Any]],
    list[dict[str, Any]],
]:
    """Converts the input data into the format expected by the model."""
    shifts = input_data["shifts"]
    demands = input_data["demands"]
    # In-place convert all times to datetime objects.
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


def custom_serial(obj):
    """JSON serializer for objects not serializable by default serializer."""

    if isinstance(obj, (datetime.datetime | datetime.date)):
        return obj.isoformat()
    raise TypeError("Type %s not serializable" % type(obj))


if __name__ == "__main__":
    main()
