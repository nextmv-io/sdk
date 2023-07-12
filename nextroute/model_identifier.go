package nextroute

// Identifier is an interface that can be used for identifying objects.
type Identifier interface {
	// ID returns the identifier of the object.
	ID() string
}

// IdentifierWriteable is an interface that can be used for identifying objects.
// Furthermore, it allows the identifier to be modified.
type IdentifierWriteable interface {
	Identifier
	// SetID sets the identifier of the object.
	SetID(string)
}
