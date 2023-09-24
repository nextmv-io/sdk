package nextroute

import (
	"math/rand"

	"github.com/nextmv-io/sdk/connect"
)

// NewSolutionPlanUnitCollection creates a new [SolutionPlanUnitCollection].
// A [SolutionPlanUnitCollection] is a collection of solution plan units. It
// can be used to randomly draw a sample of solution plan unit and remove
// solution plan unit from the collection.
func NewSolutionPlanUnitCollection(
	source *rand.Rand,
	planUnits SolutionPlanUnits,
) SolutionPlanUnitCollection {
	connect.Connect(con, &newSolutionPlanUnitCollection)
	return newSolutionPlanUnitCollection(source, planUnits)
}

// ImmutableSolutionPlanUnitCollection is a collection of solution plan
// units.
type ImmutableSolutionPlanUnitCollection interface {
	// Iterator returns a channel that can be used to iterate over the solution
	// plan units in the collection.
	// If you break out of the for loop before the channel is closed,
	// the goroutine launched by the Iterator() method will be blocked forever,
	// waiting to send the next element on the channel. This can lead to a
	// goroutine leak and potentially exhaust the system resources. Therefore,
	// it is recommended to always use the following pattern:
	//    iter := collection.Iterator()
	//    for {
	//        element, ok := <-iter
	//        if !ok {
	//            break
	//        }
	//        // do something with element, potentially break out of the loop
	//    }
	//    close(iter)
	Iterator(quit <-chan struct{}) <-chan SolutionPlanUnit
	// RandomDraw returns a random sample of n different solution plan units.
	RandomDraw(n int) SolutionPlanUnits
	// RandomElement returns a random solution plan unit.
	RandomElement() SolutionPlanUnit
	// Size return the number of solution plan units in the collection.
	Size() int
	// SolutionPlanUnit returns the solution plan units in the collection
	// which correspond to the given model plan unit. If no such solution
	// plan unit is found, nil is returned.
	SolutionPlanUnit(modelPlanUnit ModelPlanUnit) SolutionPlanUnit
	// SolutionPlanUnits returns the solution plan units in the collection.
	// The returned slice is a defensive copy of the internal slice, so
	// modifying it will not affect the collection.
	SolutionPlanUnits() SolutionPlanUnits
}

// SolutionPlanUnitCollection is a collection of solution plan units.
type SolutionPlanUnitCollection interface {
	ImmutableSolutionPlanUnitCollection
	// Add adds a [SolutionPlanUnit] to the collection.
	Add(solutionPlanUnit SolutionPlanUnit)
	// Remove removes a [SolutionPlanUnit] from the collection.
	Remove(solutionPlanUnit SolutionPlanUnit)
}
