package components

import "fmt"

// Position holds the 2d coordinates (x,y).
type Position struct {
	X int
	Y int
}

// String returns the stringified position in the format (x,y).
func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Equal compares the position with e and returns true if equal.
func (p Position) Equal(e Position) bool {
	return p.X == e.X && p.Y == e.Y
}
