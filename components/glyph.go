package components

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Glyph represents one rendereable glyph
type Glyph struct {
	T   *sdl.Texture
	Src *sdl.Rect

	Color  ColorRGB
	Shadow bool

	Width  int
	Height int

	OffsetX int
	OffsetY int
}
