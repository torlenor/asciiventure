package game

import (
	"log"
	"math/rand"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

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
			continue
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