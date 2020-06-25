package gamemap

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"strings"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/renderers"
	"github.com/torlenor/asciiventure/utils"
)

const (
	emptyChar = "·"
)

var foregroundColorEmptyDot = utils.ColorRGBA{R: 220, G: 220, B: 220, A: 100}
var foregroundColorWallVisible = utils.ColorRGBA{R: 200, G: 200, B: 200, A: 255}
var foregroundColorNotVisible = utils.ColorRGBA{R: 80, G: 80, B: 100, A: 255}
var foregroundColorEmptyDotNotVisible = utils.ColorRGBA{R: 40, G: 40, B: 40, A: 255}

// GameMap holds the data of a game map
type GameMap struct {
	T     *assets.GlyphTexture
	Tiles map[int]map[int]Tile

	Entities *[]*entity.Entity

	SpawnPoint     utils.Vec2
	MapChangePoint utils.Vec2

	notSeenGlyph renderers.Glyph

	currentOffsetX int32
	currentOffsetY int32
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
	mapChangePointSet := false
	for y, l := range lines {
		room.Tiles[int(y)] = make(map[int]Tile)
		cntX := -1
		for _, r := range l {
			cntX++
			x := int32(cntX)
			c := string(r)
			if c == "@" {
				if spawnPointSet {
					log.Printf("Warning: Player spawn point defined more than once")
				}
				room.SpawnPoint = utils.Vec2{X: x, Y: int32(y)}
				spawnPointSet = true
				c = " "
			}
			if c == "+" {
				if mapChangePointSet {
					log.Printf("Warning: Map change point defined more than once")
				}
				room.MapChangePoint = utils.Vec2{X: x, Y: int32(y)}
				mapChangePointSet = true
				c = " "
			}

			var foregroundColor utils.ColorRGBA
			if c == " " {
				c = "·"
				foregroundColor = foregroundColorEmptyDot
			} else {
				foregroundColor = foregroundColorWallVisible
			}
			opaque := false
			blocking := false
			if c == "#" {
				opaque = true
				blocking = true
			}
			room.Tiles[int(y)][int(x)] = Tile{Char: c, Opaque: opaque, Blocking: blocking, ForegroundColor: foregroundColor}
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
func (r GameMap) Distance(a utils.Vec2, b utils.Vec2) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

// Empty returns true if the specified coordinates of the room are empty
// and inside the map boundaries.
func (r *GameMap) Empty(p utils.Vec2) bool {
	maxx, maxy := r.Dimensions()
	if p.X > maxx || p.X < 0 || p.Y > maxy || p.Y < 0 {
		return false
	}

	if y, ok := r.Tiles[int(p.Y)]; ok {
		if x, ok := y[int(p.X)]; ok {
			if x.Blocking {
				return false
			}
		}
	}
	return true
}

// Opaque returns true if the specified position is not transparent.
func (r *GameMap) Opaque(p utils.Vec2) bool {
	return r.Tiles[int(p.Y)][int(p.X)].Opaque
}

// Dimensions returns the max width and height of the room.
func (r *GameMap) Dimensions() (int32, int32) {
	maxx := int32(0)
	maxy := int32(0)
	for y, r := range r.Tiles {
		if int32(y) > maxy {
			maxy = int32(y)
		}
		for x := range r {
			if int32(x) > maxx {
				maxx = int32(x)
			}
		}
	}
	return maxx, maxy
}

// InDimensions returns true if the specified position is inside the map dimensions
func (r *GameMap) InDimensions(p utils.Vec2) bool {
	maxx, maxy := r.Dimensions()
	if p.X > maxx || p.X < 0 || p.Y > maxy || p.Y < 0 {
		return false
	}
	return true
}

// Neighbors returns the empty neighbors for a given point
func (r *GameMap) Neighbors(p utils.Vec2) []utils.Vec2 {
	var neighbors []utils.Vec2
	for x := int32(-1); x <= 1; x++ {
		for y := int32(-1); y <= 1; y++ {
			n := p.Add(utils.Vec2{X: x, Y: y})
			if r.Empty(n) {
				neighbors = append(neighbors, n)
			}
		}
	}

	return neighbors
}

// IsPortal returns true if the location is a portal to another map.
func (r *GameMap) IsPortal(p utils.Vec2) bool {
	if p.Equal(r.MapChangePoint) {
		return true
	}
	return false
}

// GetPositionFromRenderCoordinates returns the position on the game map for the provided tile position.
// Returns (-1, -1) if the position is outside of the map.
func (r *GameMap) GetPositionFromRenderCoordinates(x, y int32) (ex, ey int32) {
	ex = x - r.currentOffsetX
	ey = y - r.currentOffsetY

	maxX, maxY := r.Dimensions()
	if ex > int32(maxX) || ey > int32(maxY) {
		return -1, -1
	}

	return
}

// GetRenderCoordinatesFromPosition returns the render position for a position on the game map.
func (r *GameMap) GetRenderCoordinatesFromPosition(ex, ey int32) (x, y int32) {
	x = ex + r.currentOffsetX
	y = ey + r.currentOffsetY

	return
}
