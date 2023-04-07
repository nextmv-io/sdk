package nextroute

// ModelVehicleType is a vehicle type. A vehicle type is a definition of a
// vehicle. It contains the process duration and travel duration expressions
// that are used to calculate the travel and process duration of a stop
// assignment to a vehicle of this type.
type ModelVehicleType interface {
	ModelData

	// CalculateTemporalValues calculates the arrival, start and end time of
	// stop if it would be assigned after previousStop to this vehicle type, and
	// it would leaf from previousStop at previousEnd.
	CalculateTemporalValues(
		previousEnd float64,
		previousStop ModelStop,
		stop ModelStop,
	) (arrival, start, end float64)

	// Index returns the index of the vehicle type.
	Index() int

	// Model returns the model of the vehicle type.
	Model() Model

	// Name returns the name of the vehicle.
	Name() string

	// ProcessDurationExpression returns the process duration expression of the
	// vehicle type. Is set in the factory method of the vehicle type
	// Model.NewVehicleType.
	ProcessDurationExpression() DurationExpression

	// SetName sets the name of the vehicle.
	SetName(string)

	// TravelDurationExpression returns the duration expression of the
	// vehicle type. Is set in the factory method of the vehicle type
	// Model.NewVehicleType.
	TravelDurationExpression() DurationExpression

	// Vehicles returns the vehicles of this vehicle type.
	Vehicles() ModelVehicles
}

// ModelVehicleTypes is a slice of vehicle types.
type ModelVehicleTypes []ModelVehicleType
