package game

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/torlenor/asciiventure/console"
	"github.com/torlenor/asciiventure/gamemap"
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

	if g.fullscreen {
		g.window, err = sdl.CreateWindow(windowName, 0,
			0, 0, 0, sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN_DESKTOP)
	} else {
		g.window, err = sdl.CreateWindow(windowName, sdl.WINDOWPOS_CENTERED,
			sdl.WINDOWPOS_CENTERED, int32(g.screenWidth), int32(g.screenHeight), sdl.WINDOW_SHOWN)
	}
	if err != nil {
		log.Fatalf("Failed to create window: %s", err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	if runtime.GOOS == "windows" {
		sdl.SetHint(sdl.HINT_RENDER_DRIVER, "opengl")
	}
}

func (g *Game) setupRenderer() {
	renderer, err := sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s", err)
	}
	g.renderer = renderers.NewRenderer(renderer)
	// g.renderer.GlyphWidth = latticeDX
	// g.renderer.GlyphHeight = latticeDY
}

func (g *Game) setupConsoles() {
	tileset, err := console.NewFontTileset(g.renderer, "./assets/textures/consolas12x12_gs_tc.png")
	if err != nil {
		log.Fatalf("%s", err)
	}
	availableWidth := int32(g.screenWidth - g.screenWidth/5)
	availableHeight := int32(g.screenHeight - g.screenHeight/6 - 2*16)
	g.consoleMap = console.NewMatrixConsole(g.renderer, tileset, availableWidth, availableHeight, availableWidth/tileset.GetCharWidth(), availableHeight/tileset.GetCharHeight())
	g.consoleMap.SetOffset(0, int32(g.screenHeight/6))
}

func (g *Game) setupGame() {
	rand.Seed(time.Now().UnixNano())
	g.createGlyphTexture()
	g.createPlayer()
	g.loadedGameMaps = []*gamemap.GameMap{}
	for i := 0; i < 3; i++ {
		randomMap := gamemap.NewRandomMap(10, 6, 20, 100, 60, g.glyphTexture)
		g.loadedGameMaps = append(g.loadedGameMaps, &randomMap)
	}
	g.loadGameMapsFromDirectory("./assets/rooms")
	g.selectGameMap(1)

	g.updateUI()

	g.consoleMap.Clear()
	g.currentGameMap.Render(g.consoleMap, g.player.FoV, g.player, g.entities, int32(g.renderer.OriginX), int32(g.renderer.OriginY))

	g.ui.AddLogEntry("Welcome to Lala's Quest.")
	g.ui.AddLogEntry("You are a young cat out hunting for mice.")
}
