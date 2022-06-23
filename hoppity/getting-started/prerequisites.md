# ![ears](../img/ears.png) Prerequisites

To get started you need Go 1.18.3. If you don't have Go installed, you can get 
it [here][download]. If you already have Go, you can see which version it is by
running `go version`.

```bash
$ go version
go version go1.18.3 darwin/amd64
```

If you have a different version installed, you can install 1.18.3 [using the
`go` you already have][manage].

```bash
$ go install golang.org/dl/go1.18.3@latest

$ go1.18.3 download

$ go1.18.3 version
go version go1.18.3 darwin/amd64
```

You also need the following shared object files for your architecture and 
operating system from Nextmv. Put these in a directory and add a 
`NEXTMV_LIBRARY_PATH` environment variable pointing to them. Your setup should 
look something like this.

```bash
$ ls ~/.nextmv/lib 
nextmv-run-cli-v0.16.0-dev.0-2-go1.18.3-darwin-amd64.so
nextmv-run-http-v0.16.0-dev.0-2-go1.18.3-darwin-amd64.so
nextmv-sdk-v0.16.0-dev.0-2-go1.18.3-darwin-amd64.so

$ echo $NEXTMV_LIBRARY_PATH               
~/.nextmv/lib
```

[download]: https://go.dev/dl/
[manage]:   https://go.dev/doc/manage-install
