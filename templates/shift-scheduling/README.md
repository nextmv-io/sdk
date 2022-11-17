# Nextmv shift scheduling template

The [Nurse scheduling
problem](https://en.wikipedia.org/wiki/Nurse_scheduling_problem) asks to assign
nurses or more general workers to shifts. In our example problem we have the
following setting:

* We have a number of days. Each day has three shift types: `morning`, `day` and
  `night`.
* We have a number of workers. Each shift should have a worker, but not every
  worker needs to be assigned to a shift.
* We further want to ensure that each worker has a break of at least two shifts
  between assignments.
* In addition workers can be unavailable for certain full days.
* At last, workers can have preferences as to what shift type they prefer.
* Our objective is to find a plan that minimizes the worker count and the worker
  happiness in a 10:1 ratio. Happiness is measured as the number of times the
  worker had their preferred shift assigned.

This template uses our *custom modelling* framework to model and solve
such a shift scheduling problem.

The most important files created are `main.go` and `input.json`.

* `main.go` implements a Shift scheduling solver.
* `input.json` is a sample input file that follows the input definition in
`main.go`.

Before you start customizing, run the command below to see if everything works
as expected:

```bash
nextmv run local . -- -hop.runner.input.path input.json \
  -hop.runner.output.path output.json -hop.solver.limits.duration 10s
```

A file `output.json` should have been created with the best found schedule.

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
