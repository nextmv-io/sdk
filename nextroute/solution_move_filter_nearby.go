package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"time"
)

// SolutionMoveFilterNearby is a SolutionMoveFilter that filters out moves
// that are too far away from the locations in the move compared to the
// locations the move wants to plan then next to.
//
// If a move is adding a unit A - B in an existing vehicle F - Y - L like
// this F - A - B - Y - L, then the move is filtered out if:
//   - F - A travel duration > maxFirstLegDuration
//   - A - B travel duration > maxInterLegDuration
//   - B - Y travel duration > maxDuration
type SolutionMoveFilterNearby interface {
	// FilterMove returns true if the move should be filtered out.
	FilterMove(move Move) bool

	// MaxDuration returns the maximum duration a leg connecting two locations
	// of which one is in the move and the other is planned but not the first
	// or last leg.
	MaxDuration() time.Duration

	// MaxFirstLegDuration returns the maximum duration of the first leg
	// connecting the first location in the move to the first location in the
	// vehicle.
	MaxFirstLegDuration() time.Duration

	// MaxInterLegDuration returns the maximum duration of the inter leg
	// connecting two locations in the move.
	MaxInterLegDuration() time.Duration

	// MaxLastLegDuration returns the maximum duration of the last leg
	// connecting the last location in the move to the last location in the
	// vehicle.
	MaxLastLegDuration() time.Duration
}

// NewSolutionMoveFilterNearby returns a new SolutionMoveFilterNearby.

func NewSolutionMoveFilterNearby(
	maxDuration time.Duration,
	maxFirstLegDuration time.Duration,
	maxInterLegDuration time.Duration,
	maxLastLegDuration time.Duration,
) SolutionMoveFilterNearby {
	connect.Connect(con, &newSolutionMoveFilterNearby)
	return newSolutionMoveFilterNearby(
		maxDuration,
		maxFirstLegDuration,
		maxInterLegDuration,
		maxLastLegDuration,
	)
}
