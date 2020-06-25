package utils

import "fmt"

// Vec2 is a 2d int32 vector.
type Vec2 struct {
	X int32
	Y int32
}

// String returns the stringified vector in the format (x,y).
func (p Vec2) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Equal compares the vector with e and returns true if equal.
func (p Vec2) Equal(e Vec2) bool {
	return p.X == e.X && p.Y == e.Y
}

// Add a vector and return a new vector.
func (p *Vec2) Add(e Vec2) Vec2 {
	n := Vec2{
		X: p.X + e.X,
		Y: p.Y + e.Y,
	}
	return n
}
