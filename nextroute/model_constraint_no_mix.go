package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// NoMixConstraint is a type of inter-tour stop constraint that prevents certain
// stops from being visited on the same tour.
type NoMixConstraint interface {
	ModelConstraint
	// Insert returns the mix types that are associated with a stop that is
	// inserted into a tour.
	Insert() map[ModelStop]MixIngredient
	// Remove returns the mix types that are associated with a stop that
	// is removed from a tour.
	Remove() map[ModelStop]MixIngredient
}

// MixIngredient is an ingredient that is used to specify the type of mix.
// The name is the name of the mix ingredient. The count is the number units of
// the mix ingredient are inserted or removed from a tour.
type MixIngredient struct {
	// Name is the name of the mix ingredient.
	Name string
	// Count is the number units of the mix ingredient are inserted or removed from a
	// tour.
	Count int
}

// NewNoMixConstraint creates a new no-mix constraint. The constraint
// needs to be added to the model to be taken into account.
// The insert map contains the mix ingredients that are associated with a stop
// that is inserted into a tour. The remove map contains the ingredients that
// are associated with a stop that is removed from a tour. A stop can only be in
// insert or remove, not in both. The sum of the counts of the insert map must
// be equal to the sum of the counts of the remove map for each plan unit.
func NewNoMixConstraint(
	deltas map[ModelStop]MixIngredient,
) (NoMixConstraint, error) {
	connect.Connect(con, &newNoMixConstraint)
	return newNoMixConstraint(deltas)
}
