package game

import (
	"github.com/torlenor/asciiventure/entity"
)

func (g *Game) performAction() {
	for i, e := range g.entities {
		if e.Position.Equal(g.player.Position) && e != g.player {
			result := g.player.PickUp(e)
			for _, r := range result {
				switch r.Type {
				case entity.ActionResultItemPickedUp:
					// Remove the entity from global list as it is now in the inventory
					g.entities[i] = nil
				case entity.ActionResultMessage:
					g.logWindow.AddRow(r.StringValue)
				}
			}

			result = g.player.ConsumeMutation(e)
			for _, r := range result {
				switch r.Type {
				case entity.ActionResultMutationConsumed:
					// Remove the entity from global list as it is now in the inventory
					g.entities[i] = nil
				case entity.ActionResultMessage:
					g.logWindow.AddRow(r.StringValue)
				}
			}
		}
	}
	n := 0
	for _, e := range g.entities {
		if e != nil {
			g.entities[n] = e
			n++
		}
	}
	g.entities = g.entities[:n]
}
