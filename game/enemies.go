package game

import (
	"log"
	"math/rand"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

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

func (g *Game) createMouse() *entity.Entity {
	return entity.ParseMonster("./data/monsters/mouse.json")
}

func (g *Game) createDog() *entity.Entity {
	return entity.ParseMonster("./data/monsters/dog.json")
}
