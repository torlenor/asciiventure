package game

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
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
	g.currentGamMapID = r
	r--
	if r < 0 || r >= len(g.loadedGameMaps) {
		return
	}

	g.currentGameMap = g.loadedGameMaps[r]

	g.player.FoV.ClearSeen()
	g.entities = []*entity.Entity{g.player}
	g.player.Position = &components.Position{X: (g.currentGameMap.SpawnPoint.X), Y: g.currentGameMap.SpawnPoint.Y}
	g.player.TargetPosition = *g.player.Position
	g.createEnemyEntities()
	g.createItems()
	g.createMutagens()
	g.updateFoVs()
	g.ui.AddLogEntry("Map changed.")

	g.updateUI()
}
