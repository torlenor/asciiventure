package game

import (
	"fmt"

	"github.com/torlenor/asciiventure/components"
)

func (g *Game) performAction() {
	for i, e := range g.entities {
		if e.Position.Equal(g.player.Position) && e != g.player {
			if e.Item != nil {
				if e.Item.CanPickup {
					if g.player.Mutations.Has(components.MutationEffectInventory) {
						g.player.Inventory = append(g.player.Inventory, e)
						g.entities[i] = nil
						g.logWindow.AddRow(fmt.Sprintf("%s added to inventory.", e.Name))
					} else {
						g.logWindow.AddRow(fmt.Sprintf("You are a cat, you cannot pick up things (or can you?)."))
					}
				}
			}
			if e.Mutation != nil {
				if !g.player.Mutations.Has(e.Mutation.Effect) {
					g.player.Mutations = append(g.player.Mutations, *e.Mutation)
					g.logWindow.AddRow(fmt.Sprintf("Mutation %s consumed.", g.player.Mutations[len(g.player.Mutations)-1]))
					g.entities[i] = nil
				} else {
					g.logWindow.AddRow(fmt.Sprintf("Player already has Mutation %s", e.Mutation))
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
