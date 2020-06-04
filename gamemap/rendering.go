package gamemap

import (
	"log"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/renderers"
)

// Render renders the current state of the room to the provided renderer.
func (r *GameMap) Render(renderer *renderers.Renderer, playerFoV fov.FoVMap, offsetX, offsetY int) {
	for y, l := range r.Tiles {
		for x, t := range l {
			if g, ok := r.T.Get(t.Char); ok {
				p := components.Position{X: x, Y: y}
				if playerFoV.Visible(p) {
					g.Color = t.ForegroundColor
				} else if playerFoV.Seen(p) {
					if t.Char == "·" {
						g.Color = foregroundColorEmptyDotNotVisible
					} else {
						g.Color = foregroundColorNotVisible
					}
				} else {
					// g = r.notSeenGlyph
					continue
				}

				renderer.RenderGlyph(g, x, y)
			} else {
				log.Printf("Error rendering map glyph at %dx%d: Glyph for char %s not found", x, y, t.Char)
			}
		}
	}
}
