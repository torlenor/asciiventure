package game

import (
	"fmt"

	"github.com/torlenor/asciiventure/components"
)

func (g *Game) updateUI() {
	g.updateCharacterWindow()
	g.updateInventoryPane()
	g.updateMutationsPane()
}

func (g *Game) updateStatusBar() {
	g.statusBar.Clear()
	for _, e := range g.entities {
		if e.Position.Equal(components.Position{X: g.mouseTileX, Y: g.mouseTileY}) && e != g.player {
			if e.Dead {
				g.statusBar.AddRow(e.Name + "(Dead)")
			} else {
				if e.Item != nil {
					g.statusBar.AddRow(e.Name + ": Pick up item with 'g'")
				} else if e.Mutation != nil {
					g.statusBar.AddRow(e.Mutation.String() + ": " + e.Mutation.GetDescription())
				} else {
					g.statusBar.AddRow(e.Name)
				}
			}
			return
		}
	}
	g.statusBar.Clear()
}

func (g *Game) updateMutationsPane() {
	g.mutations.Clear()
	if len(g.player.Mutations) == 0 {
		g.mutations.AddRow("No mutations")
		return
	}
	g.mutations.AddRow("Mutations:")
	g.mutations.AddRow("----------------")
	for _, m := range g.player.Mutations {
		g.mutations.AddRow(fmt.Sprintf("%s", m))
	}
}

func (g *Game) updateInventoryPane() {
	g.inventory.UpdateInventory(g.player.Inventory)
}
