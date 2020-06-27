package game

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/console"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/gamemap"
	"github.com/torlenor/asciiventure/renderers"
	"github.com/torlenor/asciiventure/ui"
	"github.com/torlenor/asciiventure/utils"
)

const (
	windowName = "Asciiventure"

	fontPath = "./assets/fonts/RobotoMono-Regular.ttf"
	fontSize = 16

	latticeDX = 19
	latticeDY = 32
)

// Game is the main struct of the game
type Game struct {
	debug bool

	quit bool

	screenWidth  int
	screenHeight int
	fullscreen   bool

	window   *sdl.Window
	renderer *renderers.Renderer

	renderScale float32

	defaultFont  *ttf.Font
	glyphTexture *assets.GlyphTexture

	currentGameMap  *gamemap.GameMap
	currentGamMapID int
	loadedGameMaps  []*gamemap.GameMap

	mouseTileX int32
	mouseTileY int32

	movementPath []utils.Vec2

	player   *entity.Entity
	entities []*entity.Entity

	time uint

	nextStep  bool
	gameState gameState

	ui             *ui.UI
	commandManager *commandManager

	consoleMap      *console.MatrixConsole
	consoleMainMenu *console.MatrixConsole

	mainMenu *MainMenu

	gameInProgress bool
}

// Setup should be called first after creating an instance of Game.
func (g *Game) Setup(windowWidth, windowHeight int, fullscreen bool) {
	g.debug = true

	g.renderScale = 1.0

	g.screenWidth = windowWidth
	g.screenHeight = windowHeight
	g.fullscreen = fullscreen

	g.setupWindow()
	g.setupRenderer()

	if fullscreen {
		dm, err := sdl.GetCurrentDisplayMode(0)
		if err != nil {
			log.Fatalf("Cannot get current display mode: %s", err)
		}
		g.screenWidth = int(dm.W)
		g.screenHeight = int(dm.H)
	}

	g.setupConsoles()

	g.ui = ui.NewUI(g.renderer, g.defaultFont, fontSize)
	g.ui.SetScreenDimensions(g.screenWidth, g.screenHeight)

	g.gameState = mainMenu

	g.mainMenu = &MainMenu{}

	g.setupInput()
	g.setupGame()
}

// Shutdown should be called when the program quits.
func (g *Game) Shutdown() {
	g.glyphTexture.Destroy()

	g.defaultFont.Close()
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
	ttf.Quit()
}

func (g *Game) createGlyphTexture() {
	var err error
	g.glyphTexture, err = assets.NewGlyphTexture(g.renderer, "./assets/textures/ascii_ext_courier.png", "./assets/textures/ascii_ext_courier.json")
	if err != nil {
		log.Fatalf("Error creating glyph texture")
	}
}

func (g *Game) createPlayer() {
	e := entity.NewEntity("Player", &components.Appearance{Char: "@", Color: utils.ColorRGBA{R: 0, G: 128, B: 255, A: 255}}, utils.Vec2{}, true)
	e.Combat = &components.Combat{Power: 5, Defense: 2}
	e.Health = &components.Health{CurrentHP: 40, HP: 40}
	e.Vision = &components.Vision{Range: 20}
	g.entities = append(g.entities, e)
	g.player = e
}

// Occupied returns true if the given tile is occupied by a blocking entity and if the tile is currently visible.
func (g *Game) Occupied(p utils.Vec2) bool {
	for _, e := range g.entities {
		if e.Position != nil && e.Position.Current.X == p.X && e.Position.Current.Y == p.Y && e.IsBlocking != nil && g.player.FoV.Seen(p) && g.player.FoV.Visible(p) {
			return true
		}
	}
	return false
}

func (g *Game) setTargetPosition(x, y int32) {
	g.movementPath = g.determinePathPlayerMouse()
	g.player.TargetPosition = utils.Vec2{X: x, Y: y}
}

func (g *Game) drawMainMenu() {
	g.renderer.GetRenderer().SetClipRect(nil)
	g.renderer.SetScale(1, 1)
	g.renderer.Clear()

	g.consoleMainMenu.Clear()
	g.mainMenu.Render(g.consoleMainMenu, g.gameInProgress)
	g.consoleMainMenu.Render()

	g.renderer.Present()
}

func (g *Game) draw() {
	g.renderer.GetRenderer().SetClipRect(nil)
	g.renderer.SetScale(g.renderScale, g.renderScale)
	g.renderer.Clear()

	g.consoleMap.Clear()
	g.currentGameMap.Render(g.consoleMap, g.player.FoV, g.player, g.entities, int32(g.renderer.OriginX), int32(g.renderer.OriginY))
	if g.gameState != gameOver && g.gameState != mainMenu {
		g.renderMouseTile()
	}
	g.consoleMap.Render()

	g.renderer.SetScale(1, 1)
	g.ui.Render()

	g.renderer.Present()
}

func (g *Game) timestep() {
	if g.nextStep {
		g.movementSystem(playersTurn)
		g.movementSystem(enemyTurn)

		g.pickupSystem()
		g.useSystem()
		g.regenerationSystem()

		g.updateFoVs()

		g.time++
		g.nextStep = false

		g.ui.SetStatusBarText("")
		g.updateUI()
	}
}

func (g *Game) updateFoVs() {
	for _, e := range g.entities {
		if e.Position == nil || e.Vision == nil {
			continue
		}
		if e.Mutations.Has(components.MutationEffectXRay) {
			fov.UpdateFoV(g.currentGameMap, e.FoV, e.Vision.Range+e.Mutations.GetData(components.MutationEffectIncreasedVision), e.Position.Current, true)
		} else {
			fov.UpdateFoV(g.currentGameMap, e.FoV, e.Vision.Range+e.Mutations.GetData(components.MutationEffectIncreasedVision), e.Position.Current, false)
		}
	}
}
