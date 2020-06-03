package pathfinding

import (
	"github.com/torlenor/asciiventure/components"
)

// Graph represents the graph on which the pathfinding shall be performed.
type Graph interface {
	Opaque(p components.Position) bool
	InDimensions(p components.Position) bool
	Neighbors(p components.Position) []components.Position
	Distance(a components.Position, b components.Position) float64
}

// Obstacles are positions which block the path finding algorithm.
// En example would be enemies or locked doors.
type Obstacles interface {
	Occupied(p components.Position) bool
}
