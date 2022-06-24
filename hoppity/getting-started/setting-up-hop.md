# ![ears](../img/ears.png) Setting up Hop

In order to run the code in this tour, you should use [a Go module][modules]. A
module manages your dependencies, including the Nextmv SDK. Let's create one
called `hoppity`.

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
hoppity$ cat << EOF > main.go
package main

import (
    "fmt"

    "github.com/nextmv-io/sdk"
    "github.com/nextmv-io/sdk/hop/store"
)

func main() {
    s := store.New()
    version := store.Var(s, sdk.VERSION)
    fmt.Println("Hello Hop", version.Get(s))
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
examples in this tour constitutes a complete `main.go` Put them in unique
directories  inside your `hoppity` folder and run them using the same `go run`
command shown above.

[modules]: https://go.dev/blog/using-go-modules
