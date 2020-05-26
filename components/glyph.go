package components

import (
	"github.com/torlenor/asciiventure/entities"
	"github.com/veandco/go-sdl2/sdl"
)

// Glyph represents one rendereable glyph
type Glyph struct {
	T   *sdl.Texture
	Src *sdl.Rect

	Color  ColorRGB
	Shadow bool

	Width  int32
	Height int32

	OffsetX int32
	OffsetY int32
}

type GlyphManager map[entities.Entity]Glyph
