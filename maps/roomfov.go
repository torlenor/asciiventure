package maps

import "math"

// This implements a simple ray-casting FoV method.

// UpdateFoV updates the map with current field of view data based on the current player postion.
// viewRange is the number of tiles the player can see.
func (r *Room) UpdateFoV(viewRange, playerX, playerY int32) {
	r.clearVisible()
	for i := 0; i < 360; i++ {
		x := math.Cos(float64(i) * 0.01745)
		y := math.Sin(float64(i) * 0.01745)
		r.doFoV(viewRange, x, y, playerX, playerY)
	}
}

// doFoV performs the actual Field of View calculation for the given view range and player coordinates
// in the direction of the provided unit vector (x,y).
func (r *Room) doFoV(viewRange int32, x float64, y float64, playerX, playerY int32) {
	ox := float64(playerX)
	oy := float64(playerY)
	for i := int32(0); i < viewRange; i++ {
		ix := int32(ox + 0.5)
		iy := int32(oy + 0.5)
		if y, ok := r.Tiles[iy]; ok {
			if t, ok := y[ix]; ok {
				t.Seen = true
				t.Visible = true
				r.Tiles[iy][ix] = t
				if t.Opaque {
					return
				}
			}
		}
		ox += float64(x)
		oy += float64(y)
	}
}

// ClearSeen removes the seen state on all tiles of the room.
func (r *Room) ClearSeen() {
	for i, y := range r.Tiles {
		for j, t := range y {
			t.Seen = false
			r.Tiles[i][j] = t
		}
	}
}

// clearVisible removes the visible state on all tiles of the room.
func (r *Room) clearVisible() {
	for i, y := range r.Tiles {
		for j, t := range y {
			t.Visible = false
			r.Tiles[i][j] = t
		}
	}
}
