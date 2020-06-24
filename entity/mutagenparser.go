package entity

import (
	"log"
)

// ParseMutagen parses a mutagen description and returns the corresponding entity.
func ParseMutagen(filename string) *Entity {
	e, err := ParseJSON(filename)
	if err != nil || e == nil {
		log.Printf("%s", err)
		return nil
	}

	if e.Mutagen == nil {
		log.Printf("Not a mutagen entity file")
		return nil
	}
	e.Blocks = true

	return e
}
