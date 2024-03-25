# Nextmv's Software Development Kit

Nextmv's SDK is a collection of Go APIs for solving decision automation
problems. Please find the following packages:

- [run][run]: convenient runners that read an input, run a solver and write an
      output.
- [measure][measure]: measures for various distances between locations.
- [golden][golden]: tools for running tests with golden files.

Please visit the official [Nextmv docs][docs] for comprehensive information.

## Versioning

We try our best to version our software thoughtfully and only break APIs and
behaviors when we have a good reason to.

- Minor (`v1.^.0`) tags: new features, might be breaking.
- Patch (`v1.0.^`) tags: bug fixes.

[run]: ./run
[measure]: ./measure
[golden]: ./golden
[docs]: https://docs.nextmv.io
