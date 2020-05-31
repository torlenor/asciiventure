package game

import (
	"fmt"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

func (g *Game) blocked(x, y int32) (*entity.Entity, bool) {
	for _, e := range g.entities {
		if e.Blocks == true && e.Position.X == x && e.Position.Y == y {
			return e, true
		}
	}
	return nil, false
}

func (g *Game) updatePositions(state gameState) {
	for _, e := range g.entities {
		if (state == playersTurn && e != g.player) || (state == enemyTurn && e == g.player) || e.IsDead {
			continue
		}
		var path []components.Position
		if e == g.player {
			path = determineLatticePath(e.Position, e.TargetPosition)
		} else {
			path = determineLatticePath(e.Position, g.player.Position)
		}
		if len(path) == 0 {
			continue
		}
		if len(path) > 0 {
			newP := path[0]
			roomEmpty := g.currentRoom.Empty(newP.X, newP.Y)
			blockingE, blocked := g.blocked(newP.X, newP.Y)
			if roomEmpty && !blocked {
				e.MoveTo(newP)
			} else if blocked {
				if e.Combat != nil && blockingE.Combat != nil && !(e != g.player && blockingE != g.player) {
					results := e.Combat.Attack(blockingE.Combat)
					for _, result := range results {
						if result.Type == components.TakeDamage {
							blockingE.Combat.HP -= result.IntegerValue
							g.logWindow.AddRow(fmt.Sprintf("%s scratches %s for %d hit points. %d/%d HP left.", e.Name, blockingE.Name, result.IntegerValue, blockingE.Combat.HP, blockingE.Combat.MaxHP))
							if blockingE.Combat.HP <= 0 {
								blockingE.Glyph, _ = g.glyphTexture.Get("%")
								blockingE.Glyph.Color = components.ColorRGB{R: 200, G: 200, B: 200}
								blockingE.Blocks = false
								blockingE.IsDead = true
								g.logWindow.AddRow(fmt.Sprintf("%s is dead.", blockingE.Name))
								if blockingE == g.player {
									g.gameState = gameOver
								}
							}
						}
					}
				}
				e.TargetPosition = e.Position
			} else if !roomEmpty {
				e.TargetPosition = e.Position
			}
		}
	}
	if state == playersTurn {
		g.markedPath = g.determineLatticePathPlayerMouse()
	}
}
