package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/torlenor/asciiventure/renderers"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TextWidget is able to render character stats
type TextWidget struct {
	renderer *renderers.Renderer
	font     *ttf.Font
	dst      *sdl.Rect

	border *sdl.Texture
	text   *sdl.Texture

	textRows []string

	textW int32
	textH int32

	wrapLength int
	maxRows    int
}

// NewTextWidget returns a new TextWidget
func NewTextWidget(r *renderers.Renderer, font *ttf.Font, dst *sdl.Rect) *TextWidget {
	return &TextWidget{
		renderer:   r,
		font:       font,
		dst:        dst,
		wrapLength: 1000,
		maxRows:    int(dst.H) / font.Height(),
	}
}

// AddRow adds a new line of text.
// If number of lines > max lines, the oldest will be removed.
func (w *TextWidget) AddRow(row string) {
	w.textRows = append(w.textRows, row)
	if len(w.textRows) > w.maxRows {
		w.textRows = w.textRows[1:]
	}
	w.createTexure()
}

// SetWrapLength defines a new wrap length on how many pixel the text should be wrapped automatically.
func (w *TextWidget) SetWrapLength(wrapLength int) {
	w.wrapLength = wrapLength
}

// SetText changes the current text to the rows provided as an argument.
func (w *TextWidget) SetText(rows []string) {
	w.textRows = rows
	w.createTexure()
}

func (w *TextWidget) createTexure() {
	if len(w.textRows) == 0 {
		return
	}
	surface, err := w.font.RenderUTF8BlendedWrapped(fmt.Sprintf("%s", strings.Join(w.textRows, "\n")), sdl.Color{R: 0, G: 255, B: 0, A: 255}, w.wrapLength) // we only want manual wrapping and therefore set the wrapLength kinda large
	if err != nil {
		log.Printf("Error rendering text: %s", err)
		return
	}
	defer surface.Free()

	if w.text != nil {
		w.text.Destroy()
		w.text = nil
	}

	w.text, err = w.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		fmt.Printf("Failed to create texture from surface when trying to render game time: %s\n", err)
		return
	}

	w.textW = surface.W
	w.textH = surface.H
}

// Render the TextWidget at the given dst rectangle with the given renderer.
func (w *TextWidget) Render() {

	r := w.renderer.GetRenderer()
	cr, cg, cb, ca, _ := r.GetDrawColor()
	var bm sdl.BlendMode
	r.GetDrawBlendMode(&bm)
	// clipRect := r.GetRenderer().GetClipRect()
	r.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	r.SetDrawColor(0, 0, 0, 255)
	r.SetClipRect(w.dst)
	r.FillRect(w.dst)
	r.SetDrawColor(100, 255, 100, 255)
	r.DrawRect(w.dst)
	r.SetDrawColor(cr, cg, cb, ca)
	r.SetDrawBlendMode(bm)

	if len(w.textRows) != 0 {
		ldst := *w.dst
		ldst.W = w.textW
		ldst.H = w.textH

		r.Copy(w.text, nil, &ldst)
	}

	r.SetClipRect(nil)
}
