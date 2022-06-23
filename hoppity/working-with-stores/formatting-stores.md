# ![ears](../img/ears.png) Formatting Stores

Hop was designed to work well with JSON data. This makes it easy to deploy Hop
models into microservices, run them as serverless functions, and many other
things. A store can be directly encoded into JSON as a representation of its
variable assignments.

Let's create a new store and encode it into JSON.

```go
import (
    "encoding/json"
    "os"

    "github.com/nextmv-io/sdk/hop/store"
)

func main() {
    enc := json.NewEncoder(os.Stdout)

    s := store.New()
    enc.Encode(s)
}
```

This should write an empty list to standard out because the store is empty.

```json
[]
```

Let's add some variables to the store before the end of our `main` function.

```go
x := store.Var(s, 42)
y := store.Var(s, "foo")
pi := store.Var(s, 3.14)
enc.Encode(s)
```

By default, the JSON representation of a store contains its variable assignments
in order of declaration.

```txt
[
  42,
  "foo",
  3.14
]
```

This may be fine, or we may rather reshape the format into something more
convenient. Let's turn that JSON list into a map with the variable names. To do
this, we use the `Format` method. `Format` is similar to `Apply`, in that it
doesn't change the existing store, but applies a change to create a new one. The
difference is that now we are adding _logic_ to the store instead of a change
to a variable assignment. Hop stores give us very specific ways to introduce
logic.

```go
s = s.Format(func(s types.Store) any {
    return map[string]any{
        "x":  x.Get(s),
        "y":  y.Get(s),
        "pi": pi.Get(s),
    }
})
enc.Encode(s)
```

Whoops! Where does that `types.Store` come from? It looks like we need a new
import statement. Luckily this also comes from the SDK. Let's change our imports
so they look like this.

```go
import (
    "encoding/json"
    "os"

    "github.com/nextmv-io/sdk/hop/store"
    "github.com/nextmv-io/sdk/hop/store/types"
)
```

Now we should see our store encoded as a map. This isn't that far from where we
started, but is more useful! Since `Format` can reshape a store into anything
that can be encoded in JSON, its easy to make Hop output match anything we
might expect in a production environment.

```json
{
  "pi": 3.14,
  "x": 42,
  "y": "foo"
}
```

Finally, let's apply a change to a variable assignment and encode the resulting
store.

```go
enc.Encode(s.Apply(y.Set("bar")))
```

Note how the new store inherits our formatting logic.

```json
{
  "pi": 3.14,
  "x": 42,
  "y": "bar"
}
```

## Exercises

* Start with the [example above][source] and reshape the output into something
  more complex than a map. Can you format it as a map of maps or as a
  user-defined structure? How do the rules of the [`encoding/json`][json]
  library apply to the output?
* What happens if you override the formatting logic of a child store? Does that
  impact the parent or any sibling stores?

[json]:   https://pkg.go.dev/encoding/json
[source]: formatting-stores/main.go
