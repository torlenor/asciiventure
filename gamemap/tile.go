package gamemap

import "github.com/torlenor/asciiventure/components"

// Tile is one segment on a game map
type Tile struct {
	Char string

	ForegroundColor components.ColorRGB
	BackgroundColor components.ColorRGBA

	Opaque bool
}
