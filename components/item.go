package components

import (
	"encoding/json"
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
	CanPickup bool `json:"CanPickup"`

	Consumable bool       `json:"Consumable"`
	Effect     ItemEffect `json:"Effect"`

	Data int32 `json:"Data"`
}

// UnmarshalJSON unmarshals a JSON into a ItemEffect.
func (d *ItemEffect) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var ok bool
	effectStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("Effect not defined or not string")
	}

	var err error
	effect, err := ItemEffectFromString(effectStr)
	if err != nil {
		return err
	}

	*d = effect

	return nil
}
