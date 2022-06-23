# ![ears](../img/ears.png) Storing Custom Data

A Hop store can manage any concrete type, including custom structs. Let's define
a `bunny` type with a few fields and a `String` method.

```go
type bunny struct {
	name       string
	fluffiness float64
	activities []string
}

func (b bunny) String() string {
	fluffy := "fluffy"
	if b.fluffiness < 0.5 {
		fluffy = "not fluffy"
	}
	return fmt.Sprintf("%s is %s and likes %v", b.name, fluffy, b.activities)
}
```

Now we create a bunny and add it to our store.

```go
s := store.New()

peter := store.Var(s, bunny{
	name:       "Peter Rabbit",
	fluffiness: 0.52,
	activities: []string{
		"stealing and eating vegetables",
		"losing his jacket and shoes",
	},
})
```

If we retrieve the value of `peter` from our store and print it, we should get a
the results of the `bunny.String` method. Note that `peter.Get` returns a value
of type `bunny`, so Go knows to call the appropriate `String` method when we
pass it to `fmt.Println`.

```go
fmt.Println(peter.Get(s))
```

Run the [source][source] above and you should see an educational message about
this Peter Rabbit fellow.

## Exercises

* Add another bunny to the store and call its `String` method directly.
* Add another method to the `bunny` type. Retrieve `peter` from the store and
  call this new method.

[source]: storing-custom-data/main.go
