package components

import (
	"fmt"
	"strings"
)

// MutationType type
type MutationType int

// Available MutationTypes
const (
	MutationUnknown MutationType = iota
	MutationInitial
	MutationInventory
	MutationXRay
	MutationIncreasedVision
)

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

func (d MutationType) String() string {
	return [...]string{"Unknown", "Initial", "Inventory", "XRay", "IncreasedVision"}[d]
}

func (d MutationCategory) String() string {
	return [...]string{"Unknown", "Core", "Eyes", "Claws", "Tail"}[d]
}

// Mutation describes properties of one mutation entity.
type Mutation struct {
	Type           MutationType
	Category       MutationCategory
	Activatable    bool
	Ready          bool
	Active         bool
	CooldownTurns  int
	ActiveTurns    int
	TurnsRemaining int
}

func (m Mutation) String() string {
	if m.Activatable {
		if m.Ready {
			return fmt.Sprintf("[%s] %s (Ready)", m.Category, m.Type)
		} else if m.Active {
			return fmt.Sprintf("[%s] %s (Active, %d turns remaining)", m.Category, m.Type, m.TurnsRemaining)
		}
		// on cooldown
		return fmt.Sprintf("[%s] %s (On cooldown, %d turns remaining)", m.Category, m.Type, m.TurnsRemaining)
	}
	// always active
	return fmt.Sprintf("[%s] %s", m.Category, m.Type)
}

// Mutations holds a list of Mutations
type Mutations []Mutation

func (m Mutation) IsCategory(category MutationCategory) bool {
	if m.Category == category {
		return true
	}
	return false
}

func (m Mutations) Has(mutation MutationType) bool {
	for _, m := range m {
		if m.Type == mutation {
			return true
		}
	}
	return false
}

func MutationTypeFromString(mutationString string) (MutationType, error) {
	switch strings.ToLower(mutationString) {
	case "inventory":
		return MutationInventory, nil
	case "xray":
		return MutationXRay, nil
	case "increasedvision":
		return MutationIncreasedVision, nil
	default:
		return MutationUnknown, fmt.Errorf("Unknown mutation '%s'", mutationString)
	}
}

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
