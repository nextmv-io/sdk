package here

import h "github.com/nextmv-io/sdk/measure/here"

// BoundingBox represents a region using four cooordinates corresponding
// to the furthest points in each of the cardinal directions within that region.
type BoundingBox = h.BoundingBox

// TransportMode represents the type of vehicle that will be used for
// the calculated routes.
type TransportMode = h.TransportMode

// Feature represents a geographical feature.
type Feature = h.Feature

// Truck captures truck-specific routing parameters.
type Truck = h.Truck

// WeightPerAxleGroup captures the weights of different axle groups.
type WeightPerAxleGroup = h.WeightPerAxleGroup

// TunnelCategory is a tunnel category used to restrict the transport of
// certain goods.
type TunnelCategory = h.TunnelCategory

// TruckType specifies the type of truck.
type TruckType = h.TruckType

// HazardousGood indicates a hazardous good that trucks can transport.
type HazardousGood = h.HazardousGood

// Scooter captures routing parameters that can be set on scooters.
type Scooter = h.Scooter

// Taxi captures routing parameters that can be set on taxis.
type Taxi = h.Taxi
