package game

import (
	"github.com/torlenor/asciiventure/components"
)

func (g *Game) cleanupEntities() {
	n := 0
	for _, e := range g.entities {
		if e != nil {
			g.entities[n] = e
			n++
		}
	}
	g.entities = g.entities[:n]
}

func (g *Game) performPlayerAction(at components.ActionType, intValue int) {
	g.player.Actor = &components.Actor{NextAction: at}

	switch at {
	case components.ActionTypeInteract:
		if g.currentGameMap.IsPortal(g.player.Position.Current) {
			g.selectGameMap(g.currentGamMapID + 1)
		}
	case components.ActionTypeDropItem:
		// TODO: Implement DropItem
	case components.ActionTypeUseItem:
		g.player.Actor.IntValue = intValue
	}
}
