package game

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entities"
	"github.com/torlenor/asciiventure/maps"
	"github.com/torlenor/asciiventure/renderers"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TODO: Some of these consts should be made configurable
const (
	windowName = "Asciiventure"

	fontPath = "./assets/fonts/Roboto-Regular.ttf"
	fontSize = 24

	screenWidth  = 820
	screenHeight = 1000

	latticeDX = 19
	latticeDY = 32

	playerViewRange = 10
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

	roomEntities   []entities.Entity
	collision      components.CollisionManager
	glyph          components.GlyphManager
	position       components.PositionManager
	targetPosition components.PositionManager
	player         entities.Entity

	time uint32

	nextStep bool
}

// Setup should be called first after creating an instance of Game.
func (g *Game) Setup() {
	g.renderScale = 1.0

	g.collision = make(components.CollisionManager)
	g.glyph = make(components.GlyphManager)
	g.position = make(components.PositionManager)
	g.targetPosition = make(components.PositionManager)

	g.setupWindow()
	g.setupRenderer()
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

	g.currentRoom.UpdateFoV(playerViewRange, g.position[g.player].X, g.position[g.player].Y)

	for !g.quit {
		start := time.Now()
		g.handleSDLEvents()
		g.timestep()
		gameLogicUpdateMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0

		start = time.Now()
		g.draw()
		drawUpdateMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0
		start = time.Now()
		<-ticker.C
		spareMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0
		if true {
			fmt.Printf("Game logic duration: %.2f ms, draw duration: %.2f ms, total: %.2f ms, spare: %.2f ms\n", gameLogicUpdateMs, drawUpdateMs, gameLogicUpdateMs+drawUpdateMs, spareMs)
		}
	}

	ticker.Stop()
}

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
}

func (g *Game) createGlyphTexture() {
	var err error
	g.glyphTexture, err = assets.NewGlyphTexture(g.renderer, "./assets/textures/ascii_ext_courier.png", "./assets/textures/ascii_ext_courier.json")
	if err != nil {
		log.Fatalf("Error creating glyph texture")
	}
}

func (g *Game) loadRoomsFromDirectory(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	g.loadedRooms = []*maps.Room{}
	for _, f := range files {
		if !f.IsDir() {
			f, err := os.Open(dir + "/" + f.Name())
			if err != nil {
				log.Printf("Error opening %s: %s", f.Name(), err)
				continue
			}
			r4 := bufio.NewReader(f)
			r, err := maps.NewRoom(r4, g.glyphTexture)
			if err != nil {
				log.Printf("Error reading room file: %s", err)
				continue
			}

			g.loadedRooms = append(g.loadedRooms, &r)
		}
	}
}
func (g *Game) createPlayer() {
	if gl, ok := g.glyphTexture.Get("@"); ok {
		e := entities.NewEntity()
		g.position[e] = components.Position{}
		gl.Color = components.ColorRGB{R: 0, G: 128, B: 255}
		g.glyph[e] = gl
		g.collision[e] = components.Collision{V: true}
		g.roomEntities = append(g.roomEntities, e)
		g.player = e
	} else {
		log.Printf("Unable to add player entity")
	}
}

func (g *Game) createEnemy(p components.Position) {
	if gl, ok := g.glyphTexture.Get("e"); ok {
		e := entities.NewEntity()
		g.position[e] = p
		gl.Color = components.ColorRGB{R: 255, G: 0, B: 0}
		g.glyph[e] = gl
		g.collision[e] = components.Collision{V: true, DestroyOnCollision: true}
		g.roomEntities = append(g.roomEntities, e)
	} else {
		log.Printf("Unable to add enemy entity")
	}
}

func (g *Game) createEnemyEntities() {
	g.createEnemy(components.Position{X: 2, Y: 2})
}

func (g *Game) selectRoom(r int) {
	r--
	if r < 0 || r >= len(g.loadedRooms) {
		return
	}

	g.currentRoom = g.loadedRooms[r]

	g.preRenderRoom()

	g.position[g.player] = components.Position{X: (g.currentRoom.SpawnPoint.X), Y: g.currentRoom.SpawnPoint.Y}
	g.targetPosition[g.player] = g.position[g.player]

	g.currentRoom.ClearSeen()
	g.currentRoom.UpdateFoV(playerViewRange, g.position[g.player].X, g.position[g.player].Y)
}

func (g *Game) removeEntity(e entities.Entity) {
	delete(g.collision, e)
	delete(g.glyph, e)
	delete(g.position, e)
}

func (g *Game) isEmpty(x, y int32) bool {
	_, o := g.occupied(x, y)
	return g.currentRoom.Empty(x, y) && !o
}

func (g *Game) updateMouse(x, y int32) {
	g.mouseTileX = int32((float32(x)+0.5)/latticeDX/g.renderScale) - g.renderer.OriginX
	g.mouseTileY = int32((float32(y)+0.5)/latticeDY/g.renderScale) - g.renderer.OriginY
	g.markedPath = g.determineLatticePathPlayerMouse()
}

func (g *Game) renderEntities() {
	for id, gl := range g.glyph {
		if p, ok := g.position[id]; ok {
			g.renderer.RenderGlyph(gl, p.X, p.Y)
		}
	}
}

func (g *Game) determineLatticePathPlayerMouse() []components.Position {
	return determineLatticePath(components.Position{X: g.position[g.player].X, Y: g.position[g.player].Y}, components.Position{X: g.mouseTileX, Y: g.mouseTileY})
}

func (g *Game) renderMouseTile() {
	if g.position[g.player].X != g.targetPosition[g.player].X ||
		g.position[g.player].Y != g.targetPosition[g.player].Y {
		path := determineLatticePath(g.position[g.player], g.targetPosition[g.player])
		for _, p := range path {
			_, occupied := g.occupied(p.X, p.Y)
			notEmpty := !g.currentRoom.Empty(p.X, p.Y)
			color := components.ColorRGBA{R: 100, G: 100, B: 255, A: 128}
			if occupied || notEmpty {
				color = components.ColorRGBA{R: 255, G: 100, B: 100, A: 128}
			}
			g.renderer.FillCharCoordinate(p.X, p.Y, color)
			if notEmpty || occupied {
				break
			}
		}
		if gl, ok := g.glyphTexture.Get("X"); ok {
			gl.Color = components.ColorRGB{R: 255, G: 0, B: 0}
			g.renderer.RenderGlyph(gl, g.targetPosition[g.player].X, g.targetPosition[g.player].Y)
		}
	}

	for _, p := range g.markedPath {
		_, occupied := g.occupied(p.X, p.Y)
		notEmpty := !g.currentRoom.Empty(p.X, p.Y)
		color := components.ColorRGBA{R: 100, G: 255, B: 100, A: 128}
		if occupied || notEmpty {
			color = components.ColorRGBA{R: 255, G: 100, B: 100, A: 128}
		}
		g.renderer.FillCharCoordinate(p.X, p.Y, color)
		if notEmpty || occupied {
			break
		}
	}
	if len(g.markedPath) > 0 {
		p := g.markedPath[len(g.markedPath)-1]
		color := components.ColorRGBA{R: 128, G: 128, B: 128, A: 180}
		g.renderer.FillCharCoordinate(p.X, p.Y, color)
	}
}

func (g *Game) preRenderRoom() {
	var err error
	g.mapTexture, err = g.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_TARGET, int32(screenWidth/g.renderScale), int32(screenHeight/g.renderScale))
	if err != nil {
		log.Printf("Error creating texture: %s", err)
	}
	err = g.renderer.SetRenderTarget(g.mapTexture)
	g.renderer.Clear()
	if err != nil {
		log.Printf("Error setting texture as render target: %s", err)
	}
	g.currentRoom.Render(g.renderer, g.renderer.OriginX, g.renderer.OriginY)
	g.renderer.Present()
	g.renderer.SetRenderTarget(nil)
}

func (g *Game) occupied(x, y int32) (entities.Entity, bool) {
	for e, p := range g.position {
		if p.X == x && p.Y == y {
			return e, true
		}
	}
	return 0, false
}

func (g *Game) setPlayerTargetPosition(x, y int32) {
	g.targetPosition[g.player] = components.Position{X: x, Y: y}
}

func (g *Game) draw() {
	g.renderer.SetScale(g.renderScale, g.renderScale)
	g.renderer.Clear()

	// g.renderer.Copy(g.mapTexture, nil, &sdl.Rect{X: 0, Y: 0, W: int32(screenWidth / g.renderScale), H: int32(screenHeight / g.renderScale)})
	// We are actually rendering it in total again because of FoV updates and some flickering which we encountered when pre-rendering
	g.currentRoom.Render(g.renderer, g.renderer.OriginX, g.renderer.OriginY)
	g.renderEntities()
	g.renderMouseTile()
	g.drawGameTime()

	g.renderer.Present()
}

func (g *Game) drawGameTime() {
	// TODO: Rendering game time is slow
	surface, err := g.defaultFont.RenderUTF8Blended(fmt.Sprintf("Time: %d", g.time), sdl.Color{R: 0, G: 255, B: 0, A: 255})
	if err != nil {
		log.Printf("Error rendering game time: %s", err)
		return
	}

	t, err := g.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		fmt.Printf("Failed to create texture from surface when trying to render game time: %s\n", err)
		return
	}

	g.renderer.Copy(t, nil, &sdl.Rect{X: int32(float32(screenWidth-surface.W) / g.renderScale), Y: 0, W: surface.W, H: surface.H})

	surface.Free()
	t.Destroy()
}

func (g *Game) timestep() {
	if g.nextStep {
		g.updatePositions()
		g.currentRoom.UpdateFoV(playerViewRange, g.position[g.player].X, g.position[g.player].Y)
		g.time++
		g.nextStep = false
	}
}
