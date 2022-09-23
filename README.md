# Nextmv's Software Development Kit

Nextmv's SDK is a collection of Go APIs for solving decision automation
problems. Please find the following packages:

- [store][store]: an all-purpose modeling kit for decision automation problems,
      serving as the core of Nextmv's SDK.
- [route][route]: a modeling kit for vehicle routing problems.
- [run][run]: convenient runners that read an input, run a solver and write an
      output.
- [model][model]: modeling components such as integer domains and ranges.

Please visit the official [Nextmv docs][docs] for comprehensive information.

## Installation

Nextmv's SDK is meant to be used in Go projects. To download please run:

```bash
go get github.com/nextmv-io/sdk
```

[store]: ./store
[route]: ./route
[run]: ./run
[model]: ./model
[docs]: https://docs.nextmv.io
