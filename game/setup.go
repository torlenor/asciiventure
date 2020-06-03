package game

import (
	"log"

	"github.com/torlenor/asciiventure/renderers"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func (g *Game) setupWindow() {
	err := ttf.Init()
	if err != nil {
		log.Fatalf("Failed to initialize ttf: %s", err)
	}

	if g.defaultFont, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		log.Fatalf("Failed to load font '%s': %s", fontPath, err)
	}

	err = sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		log.Fatalf("Failed to initialize sdl: %s", err)
	}

	g.window, err = sdl.CreateWindow(windowName, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s", err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
}

func (g *Game) setupRenderer() {
	// renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_ACCELERATED)
	renderer, err := sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s", err)
	}
	g.renderer = renderers.NewRenderer(renderer)
	g.renderer.GlyphWidth = latticeDX
	g.renderer.GlyphHeight = latticeDY

	g.renderer.OriginY = characterWindowRect.H/latticeDY + 1
}

func (g *Game) setupGame() {
	g.createGlyphTexture()
	g.createPlayer()
	g.loadGameMapsFromDirectory("./assets/rooms")
	g.selectGameMap(1)

	g.updateCharacterWindow()
	g.logWindow.SetText([]string{"Welcome to <Epic Name Here>.", "A small cat takes a stroll and ends up in an epic adventure."})
}
