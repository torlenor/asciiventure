package game

import (
	"github.com/torlenor/asciiventure/entity"
)

func (g *Game) createMouse() *entity.Entity {
	return entity.ParseMonster("./data/monsters/mouse.json")
}

func (g *Game) createDog() *entity.Entity {
	return entity.ParseMonster("./data/monsters/dog.json")
}
