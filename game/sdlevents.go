package game

import (
	"strconv"

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
					g.renderer.OriginY += 2
					g.preRenderRoom()
					continue
				case sdl.K_DOWN:
					g.renderer.OriginY -= 2
					g.preRenderRoom()
					continue
				case sdl.K_LEFT:
					g.renderer.OriginX += 2
					g.preRenderRoom()
					continue
				case sdl.K_RIGHT:
					g.renderer.OriginX -= 2
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
						g.focusPlayer()
						continue
					case "-":
						g.renderScale -= 0.1
						g.preRenderRoom()
						g.focusPlayer()
						continue
					}
				}
			}
		case *sdl.MouseMotionEvent:
			g.updateMouseTile(t.X, t.Y)
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
