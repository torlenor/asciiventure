package pathfinding

import "github.com/torlenor/asciiventure/utils"

func contains(list []utils.Vec2, p utils.Vec2) bool {
	for _, c := range list {
		if c.X == p.X && c.Y == p.Y {
			return true
		}
	}
	return false
}
