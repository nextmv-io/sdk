package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// NoMixConstraint limits the order in which stops are assigned to a vehicle
// based upon the items the stops insert or remove from a vehicle.
type NoMixConstraint interface {
	ModelConstraint
	// Insert returns the mix ingredients that are associated with a stop that
	// inserts an ingredient into a vehicle.
	Insert() map[ModelStop]MixItem
	// Remove returns the mix ingredients that are associated with a stop that
	// removes an ingredient from a vehicle.
	Remove() map[ModelStop]MixItem
}

// MixItem is an item that is used to specify the type of mix.
// The type defines the type of each item. The count is the number units of
// the item are inserted or removed from a vehicle.
type MixItem struct {
	// Name is the name of the mix item.
	Type string
	// Count is the number units of the mix items are inserted or removed from a
	// vehicle.
	Count int
}

// NewNoMixConstraint creates a new no-mix constraint. The constraint
// needs to be added to the model to be taken into account.
// The deltas map contains the information defining how many items a stop
// inserts or removes from a vehicle. If the count is positive it inserts items
// to the vehicle, if the count is negative it removes items from the vehicle.
// A stop can be present once in deltas. The constraint makes sure that at no
// point in time there are items on the vehicle of more than one type. The sum
// of the counts of all the deltas of the stops of each plan unit must be zero.
// The items of all the stops of a plan unit must be the same type if they have
// a delta.
func NewNoMixConstraint(
	deltas map[ModelStop]MixItem,
) (NoMixConstraint, error) {
	connect.Connect(con, &newNoMixConstraint)
	return newNoMixConstraint(deltas)
}
