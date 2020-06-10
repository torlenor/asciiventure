package gamemap

import "github.com/torlenor/asciiventure/utils"

// Tile is one segment on a game map
type Tile struct {
	Char string

	ForegroundColor utils.ColorRGB
	BackgroundColor utils.ColorRGBA

	Opaque bool
}
