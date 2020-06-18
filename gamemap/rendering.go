package gamemap

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/console"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/utils"
)

// Render renders the current state of the room to the provided renderer.
func (r *GameMap) Render(console *console.MatrixConsole, foV fov.FoVMap, player *entity.Entity, entities []*entity.Entity, offsetX, offsetY int32) {
	cnx, cny := console.GetDimensions()
	r.currentOffsetX = offsetX - int32(player.Position.X) + cnx/2
	r.currentOffsetY = offsetY - int32(player.Position.Y) + cny/2

	for y, l := range r.Tiles {
		for x, t := range l {
			foregroundColor := t.ForegroundColor
			p := components.Position{X: x, Y: y}
			if !foV.Visible(p) && foV.Seen(p) {
				if t.Char == "·" {
					t.Char = " "
					foregroundColor = foregroundColorEmptyDotNotVisible
				} else {
					foregroundColor = foregroundColorNotVisible
				}
			} else if !foV.Visible(p) {
				continue
			}

			console.PutCharColor(int32(x)+r.currentOffsetX, int32(y)+r.currentOffsetY, t.Char, foregroundColor, utils.ColorRGBA{})
		}
	}

	// TODO: Optimize rendering of entities on map so that we do not need three passes
	for _, e := range entities {
		if e.Position != nil && foV.Visible(*e.Position) && e.Dead {
			console.PutCharColor(int32(e.Position.X)+r.currentOffsetX, int32(e.Position.Y)+r.currentOffsetY, "%", utils.ColorRGBA{R: 150, G: 150, B: 150, A: 255}, utils.ColorRGBA{})
		}
	}
	for _, e := range entities {
		if e.Position != nil && foV.Visible(*e.Position) && (e.Item != nil || e.Mutation != nil) {
			console.PutCharColor(int32(e.Position.X)+r.currentOffsetX, int32(e.Position.Y)+r.currentOffsetY, e.Char, utils.ColorRGBA{R: e.Color.R, G: e.Color.G, B: e.Color.B, A: 255}, utils.ColorRGBA{})
		}
	}
	for _, e := range entities {
		if e.Position == nil || e.Item != nil || e.Mutation != nil || e.Dead {
			continue
		}
		if foV.Visible(*e.Position) {
			console.PutCharColor(int32(e.Position.X)+r.currentOffsetX, int32(e.Position.Y)+r.currentOffsetY, e.Char, utils.ColorRGBA{R: e.Color.R, G: e.Color.G, B: e.Color.B, A: 255}, utils.ColorRGBA{})
		}
	}
}
