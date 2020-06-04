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
)

func (d MutationType) String() string {
	return [...]string{"Unknown", "Initial", "Inventory"}[d]
}

type Mutations []MutationType

func (m Mutations) Has(mutation MutationType) bool {
	for _, m := range m {
		if m == mutation {
			return true
		}
	}
	return false
}

func MutationFromString(mutationString string) (MutationType, error) {
	switch strings.ToLower(mutationString) {
	case "inventory":
		return MutationInventory, nil
	default:
		return MutationUnknown, fmt.Errorf("Unknown mutation '%s'", mutationString)
	}
}
