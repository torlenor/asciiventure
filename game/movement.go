package game

import "github.com/torlenor/asciiventure/entities"

func (g *Game) updatePositions() {
	var toRemove []entities.Entity
	for e, p := range g.position {
		if t, ok := g.targetPosition[e]; ok {
			if t.X == p.X && t.Y == p.Y {
				continue
			}
			path := determineLatticePath(p, t)
			if len(path) > 0 {
				newP := path[0]
				occupant, occupied := g.occupied(newP.X, newP.Y)
				roomEmpty := g.currentRoom.Empty(newP.X, newP.Y)
				if !g.collision[e].V || (roomEmpty && !occupied) {
					g.position[e] = newP
				} else if g.collision[e].DestroyOnCollision && (!roomEmpty || occupied) {
					toRemove = append(toRemove, e)
					if occupied {
						if g.collision[occupant].DestroyOnCollision {
							toRemove = append(toRemove, occupant)
						}
					}
				} else {
					g.targetPosition[e] = g.position[e]
				}
			}
		}
	}
	for _, e := range toRemove {
		g.removeEntity(e)
	}
	g.markedPath = g.determineLatticePathPlayerMouse()
}
