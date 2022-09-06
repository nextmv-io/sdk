# Store

Package `store` provides a modeling kit for decision automation problems. It is
based on the paradigm of "decisions as code". The base interface is the
[Store][Store]: a space defined by variables and logic. The underlying
algorithms search that space and find the best solution possible, this is, the
best collection of variable assignments. The Store is the root node of a search
tree. Child Stores (nodes) inherit both logic and variables from the parent and
may also add new variables and logic, or overwrite existing ones. Changes to a
child do not impact its parent.

See [godocs][godocs] for package docs.

[godocs]:  https://pkg.go.dev/github.com/nextmv-io/sdk/store
[Store]: ./store.go
