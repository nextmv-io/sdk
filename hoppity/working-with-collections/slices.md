# ![ears](../img/ears.png) Slices

Go passes arguments to functions and makes variable assignments
[by value][value]. This applies to stores as well: a store will assign to a
declared variable _by value_. If that value is of a type like an `int` or a
`struct`, then it is copied before assigning it to a variable in a store. This
is unsurprising.

It gets a little hairier when assigning a pointer type like a slice or a map to
a variable. The pointer value is copied, but not the data it refers to. While
Hop doesn't keep you from storing pointers, it's usually not what you want.
Instead, Hop provides immutable containers so you can update rich data
structures in a safe and efficient manner in your model.

The simplest of these container types is an immutable slice. We can store a
slice using the special `Slice` function.

```go
s1 := store.New()
x := store.Slice(s1, "c", "d", "e") // []string{"c", "d", "e"}
```

We can assign any type to the slice. In this case, Hop infers that we are
storing a slice of strings with initial value `["c", "d", "e"]` from the
parameters. If we want to create an empty slice, Hop needs to know what type we
intend to store.

```go
x := store.Slice[string](s1) // []string{}
```

Our slice `x` has several methods that query properties of the slice for a given
store, such as `Get`, which retrieves the value of an index, and `Len`, which
returns the length of the slice. The `Slice` method returns the underlying slice
data as a standard Go slice, like `[]string`.

It also provides methods for creating changes. We may apply these changes to
create new stores.

```go
s2 := s1.Apply(
    x.Append("h", "i", "j"),
    x.Prepend("a", "y", "z"),
)

s3 := s2.Apply(
    x.Insert(6, "f", "g"),
    x.Remove(2, 2),
    x.Set(1, "b"),
)
```

## Exercises

* Try to guess what `s1`, `s2`, and `s3` contain. Run the [source][source] and
  see if you are right.
* Create a slice of a user-defined type. Insert values into it and retrieve the
  underlying slice contents.
* Use the `Len` and `Get` methods to iterate over a slice and print its values
  one at a time.

[source]: slices/main.go
[value]: https://go.dev/doc/faq#pass_by_value
