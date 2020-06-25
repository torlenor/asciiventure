package gamemap

import (
	"github.com/torlenor/asciiventure/console"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/fov"
	"github.com/torlenor/asciiventure/utils"
)

// Render renders the current state of the room to the provided renderer.
func (r *GameMap) Render(console *console.MatrixConsole, foV fov.FoVMap, player *entity.Entity, entities []*entity.Entity, offsetX, offsetY int32) {
	cnx, cny := console.GetDimensions()
	r.currentOffsetX = offsetX - int32(player.Position.Current.X) + cnx/2
	r.currentOffsetY = offsetY - int32(player.Position.Current.Y) + cny/2

	for y, l := range r.Tiles {
		for x, t := range l {
			foregroundColor := t.ForegroundColor
			p := utils.Vec2{X: int32(x), Y: int32(y)}
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
		if e.Position != nil && e.Appearance != nil && foV.Visible(e.Position.Current) && e.IsDead != nil {
			console.PutCharColor(int32(e.Position.Current.X)+r.currentOffsetX, int32(e.Position.Current.Y)+r.currentOffsetY, "%", utils.ColorRGBA{R: 150, G: 150, B: 150, A: 255}, utils.ColorRGBA{})
		}
	}
	for _, e := range entities {
		if e.Position != nil && e.Appearance != nil && foV.Visible(e.Position.Current) && (e.Item != nil || e.Mutagen != nil) {
			console.PutCharColor(int32(e.Position.Current.X)+r.currentOffsetX, int32(e.Position.Current.Y)+r.currentOffsetY, e.Appearance.Char, e.Appearance.Color, utils.ColorRGBA{})
		}
	}
	for _, e := range entities {
		if e.Position == nil || e.Appearance == nil || e.Item != nil || e.Mutagen != nil || e.IsDead != nil {
			continue
		}
		if foV.Visible(e.Position.Current) {
			console.PutCharColor(int32(e.Position.Current.X)+r.currentOffsetX, int32(e.Position.Current.Y)+r.currentOffsetY, e.Appearance.Char, e.Appearance.Color, utils.ColorRGBA{})
		}
	}
}
