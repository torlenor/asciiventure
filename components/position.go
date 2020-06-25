package components

import (
	"github.com/torlenor/asciiventure/utils"
)

// Position holds the 2d coordinates (x,y) of the current and initial position of an entity.
type Position struct {
	Current utils.Vec2
	Initial utils.Vec2
}
