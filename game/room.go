package game

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/gamemap"
)

func (g *Game) loadGameMapsFromDirectory(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			ext := path.Ext(f.Name())
			if ext != ".map" {
				continue
			}
			f, err := os.Open(dir + "/" + f.Name())
			if err != nil {
				log.Printf("Error opening %s: %s", f.Name(), err)
				continue
			}
			r4 := bufio.NewReader(f)
			r, err := gamemap.NewGameMapFromReader(r4, g.glyphTexture)
			if err != nil {
				log.Printf("Error reading room file: %s", err)
				continue
			}

			g.loadedGameMaps = append(g.loadedGameMaps, &r)
		}
	}
}

func (g *Game) selectGameMap(r int) {
	if len(g.loadedGameMaps) == 0 {
		log.Fatalf("No maps loaded")
	}
	r--
	if r < 0 || r >= len(g.loadedGameMaps) {
		return
	}

	g.currentGameMap = g.loadedGameMaps[r]

	g.player.FoV.ClearSeen()
	g.entities = []*entity.Entity{g.player}
	g.player.Position = components.Position{X: (g.currentGameMap.SpawnPoint.X), Y: g.currentGameMap.SpawnPoint.Y}
	g.player.TargetPosition = g.player.Position
	g.createEnemyEntities()
	g.createItems()
	g.createMutagens()
	fov.UpdateFoV(g.currentGameMap, g.player.FoV, g.player.VisibilityRange, g.player.Position)
	g.gameState = playersTurn
	g.logWindow.SetText([]string{})

	g.focusPlayer()
	g.updateCharacterWindow()
	g.updateInventory()
	g.updateMutations()
	g.updateMutationsPane()
}

// func (g *Game) preRenderGameMap() {
// 	var err error
// 	g.mapTexture, err = g.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888,
// 		sdl.TEXTUREACCESS_TARGET, int(screenWidth/g.renderScale), int(screenHeight/g.renderScale))
// 	if err != nil {
// 		log.Printf("Error creating texture: %s", err)
// 	}
// 	err = g.renderer.SetRenderTarget(g.mapTexture)
// 	g.renderer.Clear()
// 	if err != nil {
// 		log.Printf("Error setting texture as render target: %s", err)
// 	}
// 	g.currentGameMap.Render(g.renderer, g.player.FoV, g.renderer.OriginX, g.renderer.OriginY)
// 	g.renderer.Present()
// 	g.renderer.SetRenderTarget(nil)
// }

func (g *Game) focusPlayer() {
	g.renderer.OriginX = -g.player.Position.X + int(screenWidth/latticeDX/2/g.renderScale)
	g.renderer.OriginY = -g.player.Position.Y + int((screenHeight+screenHeight/6)/latticeDY/2/g.renderScale)
}
