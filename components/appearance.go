package components

import "github.com/torlenor/asciiventure/utils"

// Appearance holds all data related to the visual appearance of an entity.
type Appearance struct {
	Char  string          `json:"Char"`
	Color utils.ColorRGBA `json:"Color"`
}
