package game

import (
	"log"
	"math/rand"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/utils"
)

func (g *Game) createItems() {
	maxx, maxy := g.currentGameMap.Dimensions()
	for i := 0; i < 20; i++ {
		p := utils.Vec2{X: int32(rand.Intn(int(maxx))), Y: int32(rand.Intn(int(maxy)))}
		if g.Occupied(p) || !g.currentGameMap.Empty(p) {
			continue
		}
		var e *entity.Entity
		if rand.Intn(100) < 50 {
			e = entity.ParseItem("./data/items/healingpotion.json")
		} else {
			continue
		}
		if e != nil {
			e.Position = &components.Position{Current: p, Initial: p}
			e.TargetPosition = p
			g.entities = append(g.entities, e)
		} else {
			log.Printf("Error creating Item entity")
		}
	}
}
