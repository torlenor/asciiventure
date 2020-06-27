package components

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MutationEffect type
type MutationEffect int

// Available MutationEffects
const (
	MutationEffectUnknown MutationEffect = iota
	MutationEffectInventory
	MutationEffectXRay
	MutationEffectIncreasedVision // increase visibility range
	// TODO: Think of a way to indicate other entities on the map without making them visible, e.g., "hear" enemies
	MutationEffectHeightenedHearing // detect enemies even though you cannot see them
	// TODO: Implement day/night system
	MutationEffectNightVision // increased visibility range at night
	// TODO: Add a system which handles entity stats updates like healthregeneration in a turn
	MutationEffectRegeneration // increased health/whatever regeneration
	// TODO: Think of a way to influence enemy AI
	MutationEffectConfusion // confuse enemies around you
	// TODO: Implement ranged attack
	MutationEffectPyrokinesis   // sets thinks on flames
	MutationEffectPush          // pushes enemies/items away
	MutationEffectTeleport      // teleports you to a random nearby location
	MutationEffectTeleportOther // teleports target to a random nearby location
	// TODO: Make walls destructable
	MutationEffectBurrowingClaws // strengthen your claws and allows you to dig through walls
	// TODO: Create walls dynamically
	MutationEffectForceField // creates a force field which acts like a wall (for you and your enemies)
)

func (d MutationEffect) String() string {
	return [...]string{
		"Unknown",
		"Inventory",
		"XRay",
		"IncreasedVision",
		"HeightenedHearing",
		"NightVision",
		"Regeneration",
		"Confusion",
		"Pyrokinesis",
		"Push",
		"Teleport",
		"TeleportOther",
		"BurrowingClaws",
	}[d]
}

// MutationEffectFromString returns a MutationEffect from the provided string
func MutationEffectFromString(mutationString string) (MutationEffect, error) {
	switch strings.ToLower(mutationString) {
	case "inventory":
		return MutationEffectInventory, nil
	case "xray":
		return MutationEffectXRay, nil
	case "increasedvision":
		return MutationEffectIncreasedVision, nil
	case "heightenedHearing":
		return MutationEffectHeightenedHearing, nil
	case "nightVision":
		return MutationEffectNightVision, nil
	case "regeneration":
		return MutationEffectRegeneration, nil
	case "confusion":
		return MutationEffectConfusion, nil
	case "pyrokinesis":
		return MutationEffectPyrokinesis, nil
	case "push":
		return MutationEffectPush, nil
	case "teleport":
		return MutationEffectTeleport, nil
	case "teleportOther":
		return MutationEffectTeleportOther, nil
	case "burrowingClaws":
		return MutationEffectBurrowingClaws, nil
	default:
		return MutationEffectUnknown, fmt.Errorf("Unknown mutation '%s'", mutationString)
	}
}

// UnmarshalJSON unmarshals a JSON into a MutationEffect.
func (d *MutationEffect) UnmarshalJSON(data []byte) error {
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
	effect, err := MutationEffectFromString(effectStr)
	if err != nil {
		return err
	}

	*d = effect

	return nil
}
