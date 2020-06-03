package game

import (
	"fmt"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

func (g *Game) blocked(x, y int32) (*entity.Entity, bool) {
	for _, e := range g.entities {
		if e.Blocks == true && e.Position.X == x && e.Position.Y == y && g.player.FoV.Seen(components.Position{X: x, Y: y}) {
			return e, true
		}
	}
	return nil, false
}

// killEntity declares the entity dead.
func (g *Game) killEntity(e *entity.Entity) {
	e.Glyph, _ = g.glyphTexture.Get("%")
	e.Glyph.Color = components.ColorRGB{R: 150, G: 150, B: 150}
	e.Blocks = false
	e.Dead = true
	g.logWindow.AddRow(fmt.Sprintf("%s is dead.", e.Name))
}

func (g *Game) combat(e *entity.Entity, target *entity.Entity) {
	results := e.Attack(target)
	for _, result := range results {
		if result.Type == components.TakeDamage {
			target.Combat.HP -= result.IntegerValue
			g.logWindow.AddRow(fmt.Sprintf("%s scratches %s for %d hit points. %d/%d HP left.", e.Name, target.Name, result.IntegerValue, target.Combat.HP, target.Combat.MaxHP))
			if target.Combat.HP <= 0 {
				g.killEntity(target)
				if target == g.player {
					g.gameState = gameOver
				}
			}
		}
	}
}

func (g *Game) updatePositions(state gameState) {
	for _, e := range g.entities {
		if (state == playersTurn && e != g.player) || (state == enemyTurn && e == g.player) || e.Dead {
			continue
		}
		var path []components.Position
		if e == g.player {
			path = g.movementPath
		} else {
			path = determineStraightLinePath(e.Position, g.player.Position)
		}
		if len(path) == 0 {
			continue
		}
		newP := path[0]
		roomEmpty := g.currentRoom.Empty(newP.X, newP.Y)
		blockingE, blocked := g.blocked(newP.X, newP.Y)
		if roomEmpty && !blocked {
			e.MoveTo(newP)
			if e == g.player {
				g.movementPath = g.movementPath[1:]
			}
		} else if blocked {
			if e.Combat != nil && blockingE.Combat != nil && !(e != g.player && blockingE != g.player) {
				g.combat(e, blockingE)
			}
			e.TargetPosition = e.Position
		} else if !roomEmpty {
			e.TargetPosition = e.Position
		}
	}
	if state == playersTurn {
		g.focusPlayer()
		g.preRenderRoom()
	}
}
