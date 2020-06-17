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

	movementPath []components.Position

	player   *entity.Entity
	entities []*entity.Entity

	time uint

	nextStep  bool
	gameState gameState

	ui             *ui.UI
	commandManager *commandManager

	consoleMap *console.MatrixConsole
}

// Setup should be called first after creating an instance of Game.
func (g *Game) Setup(windowWidth, windowHeight int, fullscreen bool) {
	log.Printf("Setting up game...")
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

	g.gameState = playersTurn

	g.setupInput()
	g.setupGame()
	log.Printf("Done setting up game")
	g.ui.SetStatusBarText("Done setting up game")
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
	e := entity.NewEntity("Player", "@", utils.ColorRGB{R: 0, G: 128, B: 255}, components.Position{}, true)
	e.Combat = &components.Combat{CurrentHP: 40, HP: 40, Power: 5, Defense: 2}
	e.VisibilityRange = 20
	g.entities = append(g.entities, e)
	g.player = e
}

func (g *Game) createEnemy(name string, char string, color utils.ColorRGB, p components.Position) *entity.Entity {
	e := entity.NewEntity(name, char, color, p, true)
	e.TargetPosition = p
	return e
}

// Occupied returns true if the given tile is occupied by a blocking entity and if the tile is currently visible.
func (g *Game) Occupied(p components.Position) bool {
	for _, e := range g.entities {
		if e.Position != nil && e.Position.X == p.X && e.Position.Y == p.Y && e.Blocks && g.player.FoV.Seen(p) && g.player.FoV.Visible(p) {
			return true
		}
	}
	return false
}

func (g *Game) setTargetPosition(x, y int) {
	g.movementPath = g.determinePathPlayerMouse()
	g.player.TargetPosition = components.Position{X: x, Y: y}
}

func (g *Game) draw() {
	g.renderer.GetRenderer().SetClipRect(nil)
	g.renderer.SetScale(g.renderScale, g.renderScale)
	g.renderer.Clear()

	g.consoleMap.Clear()
	g.currentGameMap.Render(g.consoleMap, g.player.FoV, g.player, g.entities, int32(g.renderer.OriginX), int32(g.renderer.OriginY))
	if g.gameState != gameOver {
		g.renderMouseTile()
	}
	g.consoleMap.Render()

	g.renderer.SetScale(1, 1)
	g.ui.Render()

	g.renderer.Present()
}

func (g *Game) timestep() {
	if g.nextStep {
		g.updatePositions(playersTurn)
		g.updatePositions(enemyTurn)
		g.updateFoVs()

		g.time++
		g.nextStep = false

		g.ui.SetStatusBarText("")
		g.updateUI()
	}
}

func (g *Game) updateFoVs() {
	for _, e := range g.entities {
		if e.Position == nil {
			continue
		}
		if e.Mutations.Has(components.MutationEffectXRay) {
			fov.UpdateFoV(g.currentGameMap, e.FoV, e.VisibilityRange+e.Mutations.GetData(components.MutationEffectIncreasedVision), *e.Position, true)
		} else {
			fov.UpdateFoV(g.currentGameMap, e.FoV, e.VisibilityRange+e.Mutations.GetData(components.MutationEffectIncreasedVision), *e.Position, false)
		}
	}
}
