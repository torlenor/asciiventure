package pathfinding

import "github.com/torlenor/asciiventure/components"

func DetermineStraightLinePath(start components.Position, goal components.Position) []components.Position {
	current := start
	s := []components.Position{}
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
		lp := components.Position{X: current.X, Y: current.Y}
		s = append(s, lp)
	}
	return s
}
