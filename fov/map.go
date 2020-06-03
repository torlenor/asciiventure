package fov

import "github.com/torlenor/asciiventure/components"

type FoV struct {
	Visible bool
	Seen    bool
}

type FoVMap map[int32]map[int32]FoV

func NewFovMap() FoVMap {
	return make(FoVMap)
}

func (m FoVMap) Seen(p components.Position) bool {
	if fov, ok := m[p.Y][p.X]; ok {
		return fov.Seen
	}
	return false
}

func (m FoVMap) Visible(p components.Position) bool {
	if fov, ok := m[p.Y][p.X]; ok {
		return fov.Visible
	}
	return false
}

func (m FoVMap) UpdateSeen(p components.Position, seen bool) {
	if _, ok := m[p.Y]; !ok {
		m[p.Y] = make(map[int32]FoV)
	}
	m[p.Y][p.X] = FoV{Seen: seen, Visible: m[p.Y][p.X].Visible}
}

func (m FoVMap) UpdateVisible(p components.Position, visible bool) {
	if _, ok := m[p.Y]; !ok {
		m[p.Y] = make(map[int32]FoV)
	}
	m[p.Y][p.X] = FoV{Seen: m[p.Y][p.X].Seen, Visible: visible}
}

// ClearSeen removes the seen state on all tiles of the room.
func (m FoVMap) ClearSeen() {
	for i, y := range m {
		for j, t := range y {
			t.Seen = false
			m[i][j] = t
		}
	}
}
