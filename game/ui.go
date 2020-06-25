package game

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/utils"
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
		if e.Position.Current.Equal(utils.Vec2{X: targetX, Y: targetY}) && e != g.player {
			if e.IsDead != nil {
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
	if g.currentGameMap.IsPortal(utils.Vec2{X: targetX, Y: targetY}) {
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
