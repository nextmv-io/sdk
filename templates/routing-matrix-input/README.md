# Nextmv routing matrix input template

`routing-matrix-input` is a modeling kit for vehicle routing problems (VRP).
This template shows how to use a custom matrix that is passed via an input file.

The most important files created are `main.go` and `input.json`.

`main.go` implements a VRP solver using a custom measure matrix from the input
file. `input.json` is a sample input file that follows the input definition in
`main.go`.

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
