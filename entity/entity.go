package entity

import (
	"fmt"

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

func (e *Entity) PickUp(target *Entity) (result []ActionResult) {
	if target.Item != nil {
		if target.Item.CanPickup {
			if e.Mutations.Has(components.MutationEffectInventory) {
				e.Inventory = append(e.Inventory, target)
				result = append(result, ActionResult{Type: ActionResultItemPickedUp})
				result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s picked up %s.", e.Name, target.Name)})
			} else {
				result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("You are a cat, you cannot pick up things (or can you?).")})
			}
		}
	}
	return
}

func (e *Entity) ConsumeMutation(target *Entity) (result []ActionResult) {
	if target.Mutation != nil {
		if !e.Mutations.Has(target.Mutation.Effect) {
			e.Mutations = append(e.Mutations, *target.Mutation)
			result = append(result, ActionResult{Type: ActionResultMutationConsumed})
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s gained mutation %s.", e.Name, target.Name)})
		} else {
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s already has %s.", e.Name, target.Name)})
		}
	}

	return
}

// MoveTo moves the entity to (y,y).
func (e *Entity) MoveTo(p components.Position) {
	e.Position = p
}

func (e *Entity) Attack(target *Entity) (results []CombatResult) {
	if target.Combat == nil {
		return
	}

	dmg := e.Combat.Power - target.Combat.Defense
	if dmg < 0 {
		dmg = 0
	}
	results = append(results, CombatResult{Type: CombatResultTakeDamage, IntegerValue: dmg})

	return
}
