package game

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/pathfinding"
	"github.com/torlenor/asciiventure/utils"
)

func (g *Game) determinePathPlayerMouse() []components.Position {
	targetX, targetY := g.currentGameMap.GetPositionFromRenderCoordinates(g.mouseTileX, g.mouseTileY)
	return pathfinding.DetermineAstarPath(g.currentGameMap, g, components.Position{X: g.player.Position.X, Y: g.player.Position.Y}, components.Position{X: int(targetX), Y: int(targetY)})
}

func (g *Game) updateMouseTile(x, y int) {
	cx, cy := g.consoleMap.GetTileFromScreenCoordinates(int32(x), int32(y))

	g.mouseTileX = cx
	g.mouseTileY = cy
	g.updateStatusBar()
}

func (g *Game) renderMouseTile() {
	if !g.player.Position.Equal(g.player.TargetPosition) {
		path := g.movementPath
		for _, p := range path {
			notEmpty := !g.currentGameMap.Empty(p.X, p.Y) && g.player.FoV.Visible(p)
			_, blocked := g.blocked(p.X, p.Y)
			color := utils.ColorRGBA{R: 100, G: 100, B: 255, A: 64}
			if notEmpty || blocked {
				color = utils.ColorRGBA{R: 255, G: 80, B: 80, A: 100}
			}
			rx, ry := g.currentGameMap.GetRenderCoordinatesFromPosition(int32(p.X), int32(p.Y))
			g.consoleMap.SetBackgroundColor(rx, ry, color)
			if notEmpty || blocked {
				break
			}
		}
	}

	color := utils.ColorRGBA{R: 128, G: 128, B: 128, A: 120}
	g.consoleMap.SetBackgroundColor(g.mouseTileX, g.mouseTileY, color)
}
