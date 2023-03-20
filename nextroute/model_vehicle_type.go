package nextroute

// ModelVehicleType is a vehicle type. A vehicle type is a definition of a
// vehicle. It contains the process duration and travel duration expressions
// that are used to calculate the travel and process duration of a stop
// assignment to a vehicle of this type.
type ModelVehicleType interface {
	// Data returns the arbitrary data associated with the vehicle type. Can be
	// set using the VehicleTypeData VehicleTypeOption in the factory method
	// Model.NewVehicleType.
	Data() any
	// Index returns the index of the vehicle type.
	Index() int

	// Model returns the model of the vehicle type.
	Model() Model

	// ProcessDurationExpression returns the process duration expression of the
	// vehicle type. Is set in the factory method of the vehicle type
	// Model.NewVehicleType.
	ProcessDurationExpression() DurationExpression

	// SetData sets the arbitrary data associated with the vehicle type.
	SetData(data any)
	// TravelDurationExpression returns the travel duration expression of the
	// vehicle type. Is set in the factory method of the vehicle type
	// Model.NewVehicleType.
	TravelDurationExpression() TravelDurationExpression

	// Vehicles returns the vehicles of this vehicle type.
	Vehicles() ModelVehicles
}

// ModelVehicleTypes is a slice of vehicle types.
type ModelVehicleTypes []ModelVehicleType
