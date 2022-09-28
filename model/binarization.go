package model

// BinaryHelper is an int with an ID method to implement the Identifier
// interface.
type BinaryHelper int

// ID returns the id of a BinaryHelper.
func (b BinaryHelper) ID() int {
	return int(b)
}

// Binarized is a helper type returned by the Binarize function. It can be
// queried for the Identifier belonging to a specific number.
type Binarized []Identifier

// GetIdentifier will return the Identifier that belongs to a number.
func (b Binarized) GetIdentifier(number int) Identifier {
	return b[number-1]
}

// Binarize will take a number and return a Binarized (slice of Identifiers).
func Binarize(number int) Binarized {
	returnList := make([]Identifier, number)
	for i := 0; i < number; i++ {
		returnList[i] = BinaryHelper(i)
	}
	return returnList
}
