package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type keyboardEvent struct {
}

func (g *Game) handleSDLEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			g.commandManager.DispatchCommand(event)
		case *sdl.MouseMotionEvent:
			g.commandManager.DispatchMouseCommand(event)
		case *sdl.MouseButtonEvent:
			g.commandManager.DispatchMouseCommand(event)
		}
	}
}
