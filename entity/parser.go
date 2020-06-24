package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/torlenor/asciiventure/components"
)

type entityData struct {
	Name    string                 `json:"Name"`
	Glyph   *components.Renderable `json:"Glyph"`
	Health  *components.Health     `json:"Health"`
	Combat  *components.Combat     `json:"Combat"`
	AI      *components.AI         `json:"AI"`
	Vision  *components.Vision     `json:"Vision"`
	Item    *components.Item       `json:"Item"`
	Mutagen *components.Mutation   `json:"Mutagen"`
}

// ParseJSON parses a JSON and returns its entity.
func ParseJSON(filename string) (*Entity, error) {
	file, _ := ioutil.ReadFile(filename)
	data := entityData{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		return nil, fmt.Errorf("Error parsing entity JSON file %s: %s", filename, err)
	}

	e := NewEmptyEntity()
	e.Name = data.Name
	e.Renderable = data.Glyph
	e.Health = data.Health
	e.Combat = data.Combat
	e.AI = data.AI
	e.Vision = data.Vision
	e.Item = data.Item
	e.Mutagen = data.Mutagen

	return e, nil
}

// ParseItem parses a item description and returns the corresponding entity.
func ParseItem(filename string) *Entity {
	e, err := ParseJSON(filename)
	if err != nil || e == nil {
		log.Printf("%s", err)
		return nil
	}

	if e.Item == nil {
		log.Printf("Not a item entity file")
		return nil
	}
	e.Blocks = true

	return e
}
