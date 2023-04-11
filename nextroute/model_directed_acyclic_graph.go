package nextroute

import "github.com/nextmv-io/sdk/connect"

// Arc is a directed connection between two nodes.
type Arc interface {
	// Origin returns the origin node of the arc.
	Origin() ModelStop
	// Destination returns the destination node of the arc.
	Destination() ModelStop
}

// Arcs is a collection of arcs.
type Arcs []Arc

// DirectedAcyclicGraph is a set of nodes connected by arcs that does not
// contain cycles.
type DirectedAcyclicGraph interface {
	// NewArc creates a new arc in the graph.
	NewArc(origin, destination ModelStop) (Arc, error)
	// Arcs returns all arcs in the graph.
	Arcs() Arcs
}

// NewDirectedAcyclicGraph creates a new [DirectedAcyclicGraph].
func NewDirectedAcyclicGraph() DirectedAcyclicGraph {
	connect.Connect(con, &newDirectedAcyclicGraph)
	return newDirectedAcyclicGraph()
}
