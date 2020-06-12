package components

import (
	"fmt"
	"strings"
)

// ItemEffect type
type ItemEffect int

// Available ItemEffects
const (
	ItemEffectUnknown ItemEffect = iota
	ItemEffectHealing
)

func (d ItemEffect) String() string {
	return [...]string{"Unknown", "Healing"}[d]
}

// ItemEffectFromString returns a ItemEffect from the provided string
func ItemEffectFromString(itemString string) (ItemEffect, error) {
	switch strings.ToLower(itemString) {
	case "healing":
		return ItemEffectHealing, nil
	default:
		return ItemEffectUnknown, fmt.Errorf("Unknown item '%s'", itemString)
	}
}

// Item is used in an entity to mark it as an item.
type Item struct {
	CanPickup bool

	Consumable bool
	Effect     ItemEffect

	Data int
}
