package renderers

import (
	"log"

	"github.com/torlenor/asciiventure/components"
	"github.com/veandco/go-sdl2/sdl"
)

// Renderer is able to render actual elements onto the correct screen positions.
type Renderer struct {
	GlyphWidth  int
	GlyphHeight int

	OriginX int
	OriginY int

	renderer *sdl.Renderer
}

func NewRenderer(renderer *sdl.Renderer) *Renderer {
	return &Renderer{renderer: renderer}
}

func (r *Renderer) Destroy() {
	r.renderer.Destroy()
}

func (r *Renderer) GetRenderer() *sdl.Renderer {
	return r.renderer
}

func (r *Renderer) SetRenderTarget(texture *sdl.Texture) error {
	return r.renderer.SetRenderTarget(texture)
}

func (r *Renderer) CreateTexture(format uint, access int, w int, h int) (*sdl.Texture, error) {
	return r.renderer.CreateTexture(uint32(format), access, int32(w), int32(h))
}

func (r *Renderer) CreateTextureFromSurface(surface *sdl.Surface) (*sdl.Texture, error) {
	return r.renderer.CreateTextureFromSurface(surface)
}

func (r *Renderer) Present() {
	r.renderer.Present()
}

func (r *Renderer) Clear() {
	r.renderer.Clear()
}

func (r *Renderer) SetScale(scaleX float32, scaleY float32) error {
	return r.renderer.SetScale(scaleX, scaleY)
}

func (r *Renderer) SetDrawColor(re, g, b, a uint8) error {
	return r.renderer.SetDrawColor(re, g, b, a)
}

func (r *Renderer) Copy(texture *sdl.Texture, src *sdl.Rect, dst *sdl.Rect) error {
	return r.renderer.Copy(texture, src, dst)
}

// Render renders a texture starting at the upper left corner of the given character coordinate.
func (r *Renderer) Render(t *sdl.Texture, src *sdl.Rect, cx, cy, w, h int) {
	r.RenderWithOffset(t, src, cx, cy, w, h, 0, 0)
}

// RenderWithOffset renders a texture starting at the upper left corner at given character coordinate with the given pixel offset.
func (r *Renderer) RenderWithOffset(t *sdl.Texture, src *sdl.Rect, cx, cy, w, h, offsetX, offsetY int) {
	dst := &sdl.Rect{X: int32((cx+r.OriginX)*r.GlyphWidth - offsetX), Y: int32((cy+1+r.OriginY)*r.GlyphHeight - offsetY), W: int32(w), H: int32(h)}
	err := r.renderer.Copy(t, src, dst)
	if err != nil {
		log.Printf("Error in RenderWithOffset: %s", err)
	}
}

// RenderGlyph renders a glyph at the given character coordinate
func (r *Renderer) RenderGlyph(g components.Glyph, cx, cy int) {
	err := g.T.SetColorMod(g.Color.R, g.Color.G, g.Color.B)
	if err != nil {
		log.Printf("Error setting Color in RenderGlyph: %s", err)
	}
	r.RenderWithOffset(g.T, g.Src, cx, cy, g.Width, g.Height, g.OffsetX, g.OffsetY)
}

// FillCharCoordinate Draws a rectangle completely filling the given char coordinate
func (r *Renderer) FillCharCoordinate(cx, cy int, c components.ColorRGBA) {
	cr, cg, cb, ca, _ := r.renderer.GetDrawColor()
	var bm sdl.BlendMode
	r.renderer.GetDrawBlendMode(&bm)
	r.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	r.renderer.SetDrawColor(c.R, c.G, c.B, c.A)
	r.renderer.FillRect(&sdl.Rect{X: int32((cx + r.OriginX) * r.GlyphWidth), Y: int32((cy + r.OriginY) * r.GlyphHeight), W: int32(r.GlyphWidth), H: int32(r.GlyphHeight)})
	r.renderer.SetDrawColor(cr, cg, cb, ca)
	r.renderer.SetDrawBlendMode(bm)
}
