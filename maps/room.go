package maps

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/components"
)

const (
	emptyChar = "·"
)

type SpawnPoint struct {
	X int32
	Y int32
}

type Tile struct {
	X int32
	Y int32

	Char string

	Opaque  bool
	Visible bool
	Seen    bool
}

type Room struct {
	T          *assets.GlyphTexture
	Tiles      map[int32]map[int32]Tile
	Colors     map[int32]map[int32]components.ColorRGB
	SpawnPoint SpawnPoint

	notSeenGlyph components.Glyph
}

// NewRoomFromString constructs a room from the provided room description string
func NewRoomFromString(s string, glyphTexture *assets.GlyphTexture) (Room, error) {
	r := strings.NewReader(s)
	return NewRoom(r, glyphTexture)
}

// NewRoom constructs a room where the room description is read from the provided Reader
func NewRoom(r io.Reader, glyphTexture *assets.GlyphTexture) (Room, error) {
	b := bufio.NewReader(r)
	lines := []string{}
	for l, _, err := b.ReadLine(); err == nil; l, _, err = b.ReadLine() {
		lines = append(lines, string(l))
	}
	if len(lines) == 0 {
		return Room{}, fmt.Errorf("Not a valid room description")
	}
	room := Room{
		Tiles:  make(map[int32]map[int32]Tile),
		Colors: make(map[int32]map[int32]components.ColorRGB),
	}
	spawnPointSet := false
	for y, l := range lines {
		room.Tiles[int32(y)] = make(map[int32]Tile)
		room.Colors[int32(y)] = make(map[int32]components.ColorRGB)
		cntX := -1
		for _, r := range l {
			cntX++
			x := int32(cntX)
			c := string(r)
			if c == "@" {
				if spawnPointSet {
					log.Printf("Warning: Player spawn point defined more than once")
					c = " "
				}
				room.SpawnPoint = SpawnPoint{X: int32(x), Y: int32(y)}
				spawnPointSet = true
				c = " "
			}

			room.Colors[int32(y)][int32(x)] = components.ColorRGB{
				R: 160,
				G: 255,
				B: 160,
			}

			if c == " " {
				c = "·"
				room.Colors[int32(y)][int32(x)] = components.ColorRGB{
					R: 120,
					G: 120,
					B: 120,
				}
			}
			opaque := false
			if c == "#" {
				opaque = true
			}
			room.Tiles[int32(y)][int32(x)] = Tile{X: x, Y: int32(y), Char: c, Opaque: opaque}
		}
	}

	room.T = glyphTexture

	room.notSeenGlyph, _ = room.T.Get("#")
	room.notSeenGlyph.Color = components.ColorRGB{
		R: 20,
		G: 50,
		B: 20,
	}

	return room, nil
}

// Empty returns true if the specified coordinates of the room are empty.
func (r *Room) Empty(x, y int32) bool {
	if y, ok := r.Tiles[y]; ok {
		if x, ok := y[x]; ok {
			if x.Char != emptyChar {
				return false
			}
		}
	}
	return true
}

func (r *Room) Visible(x, y int32) bool {
	return r.Tiles[y][x].Visible
}

// Dimensions returns the max width and height of the room.
func (r *Room) Dimensions() (int32, int32) {
	maxx := int32(0)
	maxy := int32(0)
	for y, r := range r.Tiles {
		if y > maxy {
			maxy = y
		}
		for x := range r {
			if x > maxx {
				maxx = x
			}
		}
	}
	return maxx, maxy
}