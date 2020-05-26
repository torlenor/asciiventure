package maps

import (
	"log"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/renderers"
)

// Render renders the current state of the room to the provided renderer.
func (r *Room) Render(renderer *renderers.Renderer, offsetX, offsetY int32) {
	// TODO: We can probably implement a fancy way to render a background texture with "." and over that only the things which are visible/seen
	for y, l := range r.Tiles {
		for x, t := range l {
			if g, ok := r.T.Get(t.Char); ok {

				if t.Visible {
					g.Color = r.Colors[y][x]
				} else if t.Seen {
					g.Color = components.ColorRGB{R: 50, G: 100, B: 50}
				} else {
					g = r.notSeenGlyph
				}

				renderer.RenderGlyph(g, x, y)
			} else {
				log.Printf("Error rendering map glyph at %dx%d: Glyph for char %s not found", x, y, t.Char)
			}
		}
	}
}
