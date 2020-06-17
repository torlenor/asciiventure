package console

import (
	"log"

	"github.com/torlenor/asciiventure/renderers"
	"github.com/torlenor/asciiventure/utils"
	"github.com/veandco/go-sdl2/sdl"
)

// TileSet interface defines all the necessary functions for a Console
// to retreive the glyphs for rendering.
type TileSet interface {
	Get(c string) (renderers.RenderGlyph, error)
	GetCharWidth() int32
	GetCharHeight() int32
}

// MatrixConsole is used to render chars into fixed cells.
// It provides convenient functions to render ASCII or any other
// tiles onto a grid of specified size.
type MatrixConsole struct {
	renderer *renderers.Renderer
	tileset  TileSet

	consoleWidth   int32
	consoleHeight  int32
	consoleOffsetX int32
	consoleOffsetY int32

	nx int32
	ny int32

	matrix map[int32]map[int32]renderers.RenderGlyph
}

// NewMatrixConsole returns a console with the given dimensions.
// The dimensions nx x ny are number of cells, e.g., 80x50,
// while consoleWidth x consoleHeight are the number of pixels for the console to use.
func NewMatrixConsole(r *renderers.Renderer, tileset TileSet, w, h, nx, ny int32) *MatrixConsole {
	return &MatrixConsole{consoleWidth: w, consoleHeight: h, nx: nx, ny: ny, tileset: tileset, renderer: r, matrix: make(map[int32]map[int32]renderers.RenderGlyph)}
}

// GetDimensions returns the number of tiles in x and y direction of the console.
func (c *MatrixConsole) GetDimensions() (nx, ny int32) {
	return c.nx, c.ny
}

// SetOffset shifts the console by the amount of pixels provided.
func (c *MatrixConsole) SetOffset(x, y int32) {
	c.consoleOffsetX = x
	c.consoleOffsetY = y
}

// GetOffset returns the currently set offset in pixels.
func (c *MatrixConsole) GetOffset() (x, y int32) {
	return c.consoleOffsetX, c.consoleOffsetY
}

// Render the console
func (c *MatrixConsole) Render() {
	charWidth := c.tileset.GetCharWidth()
	charHeight := c.tileset.GetCharHeight()

	var borderWidthHalf int32
	var borderHeightHalf int32
	if charWidth*c.nx > c.consoleWidth {
		log.Printf("Cannot render the whole console onto the available screen width: %d > %d", charWidth*c.nx, c.consoleWidth)
	} else {
		borderWidthHalf = (c.consoleWidth - charWidth*c.nx) / 2
	}
	if charHeight*c.ny > c.consoleHeight {
		log.Printf("Cannot render the whole console onto the available screen height: %d > %d", charHeight*c.ny, c.consoleHeight)
	} else {
		borderHeightHalf = (c.consoleHeight - charHeight*c.ny) / 2
	}

	cr, cg, cb, ca, _ := c.renderer.GetRenderer().GetDrawColor()
	var bm sdl.BlendMode
	c.renderer.GetRenderer().GetDrawBlendMode(&bm)

	for x, row := range c.matrix {
		for y, g := range row {
			dst := &sdl.Rect{
				X: c.consoleOffsetX + borderWidthHalf + x*charWidth,
				Y: c.consoleOffsetY + borderHeightHalf + y*charHeight,
				W: charWidth, H: charHeight,
			}
			// Render background
			err := g.T.SetColorMod(1, 1, 1)
			if err != nil {
				log.Printf("Error setting Color in MatrixConsole Render: %s", err)
			}
			err = g.T.SetAlphaMod(255)
			if err != nil {
				log.Printf("Error setting Alpha in MatrixConsole Render: %s", err)
			}
			if g.BackgroundColor.A > 0 {
				c.renderer.GetRenderer().SetDrawBlendMode(sdl.BLENDMODE_BLEND)
				c.renderer.SetDrawColor(g.BackgroundColor.R, g.BackgroundColor.G, g.BackgroundColor.B, g.BackgroundColor.A)
				c.renderer.GetRenderer().FillRect(dst)
			}
			// Render foreground
			err = g.T.SetColorMod(g.ForegroundColor.R, g.ForegroundColor.G, g.ForegroundColor.B)
			if err != nil {
				log.Printf("Error setting Color in MatrixConsole Render: %s", err)
			}
			err = g.T.SetAlphaMod(g.ForegroundColor.A)
			if err != nil {
				log.Printf("Error setting Alpha in MatrixConsole Render: %s", err)
			}
			err = c.renderer.Copy(g.T, g.Src, dst)
			if err != nil {
				log.Printf("Error in RenderWithOffset: %s", err)
			}
		}
	}

	c.renderer.SetDrawColor(cr, cg, cb, ca)
	c.renderer.GetRenderer().SetDrawBlendMode(bm)
}

// PutChar draws a character on the console using the default colors.
// x: The x coordinate, the left-most position being 0.
// y: The y coordinate, the top-most position being 0.
func (c *MatrixConsole) PutChar(x, y int32, char string) {
	if x < 0 || x >= c.nx || y < 0 || y >= c.ny {
		return
	}
	if _, ok := c.matrix[x]; !ok {
		c.matrix[x] = make(map[int32]renderers.RenderGlyph)
	}
	glyph, err := c.tileset.Get(char)
	if err != nil {
		log.Printf("Error getting glyph: %s", err)
		return
	}
	c.matrix[x][y] = glyph
}

// PutCharColor draws a character on the console with the given colors.
// x: The x coordinate, the left-most position being 0.
// y: The y coordinate, the top-most position being 0.
func (c *MatrixConsole) PutCharColor(x, y int32, char string, foregroundColor utils.ColorRGBA, backgroundColor utils.ColorRGBA) {
	if x < 0 || x >= c.nx || y < 0 || y >= c.ny {
		return
	}
	if _, ok := c.matrix[x]; !ok {
		c.matrix[x] = make(map[int32]renderers.RenderGlyph)
	}
	glyph, err := c.tileset.Get(char)
	if err != nil {
		log.Printf("Error getting glyph: %s", err)
		return
	}
	glyph.ForegroundColor = foregroundColor
	glyph.BackgroundColor = backgroundColor
	c.matrix[x][y] = glyph
}

// HLine draws a horizontal line with length l on the console with the default colors.
// x: The starting x coordinate, the left-most position being 0.
// y: The starting y coordinate, the top-most position being 0.
// Note: It relies on the char "─" in the used tileset.
func (c *MatrixConsole) HLine(x, y, l int32, foregroundColor utils.ColorRGBA, backgroundColor utils.ColorRGBA) {
	for i := int32(0); i < l; i++ {
		c.PutCharColor(x+i, y, "─", foregroundColor, backgroundColor)
	}
}

// VLine draws a vertical line with length l on the console with the default colors.
// x: The starting x coordinate, the left-most position being 0.
// y: The starting y coordinate, the top-most position being 0.
// Note: It relies on the char "─" in the used tileset.
func (c *MatrixConsole) VLine(x, y, l int32, foregroundColor utils.ColorRGBA, backgroundColor utils.ColorRGBA) {
	for i := int32(0); i < l; i++ {
		c.PutCharColor(x, y+i, "│", foregroundColor, backgroundColor)
	}
}

// Border draws a border with ─, │, ┌, └, ┐ and ┘ characters around the whole console.
func (c *MatrixConsole) Border(foregroundColor utils.ColorRGBA, backgroundColor utils.ColorRGBA) {
	c.HLine(0, 0, c.nx, utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})
	c.HLine(0, c.ny-1, c.nx, utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})

	c.VLine(0, 0, c.ny, utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})
	c.VLine(c.nx-1, 0, c.ny, utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})

	c.PutCharColor(0, 0, "┌", utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})
	c.PutCharColor(0, c.ny-1, "└", utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})
	c.PutCharColor(c.nx-1, 0, "┐", utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})
	c.PutCharColor(c.nx-1, c.ny-1, "┘", utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})
}

// Clear the console.
func (c *MatrixConsole) Clear() {
	c.matrix = make(map[int32]map[int32]renderers.RenderGlyph)
}
