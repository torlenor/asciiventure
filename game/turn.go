package game

import (
	"fmt"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/pathfinding"
	"github.com/torlenor/asciiventure/utils"
)

func (g *Game) blocked(p utils.Vec2) (*entity.Entity, bool) {
	for _, e := range g.entities {
		if e.IsBlocking != nil && e.Position.Current.Equal(p) && g.player.FoV.Seen(p) {
			return e, true
		}
	}
	return nil, false
}

// killEntity declares the entity dead.
func (g *Game) killEntity(e *entity.Entity) {
	e.IsBlocking = nil
	e.IsDead = &components.IsDead{}
	g.ui.AddLogEntry(fmt.Sprintf("%s is dead.", e.Name))
}

func (g *Game) combat(e *entity.Entity, target *entity.Entity) {
	results := e.Attack(target)
	for _, result := range results {
		if result.Type == entity.CombatResultTakeDamage {
			target.Health.CurrentHP -= result.IntegerValue
			g.ui.AddLogEntry(fmt.Sprintf("%s scratches %s for %d hit points. %d/%d HP left.", e.Name, target.Name, result.IntegerValue, target.Health.CurrentHP, target.Health.HP))
			if target.Health.CurrentHP <= 0 {
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
		if (state == playersTurn && e != g.player) || (state == enemyTurn && e == g.player) || e.IsDead != nil || e.Position == nil {
			continue
		}
		var newPosition utils.Vec2
		if e == g.player {
			if len(g.movementPath) > 0 {
				newPosition = g.movementPath[0]
			} else {
				newPosition = g.player.TargetPosition
			}
		} else {
			if e.AI != nil {
				var path []utils.Vec2
				if g.currentGameMap.Distance(g.player.Position.Current, e.Position.Initial) <= float64(e.AI.AttackRange) && g.currentGameMap.Distance(e.Position.Current, e.Position.Initial) <= float64(e.AI.AttackRangeUntil) {
					path = pathfinding.DetermineAstarPath(g.currentGameMap, g, e.Position.Current, g.player.Position.Current)
				} else {
					path = pathfinding.DetermineAstarPath(g.currentGameMap, g, e.Position.Current, e.Position.Initial)
				}
				if len(path) > 0 {
					newPosition = path[0]
				}
			}
		}
		if newPosition.Equal(e.Position.Current) {
			continue
		}
		roomEmpty := g.currentGameMap.Empty(newPosition)
		blockingE, blocked := g.blocked(newPosition)
		if roomEmpty && !blocked {
			e.MoveTo(newPosition)
			if e == g.player {
				if len(g.movementPath) > 0 {
					g.movementPath = g.movementPath[1:]
				} else {
					g.movementPath = []utils.Vec2{}
					e.TargetPosition = e.Position.Current
				}
			}
		} else if blocked {
			if e.Combat != nil && blockingE.Combat != nil && !(e != g.player && blockingE != g.player) {
				g.combat(e, blockingE)
				g.movementPath = []utils.Vec2{}
				e.TargetPosition = e.Position.Current
			}
		} else if !roomEmpty {
			g.movementPath = []utils.Vec2{}
			e.TargetPosition = e.Position.Current
		}
	}
}
