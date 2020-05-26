package components

import "github.com/torlenor/asciiventure/entities"

type Position struct {
	X int32
	Y int32
}

type PositionManager map[entities.Entity]Position
