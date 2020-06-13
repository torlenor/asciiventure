package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type keyboardEvent struct {
}

func (g *Game) handleSDLEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			g.commandManager.DispatchCommand(event)
		case *sdl.MouseMotionEvent:
			g.updateMouseTile(int(t.X), int(t.Y))
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				switch t.Button {
				case sdl.BUTTON_LEFT:
					g.setTargetPosition(g.mouseTileX, g.mouseTileY)
				case sdl.BUTTON_RIGHT:
				}
			}
		}
	}
}
