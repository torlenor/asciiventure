package game

import "github.com/torlenor/asciiventure/components"

func (g *Game) determinePathPlayerMouse() []components.Position {
	return g.determineLatticePathAstar(components.Position{X: g.player.Position.X, Y: g.player.Position.Y}, components.Position{X: g.mouseTileX, Y: g.mouseTileY})
}

func (g *Game) updateMouseTile(x, y int32) {
	g.mouseTileX = int32((float32(x)+0.5)/latticeDX/g.renderScale) - g.renderer.OriginX
	g.mouseTileY = int32((float32(max(y, screenHeight/6))+0.5)/latticeDY/g.renderScale) - g.renderer.OriginY
}

func (g *Game) renderMouseTile() {
	if !g.player.Position.Equal(g.player.TargetPosition) {
		path := g.movementPath
		for _, p := range path {
			notEmpty := !g.currentRoom.Empty(p.X, p.Y) && g.currentRoom.Visible(p.X, p.Y)
			_, blocked := g.blocked(p.X, p.Y)
			color := components.ColorRGBA{R: 100, G: 100, B: 255, A: 64}
			if notEmpty || (blocked && g.currentRoom.Visible(p.X, p.Y) && g.currentRoom.Seen(p.X, p.Y)) {
				color = components.ColorRGBA{R: 255, G: 80, B: 80, A: 100}
			}
			g.renderer.FillCharCoordinate(p.X, p.Y, color)
			if notEmpty || (blocked && g.currentRoom.Visible(p.X, p.Y) && g.currentRoom.Seen(p.X, p.Y)) {
				break
			}
		}
	}

	for _, p := range g.markedPath {
		notEmpty := !g.currentRoom.Empty(p.X, p.Y) && g.currentRoom.Visible(p.X, p.Y)
		_, blocked := g.blocked(p.X, p.Y)
		color := components.ColorRGBA{R: 100, G: 255, B: 100, A: 64}
		if notEmpty || (blocked && g.currentRoom.Visible(p.X, p.Y) && g.currentRoom.Seen(p.X, p.Y)) {
			color = components.ColorRGBA{R: 255, G: 100, B: 100, A: 64}
		}
		g.renderer.FillCharCoordinate(p.X, p.Y, color)
		if notEmpty || (blocked && g.currentRoom.Visible(p.X, p.Y) && g.currentRoom.Seen(p.X, p.Y)) {
			break
		}
	}
	if len(g.markedPath) > 0 {
		p := g.markedPath[len(g.markedPath)-1]
		color := components.ColorRGBA{R: 128, G: 128, B: 128, A: 120}
		g.renderer.FillCharCoordinate(p.X, p.Y, color)
	}
}
