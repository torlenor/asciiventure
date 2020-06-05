package components

import (
	"fmt"
	"strings"
)

type MutationType int

const (
	MutationUnknown MutationType = iota
	MutationInitial
	MutationInventory
	MutationXRay
	MutationIncreasedVision
)

type MutationCategory int

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

type Mutation struct {
	Type     MutationType
	Category MutationCategory
}

type Mutations []Mutation

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
