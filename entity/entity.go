package entity

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/fov"
)

// Entity defines an entity in our entity-component system
type Entity struct {
	Name string

	Char  string
	Color components.ColorRGB

	InitialPosition components.Position
	Position        components.Position
	TargetPosition  components.Position

	Combat   *components.Combat
	AI       *components.AI
	Item     *components.Item
	Mutation *components.Mutation

	Blocks bool
	Dead   bool

	FoV fov.FoVMap

	VisibilityRange int

	Mutations components.Mutations
	Inventory []*Entity
}

// NewEntity creates a new unique entity.
func NewEntity(name string, char string, color components.ColorRGB, initPosition components.Position, blocks bool) *Entity {
	return &Entity{
		Name:            name,
		Char:            char,
		Color:           color,
		Position:        initPosition,
		Blocks:          blocks,
		FoV:             make(fov.FoVMap),
		InitialPosition: initPosition,
	}
}

// Move moves the entity by (dy,dy).
func (e *Entity) Move(dx, dy int) {
	e.Position.X += dx
	e.Position.Y += dy
}

// MoveTo moves the entity to (y,y).
func (e *Entity) MoveTo(p components.Position) {
	e.Position = p
}

func (e *Entity) Attack(target *Entity) (results []components.CombatResult) {
	if target.Combat == nil {
		return
	}

	dmg := e.Combat.Power - target.Combat.Defense
	if dmg < 0 {
		dmg = 0
	}
	results = append(results, components.CombatResult{Type: components.TakeDamage, IntegerValue: dmg})

	return
}
