package gamemap

import (
	"math/rand"

	"github.com/torlenor/asciiventure/assets"
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/utils"
)

//NewRandomMap returns a random game map with the specified number of rooms and sizes.
func NewRandomMap(maxRooms int, roomMinSize, roomMaxSize, mapWidth, mapHeight int, glyphTexture *assets.GlyphTexture) GameMap {
	var gameMap GameMap
	gameMap.Tiles = make(map[int]map[int]Tile)

	for y := int(0); y < mapHeight; y++ {
		if _, ok := gameMap.Tiles[y]; !ok {
			gameMap.Tiles[int(y)] = make(map[int]Tile)
		}
		for x := int(0); x < mapWidth; x++ {
			foregroundColor := components.ColorRGB{
				R: 160, G: 160, B: 160}
			gameMap.Tiles[int(y)][int(x)] = Tile{Char: "#", Opaque: true, ForegroundColor: foregroundColor}
		}
	}

	var rooms []rect

Loop:
	for i := 0; i < maxRooms; i++ {
		w := rand.Intn(int(roomMaxSize)+1) + int(roomMinSize) + 1
		h := rand.Intn(int(roomMaxSize)+1) + int(roomMinSize) + 1

		x := int(rand.Intn(int(mapWidth) - w))
		y := int(rand.Intn(int(mapHeight) - h))

		newRoom := newRect(int(x), int(y), w, h)
		for _, otherRoom := range rooms {
			if newRoom.intersect(otherRoom) {
				continue Loop
			}
		}

		createRoom(&gameMap, newRoom)
		newX, newY := newRoom.center()

		if i == 0 {
			gameMap.SpawnPoint = components.Position{X: int(newX), Y: int(newY)}
		} else {
			prevX, prevY := rooms[len(rooms)-1].center()
			if rand.Intn(2) == 0 {
				createHTunnel(&gameMap, int(prevX), int(newX), int(prevY))
				createVTunnel(&gameMap, int(prevY), int(newY), int(newX))
			} else {
				createVTunnel(&gameMap, int(prevY), int(newY), int(prevX))
				createHTunnel(&gameMap, int(prevX), int(newX), int(newY))
			}
		}

		rooms = append(rooms, newRoom)
	}

	gameMap.T = glyphTexture

	gameMap.notSeenGlyph, _ = gameMap.T.Get("#")
	gameMap.notSeenGlyph.Color = components.ColorRGB{
		R: 20,
		G: 50,
		B: 20,
	}

	return gameMap
}

func createRoom(gameMap *GameMap, room rect) {
	for y := room.y1 + 1; y < room.y2; y++ {
		if _, ok := gameMap.Tiles[int(y)]; !ok {
			gameMap.Tiles[int(y)] = make(map[int]Tile)
		}
		for x := room.x1 + 1; x < room.x2; x++ {
			foregroundColor := components.ColorRGB{
				R: 60, G: 60, B: 60}
			gameMap.Tiles[int(y)][int(x)] = Tile{Char: "·", Opaque: false, ForegroundColor: foregroundColor}
		}
	}
}

func createHTunnel(gameMap *GameMap, x1, x2, y int) {
	for x := utils.MinInt(x1, x2); x < utils.MaxInt(x1, x2)+1; x++ {
		foregroundColor := components.ColorRGB{
			R: 60, G: 60, B: 60}
		gameMap.Tiles[y][x] = Tile{Char: "·", Opaque: false, ForegroundColor: foregroundColor}
	}
}

func createVTunnel(gameMap *GameMap, y1, y2, x int) {
	for y := utils.MinInt(y1, y2); y < utils.MaxInt(y1, y2)+1; y++ {
		foregroundColor := components.ColorRGB{
			R: 60, G: 60, B: 60}
		gameMap.Tiles[y][x] = Tile{Char: "·", Opaque: false, ForegroundColor: foregroundColor}
	}
}
