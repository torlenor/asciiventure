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
					if g.player.Mutations.Has(components.MutationInventory) {
						g.player.Inventory = append(g.player.Inventory, e)
						g.entities[i] = nil
						g.logWindow.AddRow(fmt.Sprintf("%s added to inventory.", e.Name))
						g.updateInventory()
					} else {
						g.logWindow.AddRow(fmt.Sprintf("You are a cat, you cannot pick up things (or can you?)."))
					}
				} else if e.Item.Mutagen {
					g.player.Mutations = append(g.player.Mutations, e.Mutations...)
					for _, m := range e.Mutations {
						g.logWindow.AddRow(fmt.Sprintf("Mutation %s consumed.", m))
					}
					g.entities[i] = nil
					g.updateMutationsPane()
				}
				break
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
