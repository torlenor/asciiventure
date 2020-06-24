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
	MutationEffectIncreasedVision
)

func (d MutationEffect) String() string {
	return [...]string{"Unknown", "Inventory", "XRay", "IncreasedVision"}[d]
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

// MutationCategory type
type MutationCategory int

// Available MutationCategories
const (
	MutationCategoryUnknown MutationCategory = iota
	MutationCategoryCore
	MutationCategoryEyes
	MutationCategoryClaws
	MutationCategoryTail
)

func (d MutationCategory) String() string {
	return [...]string{"Unknown", "Core", "Eyes", "Claws", "Tail"}[d]
}

// MutationCategoryFromString returns a MutationCategory from the provided string.
func MutationCategoryFromString(mutationCategoryString string) (MutationCategory, error) {
	switch strings.ToLower(mutationCategoryString) {
	case "core":
		return MutationCategoryCore, nil
	case "eyes":
		return MutationCategoryEyes, nil
	case "claws":
		return MutationCategoryClaws, nil
	case "tail":
		return MutationCategoryTail, nil
	default:
		return MutationCategoryUnknown, fmt.Errorf("Unknown mutation category '%s'", mutationCategoryString)
	}
}

// UnmarshalJSON unmarshals a JSON into a MutationCategory.
func (d *MutationCategory) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var ok bool
	categoryStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("Category not defined or not string")
	}

	var err error
	category, err := MutationCategoryFromString(categoryStr)
	if err != nil {
		return err
	}

	*d = category

	return nil
}

// Mutation describes properties of one mutation entity.
type Mutation struct {
	Effect   MutationEffect   `json:"Effect"`
	Category MutationCategory `json:"Category"`
	Data     int32            `json:"Data"`
}

func (m Mutation) String() string {
	return fmt.Sprintf("[%s] %s", m.Category, m.Effect)
}

// Mutations holds a list of Mutations
type Mutations []Mutation

// IsCategory returns true if the given Mutation is of the specified category.
func (m Mutation) IsCategory(category MutationCategory) bool {
	if m.Category == category {
		return true
	}
	return false
}

// Has returns true if the Mutations list has a mutation with that effect.
func (m Mutations) Has(mutation MutationEffect) bool {
	for _, m := range m {
		if m.Effect == mutation {
			return true
		}
	}
	return false
}

// GetData returns the data for the specified MutationEffect, or always 0 if it does not exist.
// Do not forget to check first with Has(mutation)!
func (m Mutations) GetData(mutation MutationEffect) int32 {
	for _, m := range m {
		if m.Effect == mutation {
			return m.Data
		}
	}
	return 0
}

// Get returns the mutation or nil or an error if it does not exist.
func (m Mutations) Get(mutation MutationEffect) (Mutation, error) {
	for _, m := range m {
		if m.Effect == mutation {
			return m, nil
		}
	}
	return Mutation{}, fmt.Errorf("%s does not exist", mutation)
}

// GetDescription returns the description for the mutation effect.
func (m Mutation) GetDescription() string {
	switch m.Effect {
	case MutationEffectInventory:
		return fmt.Sprintf("Provides an inventory.")
	case MutationEffectXRay:
		return fmt.Sprintf("Lets you look through walls.")
	case MutationEffectIncreasedVision:
		return fmt.Sprintf("Permanently increases vision by %d.", m.Data)
	default:
		return "Unknown"
	}
}
