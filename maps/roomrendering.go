package maps

import (
	"log"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/renderers"
)

// Render renders the current state of the room to the provided renderer.
func (r *Room) Render(renderer *renderers.Renderer, playerFoV fov.FoVMap, offsetX, offsetY int32) {
	for y, l := range r.Tiles {
		for x, t := range l {
			if g, ok := r.T.Get(t.Char); ok {

				if playerFoV.Visible(components.Position{X: x, Y: y}) {
					g.Color = r.Colors[y][x]
				} else if playerFoV.Seen(components.Position{X: x, Y: y}) {
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
