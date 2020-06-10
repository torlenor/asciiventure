package game

import (
	"github.com/torlenor/asciiventure/components"
)

func (g *Game) updateUI() {
	g.updateCharacterWindow()
	g.updateInventoryPane()
	g.updateMutationsPane()
}

func (g *Game) updateStatusBar() {
	for _, e := range g.entities {
		if e.Position.Equal(components.Position{X: g.mouseTileX, Y: g.mouseTileY}) && e != g.player {
			if e.Dead {
				g.ui.SetStatusBarText(e.Name + "(Dead)")
			} else {
				if e.Item != nil {
					g.ui.SetStatusBarText(e.Name + ": Pick up item with 'g'")
				} else if e.Mutation != nil {
					g.ui.SetStatusBarText(e.Mutation.String() + ": " + e.Mutation.GetDescription())
				} else {
					g.ui.SetStatusBarText(e.Name)
				}
			}
			return
		}
	}
	if g.currentGameMap.IsPortal(components.Position{X: g.mouseTileX, Y: g.mouseTileY}) {
		g.ui.SetStatusBarText("Stairs to next map. Press 'g' to use them.")
	} else {
		g.ui.SetStatusBarText("")
	}
}

func (g *Game) updateMutationsPane() {
	g.ui.SetInventoryPaneEnabled(g.player.Mutations.Has(components.MutationEffectInventory))
	g.ui.UpdateMutationsPane(g.player.Mutations)
}

func (g *Game) updateInventoryPane() {
	g.ui.UpdateInventoryPane(g.player.Inventory)
}
