# Hopping About: A Tour of Data Store Modeling with Hop

```
Christopher Robin goes
Hoppity, hoppity,
Hoppity, hoppity, hop.
Whenever I tell him
Politely to stop it, he
Says he can't possibly stop.

If he stopped hopping,
He couldn't go anywhere,
Poor little Christopher
Couldn't go anywhere...
That's why he always goes
Hoppity, hoppity,
Hoppity,
Hoppity,
Hop.

-- A. A. Milne
```

Nextmv's decision modeling interface lets you automate any operational decision 
in a way that looks and feels like writing other code. Hop provides the 
guardrails to turn your data into automated decisions, and test and deploy them 
into production environments.

Let's see how.

## Getting Started

### Prerequisites

To get started you need Go 1.18.3. If you don't have Go installed, you can get 
it [here](https://go.dev/dl/). If you already have Go, you can see which version 
it is by running `go version`.

```bash
$ go version
go version go1.18.3 darwin/amd64
```

If you have a different version installed, you can install 1.18.3 [using the 
`go` you already have](https://go.dev/doc/manage-install).

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

### Setup

In order to run these examples in this tour, you should use
[a Go module](https://go.dev/blog/using-go-modules). A module manges your 
dependencies, including the Nextmv SDK. Let's create one called `hoppity`.

```bash
$ mkdir hoppity

$ cd hoppity

hoppity$ go mod init hoppity
go: creating new go.mod: module hoppity
```

Now we add Nextmv's SDK to our dependencies.

```bash
hoppity$ go get github.com/nextmv-io/sdk@v0.16.0-dev.0-2
go: added github.com/nextmv-io/sdk v0.16.0-dev.0-2
```

You should now have a `go.mod` file that looks like this.

```bash
hoppity$ cat go.mod
module hoppity

go 1.18

require github.com/nextmv-io/sdk v0.16.0-dev.0-2 // indirect
```

Now we can create a test file that prints Hop's version.

```bash
hoppity$ cat << EOF > ehlo/main.go
package main

import (
        "fmt"

        "github.com/nextmv-io/sdk"
)

func main() {
        fmt.Println("Hello Hop", sdk.VERSION)
}
EOF
```

We can run it using `go run`.

```bash
hoppity$ go run -trimpath ehlo/main.go 
go: downloading github.com/nextmv-io/sdk v0.16.0-dev.0-2
Hello Hop v0.16.0-dev.0-2
```

If you see see output like the above, you're ready to get hopping! Each of the 
examples linked to below is a complete `main.go` Put them in unique directories 
inside your `hoppity` folder and run them using the same `go run` command above.

## Examples

### Storing Data

1. [Create a Store](01-storing-data/01-create-a-store/README.md)
1. [Store a Custom Type](01-storing-data/02-store-a-custom-type/README.md)
