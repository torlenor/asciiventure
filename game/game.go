package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/maps"
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

	playerViewRange = 10
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

	currentRoom *maps.Room
	loadedRooms []*maps.Room
	mapTexture  *sdl.Texture

	mouseTileX int32
	mouseTileY int32

	markedPath []components.Position

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
	g.renderScale = 1.0

	g.setupWindow()
	g.setupRenderer()

	g.gameState = playersTurn

	g.characterWindow = ui.NewTextWidget(g.renderer, g.defaultFont, &characterWindowRect)
	g.characterWindow.SetWrapLength(int(characterWindowRect.W))
	g.logWindow = ui.NewTextWidget(g.renderer, g.defaultFont, &logWindowRect)
	g.logWindow.SetWrapLength(int(logWindowRect.W))
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

// GameLoop is a blocking function actually running the game.
func (g *Game) GameLoop() {
	g.createGlyphTexture()
	g.createPlayer()
	g.loadRoomsFromDirectory("./assets/rooms")
	g.selectRoom(1)

	ticker := time.NewTicker(time.Second / 15)

	g.currentRoom.UpdateFoV(playerViewRange, g.player.Position.X, g.player.Position.Y)
	g.updateCharacterWindow()
	g.logWindow.SetText([]string{"Welcome to <Epic Name Here>.", "A small cat takes a stroll and ends up in an epic adventure."})

	for !g.quit {
		start := time.Now()
		g.handleSDLEvents()
		if g.gameState != gameOver {
			g.timestep()
		}
		gameLogicUpdateMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0

		start = time.Now()
		g.draw()
		drawUpdateMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0
		start = time.Now()
		<-ticker.C
		spareMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0
		if false {
			fmt.Printf("Game logic duration: %.2f ms, draw duration: %.2f ms, total: %.2f ms, spare: %.2f ms\n", gameLogicUpdateMs, drawUpdateMs, gameLogicUpdateMs+drawUpdateMs, spareMs)
		}
	}

	ticker.Stop()
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

func (g *Game) occupied(x, y int32) bool {
	for _, e := range g.entities {
		if e.Position.X == x && e.Position.Y == y {
			return true
		}
	}
	return false
}

func (g *Game) createEnemyEntities() {
	maxx, maxy := g.currentRoom.Dimensions()
	for i := 0; i < 10; i++ {
		x := int32(rand.Intn(int(maxx)))
		y := int32(rand.Intn(int(maxy)))
		if g.occupied(x, y) || !g.currentRoom.Empty(x, y) {
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
		if g.currentRoom.Visible(e.Position.X, e.Position.Y) {
			g.renderer.RenderGlyph(e.Glyph, e.Position.X, e.Position.Y)
		}
	}
	g.renderer.RenderGlyph(g.player.Glyph, g.player.Position.X, g.player.Position.Y)
}

func (g *Game) setPlayerTargetPosition(x, y int32) {
	g.player.TargetPosition = components.Position{X: x, Y: y}
}

func (g *Game) draw() {
	g.renderer.SetScale(g.renderScale, g.renderScale)
	g.renderer.Clear()

	// g.renderer.Copy(g.mapTexture, nil, &sdl.Rect{X: 0, Y: 0, W: int32(screenWidth / g.renderScale), H: int32(screenHeight / g.renderScale)})
	// We are actually rendering it in total again because of FoV updates and some flickering which we encountered when pre-rendering
	g.currentRoom.Render(g.renderer, g.renderer.OriginX, g.renderer.OriginY)
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
		g.currentRoom.UpdateFoV(playerViewRange, g.player.Position.X, g.player.Position.Y)
		g.time++
		g.updateCharacterWindow()
		g.nextStep = false
	}
}