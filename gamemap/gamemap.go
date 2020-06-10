package gamemap

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"strings"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/renderers"
	"github.com/torlenor/asciiventure/utils"
)

const (
	emptyChar = "·"
)

var foregroundColorEmptyDot = utils.ColorRGB{R: 220, G: 220, B: 220}
var foregroundColorWallVisible = utils.ColorRGB{R: 200, G: 200, B: 200}
var foregroundColorNotVisible = utils.ColorRGB{R: 80, G: 80, B: 100}
var foregroundColorEmptyDotNotVisible = utils.ColorRGB{R: 40, G: 40, B: 40}

// GameMap holds the data of a game map
type GameMap struct {
	T     *assets.GlyphTexture
	Tiles map[int]map[int]Tile

	SpawnPoint components.Position

	notSeenGlyph renderers.Glyph
}

// NewGameMapFromString constructs a room from the provided room description string
func NewGameMapFromString(s string, glyphTexture *assets.GlyphTexture) (GameMap, error) {
	r := strings.NewReader(s)
	return NewGameMapFromReader(r, glyphTexture)
}

// NewGameMapFromReader constructs a room where the room description is read from the provided Reader
func NewGameMapFromReader(r io.Reader, glyphTexture *assets.GlyphTexture) (GameMap, error) {
	b := bufio.NewReader(r)
	lines := []string{}
	for l, _, err := b.ReadLine(); err == nil; l, _, err = b.ReadLine() {
		lines = append(lines, string(l))
	}
	if len(lines) == 0 {
		return GameMap{}, fmt.Errorf("Not a valid room description")
	}
	room := GameMap{
		Tiles: make(map[int]map[int]Tile),
	}
	spawnPointSet := false
	for y, l := range lines {
		room.Tiles[int(y)] = make(map[int]Tile)
		cntX := -1
		for _, r := range l {
			cntX++
			x := int(cntX)
			c := string(r)
			if c == "@" {
				if spawnPointSet {
					log.Printf("Warning: Player spawn point defined more than once")
					c = " "
				}
				room.SpawnPoint = components.Position{X: int(x), Y: int(y)}
				spawnPointSet = true
				c = " "
			}

			var foregroundColor utils.ColorRGB
			if c == " " {
				c = "·"
				foregroundColor = foregroundColorEmptyDot
			} else {
				foregroundColor = foregroundColorWallVisible
			}
			opaque := false
			if c == "#" {
				opaque = true
			}
			room.Tiles[int(y)][int(x)] = Tile{Char: c, Opaque: opaque, ForegroundColor: foregroundColor}
		}
	}

	room.T = glyphTexture

	room.notSeenGlyph, _ = room.T.Get("#")
	room.notSeenGlyph.Color = utils.ColorRGB{
		R: 20,
		G: 20,
		B: 20,
	}

	return room, nil
}

// Distance returns the distance between two points on the map.
func (r GameMap) Distance(a components.Position, b components.Position) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

// Empty returns true if the specified coordinates of the room are empty
// and inside the map boundaries.
func (r *GameMap) Empty(x, y int) bool {
	maxx, maxy := r.Dimensions()
	if x > maxx || x < 0 || y > maxy || y < 0 {
		return false
	}

	if y, ok := r.Tiles[y]; ok {
		if x, ok := y[x]; ok {
			if x.Char != emptyChar {
				return false
			}
		}
	}
	return true
}

// Opaque returns true if the specified position is not transparent.
func (r *GameMap) Opaque(p components.Position) bool {
	return r.Tiles[p.Y][p.X].Opaque
}

// Dimensions returns the max width and height of the room.
func (r *GameMap) Dimensions() (int, int) {
	maxx := int(0)
	maxy := int(0)
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

// InDimensions returns true if the specified position is inside the map dimensions
func (r *GameMap) InDimensions(p components.Position) bool {
	maxx, maxy := r.Dimensions()
	if p.X > maxx || p.X < 0 || p.Y > maxy || p.Y < 0 {
		return false
	}
	return true
}

// Neighbors returns the empty neighbors for a given point
func (r *GameMap) Neighbors(p components.Position) []components.Position {
	var neighbors []components.Position
	for x := int(-1); x <= 1; x++ {
		for y := int(-1); y <= 1; y++ {
			if r.Empty(p.X+x, p.Y+y) {
				neighbors = append(neighbors, components.Position{X: p.X + x, Y: p.Y + y})
			}
		}
	}

	return neighbors
}
