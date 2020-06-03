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

			// shift := false
			// ctrl := false
			alt := false

			switch t.Keysym.Mod {
			case sdl.KMOD_LSHIFT:
				fallthrough
			case sdl.KMOD_RSHIFT:
				fallthrough
			case sdl.KMOD_SHIFT:
				// shift = true
			case sdl.KMOD_LCTRL:
				fallthrough
			case sdl.KMOD_RCTRL:
				fallthrough
			case sdl.KMOD_CTRL:
				// ctrl = true
			case sdl.KMOD_RALT:
				fallthrough
			case sdl.KMOD_LALT:
				fallthrough
			case sdl.KMOD_ALT:
				alt = true
			}

			if t.State == sdl.PRESSED {
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					g.quit = true
					continue
				case sdl.K_F5:
					g.loadGameMapsFromDirectory("./assets/rooms")
					continue
				case sdl.K_UP:
					if alt {
						g.renderer.OriginY += 2
						g.preRenderGameMap()
					} else {
						g.movementPath = []components.Position{}
						g.player.TargetPosition.Y = g.player.Position.Y - 1
						g.nextStep = true
					}
					continue
				case sdl.K_DOWN:
					if alt {
						g.renderer.OriginY -= 2
						g.preRenderGameMap()
					} else {
						g.movementPath = []components.Position{}
						g.player.TargetPosition.Y = g.player.Position.Y + 1
						g.nextStep = true
					}
					continue
				case sdl.K_LEFT:
					if alt {
						g.renderer.OriginX += 2
						g.preRenderGameMap()
					} else {
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X - 1
						g.nextStep = true
					}
					continue
				case sdl.K_RIGHT:
					if alt {
						g.renderer.OriginX -= 2
						g.preRenderGameMap()
					} else {
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X + 1
						g.nextStep = true
					}
					continue
				}
			}

			keyCode := t.Keysym.Sym

			if keyCode < 10000 {
				if t.State == sdl.PRESSED {
					if keyCode == sdl.K_SPACE {
						g.nextStep = true
						continue
					}

					k := string(keyCode)
					if keyCode >= '0' && keyCode <= '9' {
						r, _ := strconv.Atoi(k)
						g.selectGameMap(r)
						continue
					}

					switch k {
					case "+":
						g.renderScale += 0.1
						g.preRenderGameMap()
						g.focusPlayer()
						continue
					case "-":
						g.renderScale -= 0.1
						g.preRenderGameMap()
						g.focusPlayer()
						continue
					// y	k	u
					// h		l
					// b	j	n
					case "y":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X - 1
						g.player.TargetPosition.Y = g.player.Position.Y - 1
						g.nextStep = true
						continue
					case "k":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.Y = g.player.Position.Y - 1
						g.nextStep = true
						continue
					case "u":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X + 1
						g.player.TargetPosition.Y = g.player.Position.Y - 1
						g.nextStep = true
						continue
					case "h":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X - 1
						g.nextStep = true
						continue
					case "l":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X + 1
						g.nextStep = true
						continue
					case "b":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X - 1
						g.player.TargetPosition.Y = g.player.Position.Y + 1
						g.nextStep = true
						continue
					case "j":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.Y = g.player.Position.Y + 1
						g.nextStep = true
						continue
					case "n":
						g.movementPath = []components.Position{}
						g.player.TargetPosition.X = g.player.Position.X + 1
						g.player.TargetPosition.Y = g.player.Position.Y + 1
						g.nextStep = true
						continue
					}
				}
			}
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
