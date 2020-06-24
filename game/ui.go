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
	targetX, targetY := g.currentGameMap.GetPositionFromRenderCoordinates(g.mouseTileX, g.mouseTileY)
	for _, e := range g.entities {
		if e == nil || e.Position == nil {
			continue
		}
		if e.Position.Equal(components.Position{X: int(targetX), Y: int(targetY)}) && e != g.player {
			if e.Dead {
				g.ui.SetStatusBarText(e.Name + "(Dead)")
			} else {
				if e.Item != nil {
					g.ui.SetStatusBarText(e.Name + ": Pick up item with 'g'")
				} else if e.Mutagen != nil {
					g.ui.SetStatusBarText(e.Mutagen.String() + ": " + e.Mutagen.GetDescription())
				} else {
					g.ui.SetStatusBarText(e.Name)
				}
			}
			return
		}
	}
	if g.currentGameMap.IsPortal(components.Position{X: int(targetX), Y: int(targetY)}) {
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

func (g *Game) updateCharacterWindow() {
	g.ui.UpdateCharacterPane(g.time, g.player.Health.CurrentHP, g.player.Health.HP, g.player.Vision.Range+g.player.Mutations.GetData(components.MutationEffectIncreasedVision), g.player.Combat.Power, g.player.Combat.Defense)
}
