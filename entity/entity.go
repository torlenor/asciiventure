package entity

import (
	"fmt"
	"log"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/utils"
)

// Entity defines an entity in our entity-component system
type Entity struct {
	Name string

	Char  string
	Color utils.ColorRGB

	InitialPosition components.Position
	TargetPosition  components.Position

	Actor    *components.Actor
	AI       *components.AI
	Combat   *components.Combat
	Item     *components.Item
	Mutation *components.Mutation
	Position *components.Position

	Blocks bool
	Dead   bool

	FoV fov.FoVMap

	VisibilityRange int

	Mutations components.Mutations
	Inventory *Inventory
}

// NewEntity creates a new unique entity.
func NewEntity(name string, char string, color utils.ColorRGB, initPosition components.Position, blocks bool) *Entity {
	return &Entity{
		Name:            name,
		Char:            char,
		Color:           color,
		Position:        &initPosition,
		Blocks:          blocks,
		FoV:             make(fov.FoVMap),
		InitialPosition: initPosition,
		Inventory:       &Inventory{MaxSlots: 4},
	}
}

// PickUpItem picks up an item defined in target and adds it to the inventory of e.
func (e *Entity) PickUpItem(target *Entity) (result []ActionResult) {
	if target.Item != nil {
		if target.Item.CanPickup {
			if e.Mutations.Has(components.MutationEffectInventory) {
				err := e.Inventory.Add(target)
				if err == nil {
					target.Position = nil
					result = append(result, ActionResult{Type: ActionResultItemPickedUp})
					result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s picked up %s.", e.Name, target.Name)})
				} else {
					result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("Cannot pick up %s. %s.", target.Name, err.Error())})
				}
			} else {
				result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("You are a cat, you cannot pick up things (or can you?).")})
			}
		}
	}
	return
}

// DropItem drops the target item if it is in the inventory.
func (e *Entity) DropItem(target *Entity) (result []ActionResult) {
	if target.Item != nil {
		if item := e.Inventory.PopOneByName(target.Name); item != nil {
			item.Position = &components.Position{X: e.Position.X, Y: e.Position.Y}
			result = append(result, ActionResult{Type: ActionResultItemDropped})
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("Item %s dropped.", target.Name)})
		}
	}
	return
}

// UseItem uses a item.
func (e *Entity) UseItem(target *Entity) (result []ActionResult) {
	if target.Item != nil {
		if target.Item.Consumable {
			var effectString string
			switch target.Item.Effect {
			case components.ItemEffectHealing:
				e.Combat.CurrentHP += target.Item.Data
				if e.Combat.CurrentHP > e.Combat.HP {
					e.Combat.CurrentHP = e.Combat.HP
				}
				effectString = fmt.Sprintf("%d healed.", target.Item.Data)
			default:
				log.Printf("Effect not implemented")
			}
			result = append(result, ActionResult{Type: ActionResultItemUsed})
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s consumed. %s", target.Name, effectString)})
		} else {
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s is not consumable.", target.Name)})
		}
	}
	return
}

// ConsumeMutation takes the mutation defined in target and adds it to e.
func (e *Entity) ConsumeMutation(target *Entity) (result []ActionResult) {
	if target.Mutation != nil {
		if !e.Mutations.Has(target.Mutation.Effect) {
			e.Mutations = append(e.Mutations, *target.Mutation)
			target.Position = nil
			result = append(result, ActionResult{Type: ActionResultMutationConsumed, MutationEffectValue: target.Mutation.Effect})
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s gained mutation %s.", e.Name, target.Name)})
		} else {
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s already has %s.", e.Name, target.Name)})
		}
	}

	return
}

// MoveTo moves the entity to (y,y).
func (e *Entity) MoveTo(p components.Position) {
	e.Position.X = p.X
	e.Position.Y = p.Y
}

// Attack the target entity.
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
