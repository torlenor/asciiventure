package renderers

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/torlenor/asciiventure/utils"
)

// Glyph represents one rendereable glyph
type Glyph struct {
	T   *sdl.Texture
	Src *sdl.Rect

	Color  utils.ColorRGB
	Shadow bool

	Width  int
	Height int

	OffsetX int
	OffsetY int
}
