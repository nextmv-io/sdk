package nextroute

// ModelPlanUnit is a plan unit. It is a unit defining what should be planned .
// For example, a unit can be a pickup and a delivery stop that are required to
// be planned together on the same vehicle.
type ModelPlanUnit interface {
	ModelData

	// Index returns the index of the invoking unit.
	Index() int

	// IsFixed returns true if the PlanUnit is fixed.
	IsFixed() bool
}

// ModelPlanUnits is a slice of plan units .
type ModelPlanUnits []ModelPlanUnit
