package game

import (
	"log"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

func (g *Game) createMouse(p components.Position) *entity.Entity {
	if gl, ok := g.glyphTexture.Get("m"); ok {
		gl.Color = components.ColorRGB{R: 200, G: 200, B: 200}
		e := g.createEnemy("Mouse", gl, p)
		e.Combat = &components.Combat{MaxHP: 2, HP: 2, Power: 1, Defense: 0}
		e.AttackRange = 4
		e.AttackRangeUntil = 10
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
		e.AttackRange = 10
		e.AttackRangeUntil = 30
		return e
	}
	log.Printf("Unable to add dog entity")
	return nil
}
