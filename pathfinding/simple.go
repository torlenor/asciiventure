package pathfinding

import (
	"github.com/torlenor/asciiventure/utils"
)

// DetermineStraightLinePath returns the straightest possible path from start to goal.
func DetermineStraightLinePath(start utils.Vec2, goal utils.Vec2) []utils.Vec2 {
	current := start
	s := []utils.Vec2{}
	for goal.X != current.X || goal.Y != current.Y {
		if current.X < goal.X {
			current.X++
		} else if current.X > goal.X {
			current.X--
		}
		if current.Y < goal.Y {
			current.Y++
		} else if current.Y > goal.Y {
			current.Y--
		}
		s = append(s, current)
	}
	return s
}
