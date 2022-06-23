# ![ears](../img/ears.png) Updating Data

We can think of a store as a lexical scope containing variable declarations,
variable assignments, and logic. Thus a store has similar mechanics to a block
in a lexically scoped programming language without destructive assignment. For 
example, say we have outer and inner blocks with the following assignments.

```txt
{
	x = 42
	y = "foo"

	{
		y = "bar"
		pi = 3.14
	}
}
```

The outer block contains two variables, `x = 42` and `y = "foo"`. The inner 
block inherits `x = 42` from the outer block which contains it, overrides
`y = "bar"`, and adds a new variable `pi = 3.14`. Assignments in the inner
block do not impact the outer block.

Let's code this example up using Hop. First we define a store `s1` and add `x`
and `y` to it.

```go
s1 := store.New()
x := store.Var(s1, 42)
y := store.Var(s1, "foo")
```

Now we apply a change set to `s1`. This results in a new store, `s2`. `s2` is 
functionally a copy of `s1` with a new value associated with `y` and a new
variable, `pi`.


```go
s2 := s1.Apply(y.Set("bar"))
pi := store.Var(s2, 3.14)
```

If we query

```go
fmt.Println("s1:", x.Get(s1), y.Get(s1))
fmt.Println("s2:", x.Get(s2), y.Get(s2), pi.Get(s2))
```

Run the [source][source] above and you should see output similar to this.

```txt
s1: 42 foo
s2: 42 bar 3.14
```

Note that calling `y.Set("bar")` creates a `Change`. The `Apply` method accepts
any number of changes and applies them in order to create a new store. Thus the
code below creates a single store with `x = 10` and `y = "abc"`.

```go
s1.Apply(
	x.Set(-3),
	y.Set("bar"),
	y.Set("abc"),
	x.Set(100),
	x.Set(10),
)
```

## Exercises

* Apply multiple changes to `s1` to create a new store `s3`. Does `s3` impact
  `s2` in any way?
* Try to store a slice on `s1` then `append` to it when applying changes to
  create both `s2` and `s3`. What do you expect to happen? What happens?

[source]: updating-data/main.go
