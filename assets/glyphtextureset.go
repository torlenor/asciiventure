package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/torlenor/asciiventure/renderers"
	"github.com/torlenor/asciiventure/utils"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Char struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	Width   int `json:"width"`
	Height  int `json:"height"`
	OriginX int `json:"originX"`
	OriginY int `json:"originY"`
	Advance int `json:"advance"`
}

type CharSet struct {
	Name       string          `json:"name"`
	Size       int             `json:"size"`
	Bold       bool            `json:"bold"`
	Italic     bool            `json:"italic"`
	Width      int             `json:"width"`
	Height     int             `json:"height"`
	Characters map[string]Char `json:"characters"`
}

type GlyphTexture struct {
	t            *sdl.Texture
	glyphCharSet CharSet
}

func createTextureFromFile(renderer *renderers.Renderer, imagePath string) (*sdl.Texture, error) {
	image, err := img.Load(imagePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load image file: %s", err)
	}
	defer image.Free()

	glyphTexture, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		return nil, fmt.Errorf("Failed to create texture: %s", err)
	}

	return glyphTexture, nil
}

// NewGlyphTexture generates a new GlyphTexture from the provided
// image and description file.
func NewGlyphTexture(renderer *renderers.Renderer, imagePath string, descriptionPath string) (*GlyphTexture, error) {
	g := &GlyphTexture{}
	jsonFile, err := os.Open("./assets/textures/ascii_ext_courier.json")
	if err != nil {
		return g, fmt.Errorf("Unable to open description file for glyph texture: %s", err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return g, fmt.Errorf("Unable to read description file for glyph texture: %s", err)
	}
	err = json.Unmarshal(byteValue, &g.glyphCharSet)
	if err != nil {
		return g, fmt.Errorf("Unable to parse description file for glyph texture: %s", err)
	}
	g.t, err = createTextureFromFile(renderer, "./assets/textures/ascii_ext_courier.png")
	if err != nil {
		return g, fmt.Errorf("Unable to load image texture: %s", err)
	}

	return g, nil
}

// Get returns a glyph with Dst set to render at origin (0,0).
// Returns true as second value if the operation was successfull.
func (g *GlyphTexture) Get(c string) (renderers.Glyph, bool) {
	if a, ok := g.glyphCharSet.Characters[c]; ok {
		return renderers.Glyph{T: g.t, Src: &sdl.Rect{X: int32(a.X), Y: int32(a.Y), W: int32(a.Width), H: int32(a.Height)},
			OffsetX: a.OriginX, OffsetY: a.OriginY,
			Width: a.Width, Height: a.Height,
			Color: utils.ColorRGB{R: 255, G: 255, B: 255}}, true
	}
	return renderers.Glyph{}, false
}

// Destroy the GlyphTexture (do not use it afterwards)
func (g *GlyphTexture) Destroy() {
	g.t.Destroy()
}
