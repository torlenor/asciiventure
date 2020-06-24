package components

import "github.com/torlenor/asciiventure/utils"

// Renderable holds all data needed to display the entity on the screen.
type Renderable struct {
	Char  string          `json:"Char"`
	Color utils.ColorRGBA `json:"Color"`
}
