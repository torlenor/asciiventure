package console

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/torlenor/asciiventure/renderers"
)

type charData struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	Width   int `json:"width"`
	Height  int `json:"height"`
	OriginX int `json:"originX"`
	OriginY int `json:"originY"`
	Advance int `json:"advance"`
}

type charSetData struct {
	Name       string          `json:"name"`
	Size       int             `json:"size"`
	Bold       bool            `json:"bold"`
	Italic     bool            `json:"italic"`
	Width      int             `json:"width"`
	Height     int             `json:"height"`
	Characters map[string]Char `json:"characters"`
}

// NewFontTilesetFromJSON generates a new FontTileSet from the provided
// image and description file.
func NewFontTilesetFromJSON(renderer *renderers.Renderer, imagePath string, descriptionPath string) (*FontTileSet, error) {
	jsonFile, err := os.Open("./assets/textures/ascii_ext_courier.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to open description file for glyph texture: %s", err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read description file for glyph texture: %s", err)
	}
	var data charSetData
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse description file for glyph texture: %s", err)
	}
	texture, err := createTextureFromFile(renderer, imagePath)
	if err != nil {
		return nil, fmt.Errorf("Error creating font: %s", err)
	}

	font := FontTileSet{
		t: texture,
	}

	font.characters = make(map[string]Char)

	maxWidth := int32(0)
	maxHeight := int32(0)
	for char, chardata := range data.Characters {
		font.characters[char] = Char{X: chardata.X, Y: chardata.Y, Width: chardata.Width, Height: chardata.Height}
		if chardata.Width > maxWidth {
			maxWidth = chardata.Width
		}
		if chardata.Height > maxHeight {
			maxHeight = chardata.Height
		}
	}

	font.charWidth = maxWidth
	font.charHeight = maxHeight

	return &font, nil
}
