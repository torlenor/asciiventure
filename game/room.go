package game

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/maps"
	"github.com/veandco/go-sdl2/sdl"
)

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

func (g *Game) selectRoom(r int) {
	r--
	if r < 0 || r >= len(g.loadedRooms) {
		return
	}

	g.currentRoom = g.loadedRooms[r]

	g.preRenderRoom()

	g.currentRoom.ClearSeen()
	g.entities = []*entity.Entity{}
	// TODO: Should not create a new player when this is going to be used as map transition
	g.createPlayer()
	g.player.Position = components.Position{X: (g.currentRoom.SpawnPoint.X), Y: g.currentRoom.SpawnPoint.Y}
	g.player.TargetPosition = g.player.Position
	g.createEnemyEntities()
	g.currentRoom.UpdateFoV(playerViewRange, g.player.Position.X, g.player.Position.Y)
	g.gameState = playersTurn
	g.logWindow.SetText([]string{})
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
