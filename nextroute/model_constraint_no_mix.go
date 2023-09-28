package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// NoMixConstraint limits the order in which stops are assigned to a vehicle
// based upon the ingredients the stops insert or remove from a vehicle.
type NoMixConstraint interface {
	ModelConstraint
	// Insert returns the mix ingredients that are associated with a stop that
	// inserts an ingredient into a vehicle.
	Insert() map[ModelStop]MixIngredient
	// Remove returns the mix ingredients that are associated with a stop that
	// removes an ingredient from a vehicle.
	Remove() map[ModelStop]MixIngredient
}

// MixIngredient is an ingredient that is used to specify the type of mix.
// The name is the name of the mix ingredient. The count is the number units of
// the mix ingredient are inserted or removed from a vehicle.
type MixIngredient struct {
	// Name is the name of the mix ingredient.
	Name string
	// Count is the number units of the mix ingredient are inserted or removed from a
	// tour.
	Count int
}

// NewNoMixConstraint creates a new no-mix constraint. The constraint
// needs to be added to the model to be taken into account.
// The deltas map contains the information defining how much ingredient a stop
// inserts to the mix on a vehicle. If the count is positive it adds to the mix,
// if the count is negative it removes from the mix. A stop can only insert or
// remove an ingredient from the mix. The constraint limits the number of
// ingredients that can be on the vehicle to be at most one ingredient. The sum
// of the counts of all the stops of each plan unit must be zero.
func NewNoMixConstraint(
	deltas map[ModelStop]MixIngredient,
) (NoMixConstraint, error) {
	connect.Connect(con, &newNoMixConstraint)
	return newNoMixConstraint(deltas)
}
