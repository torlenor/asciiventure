package pathfinding

import "github.com/torlenor/asciiventure/utils"

// Graph represents the graph on which the pathfinding shall be performed.
type Graph interface {
	Opaque(p utils.Vec2) bool
	InDimensions(p utils.Vec2) bool
	Neighbors(p utils.Vec2) []utils.Vec2
	Distance(a utils.Vec2, b utils.Vec2) float64
}

// Obstacles are positions which block the path finding algorithm.
// En example would be enemies or locked doors.
type Obstacles interface {
	Occupied(p utils.Vec2) bool
}
