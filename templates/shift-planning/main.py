"""
Template for working with Google OR-Tools.
"""

import argparse
import datetime
import json
import sys
from typing import Any

from ortools.linear_solver import pywraplp

# Status of the solver after optimizing.
STATUS = {
    pywraplp.Solver.FEASIBLE: "feasible",
    pywraplp.Solver.INFEASIBLE: "infeasible",
    pywraplp.Solver.OPTIMAL: "optimal",
    pywraplp.Solver.UNBOUNDED: "unbounded",
}
ANY_SOLUTION = [pywraplp.Solver.FEASIBLE, pywraplp.Solver.OPTIMAL]


def main() -> None:
    """Entry point for the template."""

    parser = argparse.ArgumentParser(
        description="Solve shift-creation with OR-Tools MIP."
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
    args = parser.parse_args()

    # Read input data, solve the problem and write the solution.
    input_data = read_input(args.input)
    log("Solving shift-creation:")
    log(f"  - shifts-templates: {len(input_data.get('shifts', []))}")
    log(f"  - demands: {len(input_data.get('demands', []))}")
    log(f"  - max duration: {args.duration} seconds")
    solution = solve(input_data, args.duration)
    write_output(args.output, solution)


def solve(input_data: dict[str, Any], duration: int) -> dict[str, Any]:
    """Solves the given problem and returns the solution."""

    # Creates the solver.
    provider = "SCIP"
    solver = pywraplp.Solver.CreateSolver(provider)
    solver.SetTimeLimit(duration * 1000)

    # Prepare data
    shifts, demands = convert_data(input_data)

    # Generate concrete shifts from shift templates.
    concrete_shifts = get_concrete_shifts(shifts)

    # Determine all unique time periods in which demands occur and the shifts covering them.
    periods = get_coverage(concrete_shifts, demands)

    # Create integer variables indicating how many times a shift is planned.
    x_assign = {}
    for s in concrete_shifts:
        x_assign[(s["id"])] = solver.IntVar(
            s["min_workers"],
            s["max_workers"] if s["max_workers"] >= 0 else solver.infinity(),
            f'Planned_{s["id"]}',
        )
    # x_over, x_under = {}, {}
    # for p in periods:
    #     x_over[p] = solver.IntVar(0, solver.infinity(), f"Over_{p}")
    #     x_under[p] = solver.IntVar(0, solver.infinity(), f"Under_{p}")

    # Objective function: minimize the cost of the planned shifts
    solver.Minimize(
        solver.Sum([x_assign[s["id"]] * s["cost"] for s in concrete_shifts])
        # + solver.Sum(
        #     [
        #         x_over[p] * (p.end_time - p.start_time) * p.demands[0]["over_cost"]
        #         for p in periods
        #     ]
        # )
    )

    # TODO: over-supply / under-supply costs

    # >> Constraints

    # We need to make sure that all demands are covered
    for p in periods:
        solver.Add(
            solver.Sum([x_assign[s["id"]] for s in p.covering_shifts])
            # + x_over[p]
            # - x_under[p]
            >= sum(d["count"] for d in p.demands),
            f"DemandCover_{p.start_time}_{p.end_time}_{p.qualification}",
        )

    # Solves the problem.
    status = solver.Solve()

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
                "count": int(round(x_assign[(s["id"])].solution_value())),
            }
            for s in concrete_shifts
            if x_assign[(s["id"])].solution_value() > 0.5
        ]
        if status in ANY_SOLUTION
        else [],
    }

    # Creates the statistics.
    statistics = {
        "result": {
            "custom": {
                "constraints": solver.NumConstraints(),
                "provider": provider,
                "status": STATUS.get(status, "unknown"),
                "variables": solver.NumVariables(),
            },
            "duration": solver.WallTime() / 1000,
            "value": solver.Objective().Value() if status in ANY_SOLUTION else None,
        },
        "run": {
            "duration": solver.WallTime() / 1000,
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


def get_coverage(
    concrete_shifts: list[dict[str, Any]], demands: list[dict[str, Any]]
) -> list[UniqueTimeQualificationPeriod]:
    """
    Determines all unique time-period and qualification combinations. Furthermore, returns
    the demands contributing to each one of them as well as the shifts covering them.
    """

    # Determine all unique times
    times = set()
    for d in demands:
        times.add(d["start_time"])
        times.add(d["end_time"])
    for s in concrete_shifts:
        times.add(s["start_time"])
        times.add(s["end_time"])
    times = sorted(times)

    # Create unique time periods
    periods = []
    for i in range(len(times) - 1):
        start, end = times[i], times[i + 1]
        # Determine demands contributing to this time period
        associated_demands = [
            d for d in demands if d["start_time"] <= start and d["end_time"] >= end
        ]
        # Determine all qualification dimensions we need to consider, use empty string for
        # no qualification required
        qualifications = set()
        for d in associated_demands:
            qualifications.add(d["qualification"] if "qualification" in d else "")
        # Create a unique time period for each qualification dimension
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

    if isinstance(obj | (datetime.datetime, datetime.date)):
        return obj.isoformat()
    raise TypeError("Type %s not serializable" % type(obj))


if __name__ == "__main__":
    main()
