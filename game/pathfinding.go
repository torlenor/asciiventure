package game

import "github.com/torlenor/asciiventure/components"

func determineLatticePath(origin components.Position, target components.Position) []components.Position {
	// TODO: determineLatticePath should take into account occupied tiles and be able to suggest a way around it.
	// https://en.wikipedia.org/wiki/A*_search_algorithm
	current := origin
	s := []components.Position{}
	for target.X != current.X || target.Y != current.Y {
		if current.X < target.X {
			current.X++
		} else if current.X > target.X {
			current.X--
		}
		if current.Y < target.Y {
			current.Y++
		} else if current.Y > target.Y {
			current.Y--
		}
		lp := components.Position{X: current.X, Y: current.Y}
		s = append(s, lp)
	}
	return s
}
