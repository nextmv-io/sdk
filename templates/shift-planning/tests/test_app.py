import argparse
import collections
import os
import pathlib
import pytest
import re
import subprocess
import sys

READY = False
UPDATE = False
OUTPUT_DIR = None
DATA_DIR = None
APP_PATH = None

AppTest = collections.namedtuple(
    "Test",
    [
        "name",
        "args",
        "input",
        "output",
        "golden_output",
    ],
)


@pytest.fixture(autouse=True)
def run_around_tests() -> None:
    """
    Prepare tests and clean up.
    """
    global READY, UPDATE, OUTPUT_DIR, DATA_DIR, APP_PATH
    # If it's the first test, setup all variables
    if not READY:
        # Read arguments
        parser = argparse.ArgumentParser(description="app tests")
        parser.add_argument(
            "--update",
            dest="update",
            action="store_true",
            default=False,
            help="updates the golden files",
        )
        args, _ = parser.parse_known_args()  # Ignore potentially forwarded pytest args
        UPDATE = args.update

        # Update if requested by env var
        if os.environ.get("UPDATE", "0") == "1":
            UPDATE = True

        # Set paths
        OUTPUT_DIR = str(pathlib.Path(__file__).parent.joinpath("./output").resolve())
        DATA_DIR = str(pathlib.Path(__file__).parent.joinpath("./testdata").resolve(strict=True))
        APP_PATH = str(pathlib.Path(__file__).parent.joinpath("../main.py").resolve(strict=True))

        # Prepare output directory
        os.makedirs(OUTPUT_DIR, exist_ok=True)

        # Mark as ready
        READY = True

    # Run a test
    yield

    # Clean up
    # - nothing to do, we keep the output for manual inspection


def _stabilize(data: str) -> str:
    """
    Removes volatile output to ensure a stable comparison.
    """
    data = re.sub(r'"duration":.*[0-9.]+', '"duration": 0.123', data)
    return data


def _run_test(test: AppTest) -> None:
    # Clear old results
    if os.path.isfile(test.output):
        os.remove(test.output)

    # Assemble command and arguments
    base = [sys.executable, APP_PATH]
    cmd = [*base, *test.args]

    # Log
    cmd_string = " ".join(test.args)
    print(f"Invoking: {cmd_string}")

    # Run command
    result = subprocess.run(cmd, stdout=subprocess.PIPE, stdin=open(test.input, "rb"))

    # Expect no errors
    assert result.returncode == 0

    # Make sure stderr is empty
    assert result.stderr is None or result.stderr.decode("utf-8") == ""

    # Expect solution on stdout
    assert result.stdout is not None
    output = result.stdout.decode("utf-8")
    output = _stabilize(output)

    # Write output to file (for manual inspection)
    with open(test.output, "w") as f:
        f.write(output)

    if UPDATE:
        # Update golden output file
        with open(test.golden_output, "w") as fw:
            fw.write(output)
    else:
        # Compare output with golden output file
        expected = ""
        with open(test.golden_output, "r") as f:
            expected = f.read()
        assert output == expected


def test_sample_input():
    test = AppTest(
        "sample-input",
        [],
        os.path.join(DATA_DIR, "input.json"),
        os.path.join(OUTPUT_DIR, "input.output.json"),
        os.path.join(DATA_DIR, "input.json.golden"),
    )
    _run_test(test)


if __name__ == "__main__":
    run_around_tests()
    test_sample_input()
    print("Everything passed")
