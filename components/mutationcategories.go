package components

import (
	"encoding/json"
	"fmt"
	"strings"
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
