# Nextmv unique matrix measure template

The `unique-matrix-measure` template shows you how to use a smaller input file
in cases where you need to pass a matrix for many stops that are using few
unique locations.

The most important files created are `main.go` and `input.json`.

`main.go` reads a unique set of locations and a matrix that represents the costs
for going from each unique location to another. `input.json` is a sample input
file that follows the input definition in `main.go`.

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
