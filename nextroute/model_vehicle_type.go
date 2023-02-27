package nextroute

import "github.com/nextmv-io/sdk/connect"

// ModelVehicleType is a vehicle type. A vehicle type is a definition of a
// vehicle. It contains the process duration and travel duration expressions
// that are used to calculate the travel and process duration of a stop
// assignment to a vehicle of this type.
type ModelVehicleType interface {
	// Data returns the arbitrary data associated with the vehicle type. Can be
	// set using the VehicleTypeData VehicleTypeOption in the factory method
	// Model.NewVehicleType.
	Data() interface{}

	// Index returns the index of the vehicle type.
	Index() int

	// Model returns the model of the vehicle type.
	Model() Model

	// ProcessDurationExpression returns the process duration expression of the
	// vehicle type. Is set in the factory method of the vehicle type
	// Model.NewVehicleType.
	ProcessDurationExpression() ModelExpression

	// TravelDurationExpression returns the travel duration expression of the
	// vehicle type. Is set in the factory method of the vehicle type
	// Model.NewVehicleType.
	TravelDurationExpression() ModelExpression
}

// ModelVehicleTypes is a slice of vehicle types.
type ModelVehicleTypes []ModelVehicleType

// VehicleTypeOption is an option for a vehicle type. Can be used in the
// factory method of a vehicle type Model.NewVehicleType.
type VehicleTypeOption func(ModelVehicleType) error

// VehicleTypeData is an option for a vehicle type. Can be used in the factory
// method of a vehicle type Model.NewVehicleType. The data is arbitrary data
// associated with the vehicle type.
func VehicleTypeData(
	data interface{},
) VehicleTypeOption {
	connect.Connect(con, &vehicleTypeDataVehicleTypeOption)
	return vehicleTypeDataVehicleTypeOption(data)
}
