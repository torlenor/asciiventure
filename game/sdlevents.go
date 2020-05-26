package game

import (
	"strconv"

	"github.com/torlenor/asciiventure/components"
	"github.com/veandco/go-sdl2/sdl"
)

func (g *Game) handleSDLEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			g.quit = true
		case *sdl.KeyboardEvent:
			if t.State == sdl.PRESSED {
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					g.quit = true
					continue
				case sdl.K_F5:
					g.loadRoomsFromDirectory("./assets/rooms")
					continue
				case sdl.K_UP:
					g.renderer.OriginY++
					g.preRenderRoom()
					continue
				case sdl.K_DOWN:
					g.renderer.OriginY--
					g.preRenderRoom()
					continue
				case sdl.K_LEFT:
					g.renderer.OriginX++
					g.preRenderRoom()
					continue
				case sdl.K_RIGHT:
					g.renderer.OriginX--
					g.preRenderRoom()
					continue
				}
			}

			keyCode := t.Keysym.Sym

			// upperCase := false
			// rightAlt := false

			switch t.Keysym.Mod {
			case sdl.KMOD_LSHIFT:
				// upperCase = true
			case sdl.KMOD_RALT:
				// rightAlt = true
			}

			if keyCode < 10000 {
				if t.State == sdl.PRESSED {
					if keyCode == sdl.K_SPACE {
						// TODO: Marking next time step should not be event based
						g.nextStep = true
						continue
					}

					k := string(keyCode)
					if keyCode >= '0' && keyCode <= '9' {
						r, _ := strconv.Atoi(k)
						g.selectRoom(r)
						continue
					}

					switch k {
					case "+":
						g.renderScale += 0.1
						g.preRenderRoom()
						continue
					case "-":
						g.renderScale -= 0.1
						g.preRenderRoom()
						continue
					case "e":
						g.createEnemy(components.Position{X: g.mouseTileX, Y: g.mouseTileY})
						continue
					}
				}
			}
		case *sdl.MouseMotionEvent:
			g.updateMouse(t.X, t.Y)
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				switch t.Button {
				case sdl.BUTTON_LEFT:
					g.setPlayerTargetPosition(g.mouseTileX, g.mouseTileY)
				case sdl.BUTTON_RIGHT:
				}
			}
		}
	}
}
