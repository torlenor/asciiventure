package entity

import (
	"fmt"
	"log"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/utils"
)

// Entity defines an entity in our entity-component-system system
type Entity struct {
	TargetPosition utils.Vec2

	Actor      *components.Actor
	AI         *components.AI
	Appearance *components.Appearance
	Combat     *components.Combat
	Health     *components.Health
	IsBlocking *components.IsBlocking
	IsDead     *components.IsDead
	Item       *components.Item
	Mutagen    *components.Mutation
	Name       string // Every entity has a name, even when it's empty
	Position   *components.Position
	Vision     *components.Vision

	FoV fov.FoVMap

	Inventory *Inventory
	Mutations components.Mutations
}

// NewEntity creates a new unique entity.
func NewEntity(name string, appearance *components.Appearance, initPosition utils.Vec2, blocks bool) *Entity {
	e := &Entity{
		Name:       name,
		Appearance: appearance,
		Position:   &components.Position{Current: initPosition, Initial: initPosition},
		FoV:        make(fov.FoVMap),
		Inventory:  &Inventory{MaxSlots: 4},
	}
	if blocks {
		e.IsBlocking = &components.IsBlocking{}
	}
	return e
}

// NewEmptyEntity creates an empty entity.
func NewEmptyEntity() *Entity {
	return &Entity{
		FoV: make(fov.FoVMap),
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
			item.Position = &components.Position{Current: e.Position.Current}
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
				if e.Health != nil {
					e.Health.CurrentHP += target.Item.Data
					if e.Health.CurrentHP > e.Health.HP {
						e.Health.CurrentHP = e.Health.HP
					}
					effectString = fmt.Sprintf("%d healed.", target.Item.Data)
				} else {
					return
				}
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
	if target.Mutagen != nil {
		if !e.Mutations.Has(target.Mutagen.Effect) {
			e.Mutations = append(e.Mutations, *target.Mutagen)
			target.Position = nil
			result = append(result, ActionResult{Type: ActionResultMutationConsumed, MutationEffectValue: target.Mutagen.Effect})
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s gained mutation %s.", e.Name, target.Name)})
		} else {
			result = append(result, ActionResult{Type: ActionResultMessage, StringValue: fmt.Sprintf("%s already has %s.", e.Name, target.Name)})
		}
	}

	return
}

// MoveTo moves the entity to (y,y).
func (e *Entity) MoveTo(p utils.Vec2) {
	if e.Position != nil {
		e.Position.Current = p
	} else {
		e.Position = &components.Position{Current: p, Initial: p}
	}
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
