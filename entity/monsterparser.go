package entity

import (
	"log"
)

// ParseMonster parses a monster description and returns the corresponding entity.
func ParseMonster(filename string) *Entity {
	e, err := ParseJSON(filename)
	if err != nil || e == nil {
		log.Printf("%s", err)
		return nil
	}

	if e.Combat == nil || e.Health == nil || e.AI == nil || e.Vision == nil {
		log.Printf("Not a monster entity file")
		return nil
	}
	e.Blocks = true

	return e
}
