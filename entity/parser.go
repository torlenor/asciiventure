package entity

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/torlenor/asciiventure/components"
)

type EntityData struct {
	Name  string `json:"Name"`
	Glyph struct {
		Char  string `json:"Char"`
		Color struct {
			R uint8 `json:"R"`
			G uint8 `json:"G"`
			B uint8 `json:"B"`
		} `json:"Color"`
	} `json:"Glyph"`
	Combat struct {
		HP      int32 `json:"HP"`
		Defense int32 `json:"Defense"`
		Power   int32 `json:"Power"`
	} `json:"Combat"`
	AI struct {
		AttackRange      int `json:"AttackRange"`
		AttackRangeUntil int `json:"AttackRangeUntil"`
	} `json:"AI"`
	Vision struct {
		Range int `json:"Range"`
	} `json:"Vision"`
	Item struct {
		CanPickup   bool   `json:"CanPickup"`
		Consumable  bool   `json:"Consumable"`
		Mutagen     bool   `json:"Mutagen"`
		MutagenType string `json:"MutagenType"`
	} `json:"Item"`
}

// ParseMonster parses a monster description and returns the corresponding entity.
func ParseMonster(filename string) *Entity {
	file, _ := ioutil.ReadFile(filename)
	data := EntityData{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Printf("Error parsing monster file %s: %s", filename, err)
	}

	color := components.ColorRGB{R: data.Glyph.Color.R, G: data.Glyph.Color.G, B: data.Glyph.Color.B}
	e := NewEntity(data.Name, data.Glyph.Char, color, components.Position{}, true)
	e.Combat = &components.Combat{CurrentHP: data.Combat.HP, HP: data.Combat.HP, Power: data.Combat.Power, Defense: data.Combat.Defense}
	e.AI = &components.AI{AttackRange: data.AI.AttackRange, AttackRangeUntil: data.AI.AttackRangeUntil}
	e.VisibilityRange = data.Vision.Range

	return e
}

// ParseItem parses a item description and returns the corresponding entity.
func ParseItem(filename string) *Entity {
	file, _ := ioutil.ReadFile(filename)
	data := EntityData{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Printf("Error parsing item file %s: %s", filename, err)
	}

	color := components.ColorRGB{R: data.Glyph.Color.R, G: data.Glyph.Color.G, B: data.Glyph.Color.B}
	e := NewEntity(data.Name, data.Glyph.Char, color, components.Position{}, true)
	e.Item = &components.Item{CanPickup: data.Item.CanPickup, Consumable: data.Item.Consumable, Mutagen: data.Item.Mutagen}
	if data.Item.Mutagen {
		if m, err := components.MutationFromString(data.Item.MutagenType); err == nil {
			e.Mutations = append(e.Mutations, m)
		} else {
			log.Printf("%s", err)
			return nil
		}
	}

	return e
}
