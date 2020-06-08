package fov

import (
	"math"

	"github.com/torlenor/asciiventure/components"
)

// TODO: Raycasting FoV has artifacts in low angle situations

// OpaqueGraph is a interface describing a graph which provides a Opaque function.
type OpaqueGraph interface {
	Opaque(p components.Position) bool
}

// UpdateFoV updates the map with current field of view data based on the provided entity postion.
// viewRange is the number of tiles the entity can see.
func UpdateFoV(r OpaqueGraph, fovMap FoVMap, viewRange int, entityPosition components.Position, ignoreOpaque bool) {
	fovMap.ClearVisible()
	for i := 0; i < 360; i++ {
		uvecX := math.Cos(float64(i) * 0.01745)
		uvecY := math.Sin(float64(i) * 0.01745)
		doFoV(r, fovMap, viewRange, uvecX, uvecY, entityPosition, ignoreOpaque)
	}
}

// doFoV performs the actual Field of View calculation for the given view range and player coordinates
// in the direction of the provided unit vector (x,y).
func doFoV(r OpaqueGraph, fovMap FoVMap, viewRange int, x float64, y float64, entityPosition components.Position, ignoreOpaque bool) {
	ox := float64(entityPosition.X)
	oy := float64(entityPosition.Y)
	for i := 0; i < viewRange; i++ {
		ix := int(ox + 0.5)
		iy := int(oy + 0.5)
		fovMap.UpdateSeen(components.Position{X: ix, Y: iy}, true)
		fovMap.UpdateVisible(components.Position{X: ix, Y: iy}, true)
		if r.Opaque(components.Position{X: ix, Y: iy}) && !ignoreOpaque {
			return
		}
		ox += float64(x)
		oy += float64(y)
	}
}
