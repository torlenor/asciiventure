package entity

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/torlenor/asciiventure/components"
)

type MonsterData struct {
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
}

// ParseMonster parses a monster description and returns the corresponding entity.
func ParseMonster(filename string) *Entity {
	file, _ := ioutil.ReadFile(filename)
	data := MonsterData{}

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Printf("Error parsing monster file %s: %s", filename, err)
	}

	color := components.ColorRGB{R: data.Glyph.Color.R, G: data.Glyph.Color.G, B: data.Glyph.Color.B}
	e := NewEntity(data.Name, data.Glyph.Char, color, components.Position{}, true)
	e.Combat = &components.Combat{CurrentHP: data.Combat.HP, HP: data.Combat.HP, Power: data.Combat.Power, Defense: data.Combat.Defense}
	e.AttackRange = data.AI.AttackRange
	e.AttackRangeUntil = data.AI.AttackRangeUntil
	e.VisibilityRange = data.Vision.Range

	return e
}
