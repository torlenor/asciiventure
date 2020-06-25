package entity

import "log"

// ParseItem parses a item description and returns the corresponding entity.
func ParseItem(filename string) *Entity {
	e, err := ParseJSON(filename)
	if err != nil || e == nil {
		log.Printf("%s", err)
		return nil
	}

	if e.Item == nil || e.Appearance == nil {
		log.Printf("Not a item entity file")
		return nil
	}

	return e
}
