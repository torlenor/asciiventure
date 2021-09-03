package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/torlenor/asciiventure/renderers"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type rowEntry struct {
	text  string
	count int
}

// TextWidget is able to render character stats
type TextWidget struct {
	renderer *renderers.Renderer
	font     *ttf.Font
	dst      *sdl.Rect

	border *sdl.Texture
	text   *sdl.Texture

	textRows []rowEntry

	textW int32
	textH int32

	wrapLength int
	maxRows    int

	drawBorder bool
}

// NewTextWidget returns a new TextWidget
func NewTextWidget(r *renderers.Renderer, font *ttf.Font, dst *sdl.Rect, drawBorder bool) *TextWidget {
	return &TextWidget{
		renderer:   r,
		font:       font,
		dst:        dst,
		wrapLength: 1000,
		maxRows:    int(dst.H) / font.Height(),
		drawBorder: drawBorder,
	}
}

// AddRow adds a new line of text.
// If number of lines > max lines, the oldest will be removed.
func (w *TextWidget) AddRow(row string) {
	if len(w.textRows) > 0 {
		if w.textRows[len(w.textRows)-1].text == row {
			w.textRows[len(w.textRows)-1].count++
		} else {
			w.textRows = append(w.textRows, rowEntry{text: row, count: 1})
		}
	} else {
		w.textRows = append(w.textRows, rowEntry{text: row, count: 1})
	}
	if len(w.textRows) > w.maxRows {
		w.textRows = w.textRows[1:]
	}
	w.createTexture()
}

// SetWrapLength defines a new wrap length on how many pixel the text should be wrapped automatically.
func (w *TextWidget) SetWrapLength(wrapLength int) {
	w.wrapLength = wrapLength
}

// SetText changes the current text to the rows provided as an argument.
func (w *TextWidget) SetText(rows []string) {
	w.textRows = []rowEntry{}
	for _, s := range rows {
		w.textRows = append(w.textRows, rowEntry{text: s, count: 1})
	}
	w.createTexture()
}

func getText(r rowEntry) string {
	if r.count > 1 {
		return fmt.Sprintf("%s (%dx)", r.text, r.count)
	}
	return fmt.Sprintf("%s", r.text)
}

func getJoinedText(r []rowEntry) string {
	var lines []string
	for _, entry := range r {
		lines = append(lines, getText(entry))
	}
	return strings.Join(lines, "\n")
}

func (w *TextWidget) createTexture() {
	if len(w.textRows) == 0 {
		return
	}
	surface, err := w.font.RenderUTF8BlendedWrapped(fmt.Sprintf("%s", getJoinedText(w.textRows)), sdl.Color{R: 255, G: 255, B: 255, A: 255}, w.wrapLength) // we only want manual wrapping and therefore set the wrapLength kinda large
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
	r.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	r.SetDrawColor(0, 0, 0, 255)
	// r.SetClipRect(w.dst) // This breaks rendering somehow
	r.FillRect(w.dst)
	if w.drawBorder {
		r.SetDrawColor(255, 255, 255, 255)
		r.DrawRect(w.dst)
	}
	r.SetDrawColor(cr, cg, cb, ca)
	r.SetDrawBlendMode(bm)

	if len(w.textRows) != 0 {
		ldst := *w.dst
		ldst.X += 4
		ldst.Y += 4
		ldst.W = w.textW
		ldst.H = w.textH

		r.Copy(w.text, nil, &ldst)
	}

	// r.SetClipRect(nil)
}

// Clear clears the content of the TextWidget.
func (w *TextWidget) Clear() {
	w.SetText([]string{})
}
