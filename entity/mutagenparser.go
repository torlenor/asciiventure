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
	if data.Mutagen.IsMutagen {
		if t, err := components.MutationEffectFromString(data.Mutagen.Type); err == nil {
			if c, err := components.MutationCategoryFromString(data.Mutagen.Category); err == nil {
				e.Mutation = &components.Mutation{Effect: t, Category: c, Data: data.Mutagen.Data}
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
