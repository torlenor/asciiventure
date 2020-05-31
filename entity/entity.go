package entity

import "github.com/torlenor/asciiventure/components"

// Entity defines an entity in our entity-component system
type Entity struct {
	Name  string
	Glyph components.Glyph

	Position       components.Position
	TargetPosition components.Position

	Combat *components.Combat

	Blocks bool
	IsDead bool
}

// NewEntity creates a new unique entity
func NewEntity(name string, glyph components.Glyph, initPosition components.Position, blocks bool) *Entity {
	return &Entity{
		Name:     name,
		Glyph:    glyph,
		Position: initPosition,
		Blocks:   blocks,
	}
}

// Move moves the entity by (dy,dy)
func (e *Entity) Move(dx, dy int32) {
	e.Position.X += dx
	e.Position.Y += dy
}

//MoveTo moves the entity to (x,y)
func (e *Entity) MoveTo(p components.Position) {
	e.Position = p
}
