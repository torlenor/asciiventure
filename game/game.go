package game

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/gamemap"
	"github.com/torlenor/asciiventure/renderers"
	"github.com/torlenor/asciiventure/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TODO: Some of these consts should be made configurable
const (
	windowName = "Asciiventure"

	fontPath = "./assets/fonts/RobotoMono-Regular.ttf"
	fontSize = 16

	screenWidth  = 1366
	screenHeight = 960

	latticeDX = 19
	latticeDY = 32

	playerViewRange = 20
)

type gameState int

const (
	playersTurn gameState = iota
	enemyTurn
	gameOver
)

func (d gameState) String() string {
	return [...]string{"playersTurn", "enemyTurn", "gameOver"}[d]
}

var (
	roomRenderPane      = sdl.Rect{X: screenHeight / 6, Y: 0, W: screenWidth, H: screenHeight - screenHeight/6}
	characterWindowRect = sdl.Rect{X: 0, Y: 0, W: screenWidth / 2, H: screenHeight / 6}
	logWindowRect       = sdl.Rect{X: screenWidth - screenWidth/2, Y: 0, W: screenWidth / 2, H: screenHeight / 6}
	statusBarRec        = sdl.Rect{X: 0, Y: screenHeight - fontSize - 16, W: screenWidth, H: fontSize + 16}

	mutationsRect = sdl.Rect{X: screenWidth - screenWidth/4, Y: screenHeight / 6, W: screenWidth / 4, H: 3 * screenHeight / 6}
	inventoryRect = sdl.Rect{X: screenWidth - screenWidth/4, Y: 4 * screenHeight / 6, W: screenWidth / 4, H: 2*screenHeight/6 - statusBarRec.H}
)

// Game is the main struct of the game
type Game struct {
	quit bool

	window   *sdl.Window
	renderer *renderers.Renderer

	renderScale float32

	defaultFont  *ttf.Font
	glyphTexture *assets.GlyphTexture

	currentGameMap *gamemap.GameMap
	loadedGameMaps []*gamemap.GameMap
	mapTexture     *sdl.Texture

	mouseTileX int
	mouseTileY int

	movementPath []components.Position

	player   *entity.Entity
	entities []*entity.Entity

	time uint

	nextStep  bool
	gameState gameState

	characterWindow *ui.TextWidget
	logWindow       *ui.TextWidget
	statusBar       *ui.TextWidget
	mutations       *ui.TextWidget
	inventory       *ui.TextWidget
}

// Setup should be called first after creating an instance of Game.
func (g *Game) Setup() {
	log.Printf("Setting up game...")
	g.renderScale = 0.8

	g.setupWindow()
	g.setupRenderer()

	g.gameState = playersTurn

	g.characterWindow = ui.NewTextWidget(g.renderer, g.defaultFont, &characterWindowRect, true)
	g.characterWindow.SetWrapLength(int(characterWindowRect.W - 8))
	g.logWindow = ui.NewTextWidget(g.renderer, g.defaultFont, &logWindowRect, true)
	g.logWindow.SetWrapLength(int(logWindowRect.W - 8))
	g.statusBar = ui.NewTextWidget(g.renderer, g.defaultFont, &statusBarRec, true)
	g.statusBar.SetWrapLength(int(statusBarRec.W - 8))

	g.mutations = ui.NewTextWidget(g.renderer, g.defaultFont, &mutationsRect, true)
	g.mutations.SetWrapLength(int(mutationsRect.W - 8))
	g.mutations.AddRow("No mutations")
	g.inventory = ui.NewTextWidget(g.renderer, g.defaultFont, &inventoryRect, true)
	g.inventory.SetWrapLength(int(inventoryRect.W - 8))
	g.inventory.AddRow("Inventory Empty")

	g.setupGame()
	log.Printf("Done setting up game")
	g.statusBar.AddRow("Done setting up game")
}

// Shutdown should be called when the program quits.
func (g *Game) Shutdown() {
	g.glyphTexture.Destroy()
	g.mapTexture.Destroy()

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
	if gl, ok := g.glyphTexture.Get("@"); ok {
		gl.Color = components.ColorRGB{R: 0, G: 128, B: 255}
		e := entity.NewEntity("Player", "@", components.ColorRGB{R: 0, G: 128, B: 255}, components.Position{}, true)
		e.Combat = &components.Combat{CurrentHP: 40, HP: 40, Power: 5, Defense: 2}
		g.entities = append(g.entities, e)
		g.player = e
	} else {
		log.Printf("Unable to add player entity")
	}
}

func (g *Game) createEnemy(name string, char string, color components.ColorRGB, p components.Position) *entity.Entity {
	e := entity.NewEntity(name, char, color, p, true)
	e.TargetPosition = p
	return e
}

// Occupied returns true if the given tile is occupied by a blocking entity and if the tile is currently visible.
func (g *Game) Occupied(p components.Position) bool {
	for _, e := range g.entities {
		if e.Position.X == p.X && e.Position.Y == p.Y && e.Blocks && g.player.FoV.Seen(p) && g.player.FoV.Visible(p) {
			return true
		}
	}
	return false
}

func (g *Game) createEnemyEntities() {
	maxx, maxy := g.currentGameMap.Dimensions()
	for i := 0; i < 5; i++ {
		p := components.Position{X: rand.Intn(maxx), Y: rand.Intn(maxy)}
		if g.Occupied(p) || !g.currentGameMap.Empty(p.X, p.Y) {
			continue
		}
		var e *entity.Entity
		if rand.Intn(100) < 50 {
			e = g.createMouse()

		} else {
			e = g.createDog()
		}
		if e != nil {
			e.Position = p
			e.InitialPosition = p
			e.TargetPosition = p
			g.entities = append(g.entities, e)
		} else {
			log.Printf("Error creating Mouse entity")
		}
	}
}

func (g *Game) createItems() {
	maxx, maxy := g.currentGameMap.Dimensions()
	for i := 0; i < 20; i++ {
		p := components.Position{X: rand.Intn(maxx), Y: rand.Intn(maxy)}
		if g.Occupied(p) || !g.currentGameMap.Empty(p.X, p.Y) {
			continue
		}
		var e *entity.Entity
		if rand.Intn(100) < 50 {
			e = entity.ParseItem("./data/items/healingpotion.json")
		} else {
			e = entity.ParseItem("./data/items/inventory_mutagen.json")
		}
		if e != nil {
			e.Position = p
			e.InitialPosition = p
			e.TargetPosition = p
			e.Blocks = false
			g.entities = append(g.entities, e)
		} else {
			log.Printf("Error creating Item entity")
		}
	}
}

func (g *Game) renderChar(char string, color components.ColorRGB, p components.Position) {
	if gl, ok := g.glyphTexture.Get(char); ok {
		gl.Color = color
		g.renderer.RenderGlyph(gl, p.X, p.Y)
	} else {
		log.Printf("Unable to render '%s'. Glyph not found.", char)
	}
}

func (g *Game) renderEntity(e *entity.Entity) {
	if e.Dead {
		g.renderChar("%", components.ColorRGB{R: 150, G: 150, B: 150}, e.Position)
	} else {
		g.renderChar(e.Char, e.Color, e.Position)
	}
}

func (g *Game) renderEntities() {
	for _, e := range g.entities {
		if e == g.player || (e.Position.X == g.player.Position.X && e.Position.Y == g.player.Position.Y) {
			continue
		}
		if g.player.FoV.Visible(e.Position) {
			g.renderEntity(e)
		}
	}
	g.renderEntity(g.player)
}

func (g *Game) setTargetPosition(x, y int) {
	g.movementPath = g.determinePathPlayerMouse()
	g.player.TargetPosition = components.Position{X: x, Y: y}
}

func (g *Game) draw() {
	g.renderer.SetScale(g.renderScale, g.renderScale)
	g.renderer.Clear()

	// g.renderer.Copy(g.mapTexture, nil, &sdl.Rect{X: 0, Y: 0, W: int(screenWidth / g.renderScale), H: int(screenHeight / g.renderScale)})
	// We are actually rendering it in total again because of FoV updates and some flickering which we encountered when pre-rendering
	g.currentGameMap.Render(g.renderer, g.player.FoV, g.renderer.OriginX, g.renderer.OriginY)
	g.renderEntities()
	if g.gameState != gameOver {
		g.renderMouseTile()
	}
	g.renderer.SetScale(1, 1)
	g.characterWindow.Render()
	g.logWindow.Render()
	g.statusBar.Render()
	g.mutations.Render()
	if g.player.Mutations.Has(components.MutationInventory) {
		g.inventory.Render()
	}

	g.renderer.Present()
}

func (g *Game) updateCharacterWindow() {
	g.characterWindow.SetText([]string{fmt.Sprintf("Time: %d", g.time), fmt.Sprintf("HP: %d/%d", g.player.Combat.CurrentHP, g.player.Combat.HP)})
}

func (g *Game) timestep() {
	if g.nextStep {
		g.updatePositions(playersTurn)
		g.updatePositions(enemyTurn)
		fov.UpdateFoV(g.currentGameMap, g.player.FoV, playerViewRange, g.player.Position)
		g.time++
		g.updateCharacterWindow()
		g.statusBar.SetText([]string{})
		g.nextStep = false
	}
}

func (g *Game) updateStatusBar() {
	g.statusBar.Clear()
	for _, e := range g.entities {
		if e.Position.Equal(components.Position{X: g.mouseTileX, Y: g.mouseTileY}) && e != g.player {
			if e.Dead {
				g.statusBar.AddRow(e.Name + "(Dead)")
			} else if e.Item != nil {
				if e.Item.CanPickup {
					g.statusBar.AddRow(e.Name + ": Pick up item with 'g'")
				} else if e.Item.Mutagen {
					g.statusBar.AddRow(e.Name + ": Consume mutagen with 'g'")
				}
			} else {
				g.statusBar.AddRow(e.Name)
			}
			return
		}
	}
	g.statusBar.Clear()
}

func (g *Game) updateMutationsPane() {
	g.mutations.Clear()
	g.mutations.AddRow("Mutations:")
	g.mutations.AddRow("----------------")
	for _, m := range g.player.Mutations {
		g.mutations.AddRow(m.String())
	}
}
