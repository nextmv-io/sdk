# Nextmv cluster template

The most important files created are `main.go` and `input.json`.

* `main.go` implements a cluster solver.
* `input.json` is a sample input file that follows the input definition in
`main.go`.

Before you start customizing, run the command below to see if everything works
as expected:

```bash
nextmv run local . -- -runner.input.path input.json \
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with a clustering solution

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
