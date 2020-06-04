package game

import (
	"fmt"

	"github.com/torlenor/asciiventure/components"
)

func (g *Game) performAction() {
	for _, e := range g.entities {
		if e.Position.Equal(g.player.Position) && e != g.player {
			if e.Item != nil {
				if e.Item.CanPickup {
					if g.player.Mutations.Has(components.MutationInventory) {
						g.logWindow.AddRow(fmt.Sprintf("%s added to inventory.", e.Name))
					} else {
						g.logWindow.AddRow(fmt.Sprintf("You are a cat, you cannot pick up things (or can you?)."))
					}
				} else if e.Item.Mutagen {
					g.player.Mutations = append(g.player.Mutations, e.Mutations...)
					for _, m := range e.Mutations {
						g.logWindow.AddRow(fmt.Sprintf("Mutation %s consumed.", m))
					}
					g.updateMutationsPane()
				}
				return
			}
		}
	}
}
