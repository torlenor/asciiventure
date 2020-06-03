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
	screenHeight = 768

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

	mouseTileX int32
	mouseTileY int32

	markedPath   []components.Position
	movementPath []components.Position

	player   *entity.Entity
	entities []*entity.Entity

	time uint32

	nextStep  bool
	gameState gameState

	characterWindow *ui.TextWidget
	logWindow       *ui.TextWidget
}

// Setup should be called first after creating an instance of Game.
func (g *Game) Setup() {
	log.Printf("Setting up game...")
	g.renderScale = 0.8

	g.setupWindow()
	g.setupRenderer()

	g.gameState = playersTurn

	g.characterWindow = ui.NewTextWidget(g.renderer, g.defaultFont, &characterWindowRect)
	g.characterWindow.SetWrapLength(int(characterWindowRect.W))
	g.logWindow = ui.NewTextWidget(g.renderer, g.defaultFont, &logWindowRect)
	g.logWindow.SetWrapLength(int(logWindowRect.W))

	g.setupGame()
	log.Printf("Done setting up game")
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
		e := entity.NewEntity("Player", gl, components.Position{}, true)
		e.Combat = &components.Combat{MaxHP: 20, HP: 20, Power: 5, Defense: 2}
		g.entities = append(g.entities, e)
		g.player = e
	} else {
		log.Printf("Unable to add player entity")
	}
}

func (g *Game) createMouse(p components.Position) *entity.Entity {
	if gl, ok := g.glyphTexture.Get("m"); ok {
		gl.Color = components.ColorRGB{R: 200, G: 200, B: 200}
		e := g.createEnemy("Mouse", gl, p)
		e.Combat = &components.Combat{MaxHP: 2, HP: 2, Power: 1, Defense: 0}
		return e
	}
	log.Printf("Unable to add mouse entity")
	return nil
}

func (g *Game) createDog(p components.Position) *entity.Entity {
	if gl, ok := g.glyphTexture.Get("d"); ok {
		gl.Color = components.ColorRGB{R: 255, G: 0, B: 0}
		e := g.createEnemy("Dog", gl, p)
		e.Combat = &components.Combat{MaxHP: 10, HP: 10, Power: 5, Defense: 2}
		return e
	}
	log.Printf("Unable to add dog entity")
	return nil
}

func (g *Game) createEnemy(name string, gl components.Glyph, p components.Position) *entity.Entity {
	e := entity.NewEntity(name, gl, p, true)
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
		x := int32(rand.Intn(int(maxx)))
		y := int32(rand.Intn(int(maxy)))
		if g.Occupied(components.Position{X: x, Y: y}) || !g.currentGameMap.Empty(x, y) {
			continue
		}
		if rand.Intn(100) < 50 {
			e := g.createMouse(components.Position{X: x, Y: y})
			if e != nil {
				g.entities = append(g.entities, e)
			}
		} else {
			e := g.createDog(components.Position{X: x, Y: y})
			if e != nil {
				g.entities = append(g.entities, e)
			}
		}
	}
}

func (g *Game) renderEntities() {
	for _, e := range g.entities {
		if e == g.player || (e.Position.X == g.player.Position.X && e.Position.Y == g.player.Position.Y) {
			continue
		}
		if g.player.FoV.Visible(e.Position) {
			g.renderer.RenderGlyph(e.Glyph, e.Position.X, e.Position.Y)
		}
	}
	g.renderer.RenderGlyph(g.player.Glyph, g.player.Position.X, g.player.Position.Y)
}

func (g *Game) setTargetPosition(x, y int32) {
	g.movementPath = g.markedPath
	g.player.TargetPosition = components.Position{X: x, Y: y}
}

func (g *Game) draw() {
	g.renderer.SetScale(g.renderScale, g.renderScale)
	g.renderer.Clear()

	// g.renderer.Copy(g.mapTexture, nil, &sdl.Rect{X: 0, Y: 0, W: int32(screenWidth / g.renderScale), H: int32(screenHeight / g.renderScale)})
	// We are actually rendering it in total again because of FoV updates and some flickering which we encountered when pre-rendering
	g.currentGameMap.Render(g.renderer, g.player.FoV, g.renderer.OriginX, g.renderer.OriginY)
	g.renderEntities()
	if g.gameState != gameOver {
		g.renderMouseTile()
	}
	g.renderer.SetScale(1, 1)
	g.characterWindow.Render()
	g.logWindow.Render()

	g.renderer.Present()
}

func (g *Game) updateCharacterWindow() {
	g.characterWindow.SetText([]string{fmt.Sprintf("Time: %d", g.time), fmt.Sprintf("Health: %d/%d", g.player.Combat.HP, g.player.Combat.MaxHP)})
}

func (g *Game) timestep() {
	if g.nextStep {
		g.updatePositions(playersTurn)
		g.updatePositions(enemyTurn)
		fov.UpdateFoV(g.currentGameMap, g.player.FoV, playerViewRange, g.player.Position)
		g.time++
		g.updateCharacterWindow()
		g.nextStep = false
	}
}
