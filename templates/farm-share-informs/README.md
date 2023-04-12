# Nextmv INFORMS Analytics 2023 routing template

The INFORMS template sets up a vehicle routing problem (VRP) solver that is
compatible with our [cloud](https://docs.nextmv.io/cloud/get-started) input
files. In addition to almost being a drop-in replacement for the cloud endpoint
it also demonstrates some of the more advanced router options.

The most important files created are `main.go` and several data input files for
different uses cases: fleet management, delivery, distribution and sourcing. By
default the workspace file points to `data/denv_s.json`. In addition there is the
`schema.go` file that defines needed data structures and the `helper.go` file in
which helper functions are defined, e.g. for data handling.

`main.go` is the entry point for the VRP solver. The actual configuration can be
found in `solver.go`.

Before you start customizing run the command below to see if everything works as
expected:

``` bash
nextmv sdk run . -- -runner.input.path data/denv_s.json \
  -runner.output.path output.json -solveoptions.maximumduration 10s
```

A file `output.json` should have been created with a VRP solution.

## App management and experiments

Two bash scripts have been provided for you to use during the conference
(`part1` and `part2`). The scripts will generate and run the required Nextmv CLI
commands for apps and experiments.

## Next steps

* For more information about our platform, please visit our [docs][docs].
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!

[docs]: https://docs.nextmv.io
