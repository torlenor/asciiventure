package entity

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/torlenor/asciiventure/components"
)

// ParseMutagen parses a mutagen description and returns the corresponding entity.
func ParseMutagen(filename string) *Entity {
	file, _ := ioutil.ReadFile(filename)
	data := EntityData{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Printf("Error parsing mutagen file %s: %s", filename, err)
	}

	color := components.ColorRGB{R: data.Glyph.Color.R, G: data.Glyph.Color.G, B: data.Glyph.Color.B}
	e := NewEntity(data.Name, data.Glyph.Char, color, components.Position{}, true)
	e.Item = &components.Item{CanPickup: data.Item.CanPickup, Consumable: data.Item.Consumable}
	if data.Mutagen.IsMutagen {
		if t, err := components.MutationTypeFromString(data.Mutagen.Type); err == nil {
			if c, err := components.MutationCategoryFromString(data.Mutagen.Category); err == nil {
				e.Mutations = append(e.Mutations, components.Mutation{Type: t, Category: c, Activatable: data.Mutagen.Activatable, CooldownTurns: data.Mutagen.CooldownTurns, ActiveTurns: data.Mutagen.ActiveTurns, Ready: true})
			} else {
				log.Printf("%s", err)
				return nil
			}
		} else {
			log.Printf("%s", err)
			return nil
		}
	} else {
		log.Printf("Not a mutagen")
		return nil
	}

	return e
}
