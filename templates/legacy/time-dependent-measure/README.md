# Nextmv time dependent measure template

In real world vehicle routing problems (VRP) it is often needed to take the time
of the day into account when calculating routes, e.g. to calculate correct ETAs
during rush hours.
The `time dependent measure` template shows you how to make use of time
dependent measures and apply different costs according to time. This template is
based on the `routing` template and can easily be combined with more options
that are shown there like time windows or capacities (among many more).

The most important files created are `main.go` and `input.json`.

`main.go` implements a VRP solver with many real world features already
configured. `input.json` is a sample input file that follows the input
definition in `main.go`.

Before you start customizing run the command below to see if everything works as
expected:

```bash
nextmv sdk run . -- -runner.input.path input.json\
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with a VRP solution.

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
