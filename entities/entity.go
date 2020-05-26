package entities

// Entity defines an entity in our entity-component system
type Entity uint64

var lastEntity Entity = 0

// NewEntity creates a new unique entity
func NewEntity() Entity {
	lastEntity++
	return lastEntity
}
