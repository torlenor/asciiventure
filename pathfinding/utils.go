package pathfinding

import "github.com/torlenor/asciiventure/components"

func contains(list []components.Position, p components.Position) bool {
	for _, c := range list {
		if c.X == p.X && c.Y == p.Y {
			return true
		}
	}
	return false
}
