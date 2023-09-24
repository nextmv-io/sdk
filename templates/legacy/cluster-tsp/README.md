# Nextmv cluster-tsp template

The cluster-tsp template reads a routing input data format but instead of
solving a vehicle routing problem directly, it creates a set of clusters in a
pre-processing step. The number of clusters created match the number of vehicles
available in the input file, the cluster size is calculated such that points are
distributed evenly among them.
Then for each vehicle a TSP is solved. This is achieved by using the `Attributes`
option in the routing engine.

The most important files created are `solver.go` and an input file which
represents the instance to be solved. In addition there is the `schema.go` file
that defines needed data structures and the `helper.go` file in which helper
functions are defined, e.g. for data handling.

`main.go` is the entry point for the solver. The actual configuration can be
found in `solver.go`.

Before you start customizing run the command below to see if everything works as
expected:

``` bash
nextmv sdk run . -- -runner.input.path input.json\
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with a solution.

## Next steps

* For more information about our platform, please visit our [docs][docs].
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!

[docs]: https://docs.nextmv.io
