package gamemap

type rect struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func newRect(x, y, w, h int) rect {
	return rect{x1: x, y1: y, x2: x + w, y2: y + h}
}

// center returns the center of the room
func (r rect) center() (centerX, centerY int) {
	centerX = (r.x1 + r.x2) / 2
	centerY = (r.y1 + r.y2) / 2
	return
}

func (r rect) intersect(other rect) bool {
	return (r.x1 <= other.x2 && r.x2 >= other.x1 &&
		r.y1 <= other.y2 && r.y2 >= other.y1)
}
