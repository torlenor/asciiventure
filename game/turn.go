package game

import (
	"fmt"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/pathfinding"
)

func (g *Game) blocked(x, y int) (*entity.Entity, bool) {
	for _, e := range g.entities {
		if e.Blocks == true && e.Position.X == x && e.Position.Y == y && g.player.FoV.Seen(components.Position{X: x, Y: y}) {
			return e, true
		}
	}
	return nil, false
}

// killEntity declares the entity dead.
func (g *Game) killEntity(e *entity.Entity) {
	e.Blocks = false
	e.Dead = true
	g.logWindow.AddRow(fmt.Sprintf("%s is dead.", e.Name))
}

func (g *Game) combat(e *entity.Entity, target *entity.Entity) {
	results := e.Attack(target)
	for _, result := range results {
		if result.Type == entity.CombatResultTakeDamage {
			target.Combat.CurrentHP -= result.IntegerValue
			g.logWindow.AddRow(fmt.Sprintf("%s scratches %s for %d hit points. %d/%d HP left.", e.Name, target.Name, result.IntegerValue, target.Combat.CurrentHP, target.Combat.HP))
			if target.Combat.CurrentHP <= 0 {
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
		var newPosition components.Position
		if e == g.player {
			if len(g.movementPath) > 0 {
				newPosition = g.movementPath[0]
			} else {
				newPosition = g.player.TargetPosition
			}
		} else {
			if e.AI != nil {
				var path []components.Position
				if g.currentGameMap.Distance(g.player.Position, e.InitialPosition) <= float64(e.AI.AttackRange) && g.currentGameMap.Distance(e.Position, e.InitialPosition) <= float64(e.AI.AttackRangeUntil) {
					path = pathfinding.DetermineAstarPath(g.currentGameMap, g, e.Position, g.player.Position)
				} else {
					path = pathfinding.DetermineAstarPath(g.currentGameMap, g, e.Position, e.InitialPosition)
				}
				if len(path) > 0 {
					newPosition = path[0]
				}
			}
		}
		if newPosition.Equal(e.Position) {
			continue
		}
		roomEmpty := g.currentGameMap.Empty(newPosition.X, newPosition.Y)
		blockingE, blocked := g.blocked(newPosition.X, newPosition.Y)
		if roomEmpty && !blocked {
			e.MoveTo(newPosition)
			if e == g.player {
				if len(g.movementPath) > 0 {
					g.movementPath = g.movementPath[1:]
				}
			}
		} else if blocked {
			if e.Combat != nil && blockingE.Combat != nil && !(e != g.player && blockingE != g.player) {
				g.combat(e, blockingE)
				e.TargetPosition = e.Position
				g.movementPath = []components.Position{}
			}
		} else if !roomEmpty {
			e.TargetPosition = e.Position
		}
	}
	if state == playersTurn {
		g.focusPlayer()
	}
}
