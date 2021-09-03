package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/renderers"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type inventoryEntries map[string]int

// InventoryWidget is able to render a inventory.
type InventoryWidget struct {
	renderer *renderers.Renderer
	font     *ttf.Font
	dst      *sdl.Rect

	border *sdl.Texture
	text   *sdl.Texture

	inventoryEntries []string

	textW int32
	textH int32

	wrapLength int
	maxRows    int

	drawBorder bool
}

// NewInventoryWidget returns a new NewInventoryWidget
func NewInventoryWidget(r *renderers.Renderer, font *ttf.Font, dst *sdl.Rect, drawBorder bool) *InventoryWidget {
	return &InventoryWidget{
		renderer:   r,
		font:       font,
		dst:        dst,
		wrapLength: 1000,
		maxRows:    int(dst.H) / font.Height(),
		drawBorder: drawBorder,
	}
}

// UpdateInventory updates the inventory with the new list of entities.
func (w *InventoryWidget) UpdateInventory(inventory *entity.Inventory) {
	w.inventoryEntries = []string{}
	for _, item := range inventory.GetInventoryList() {
		w.inventoryEntries = append(w.inventoryEntries, item)
	}
	w.createTexture()
}

// SetWrapLength defines a new wrap length on how many pixel the text should be wrapped automatically.
func (w *InventoryWidget) SetWrapLength(wrapLength int) {
	w.wrapLength = wrapLength
}

func getJoinedInventoryText(r []string) string {
	var lines []string
	cnt := 0
	for _, name := range r {
		cnt++
		lines = append(lines, fmt.Sprintf("%d) %s", cnt, name))
	}
	return strings.Join(lines, "\n")
}

func (w *InventoryWidget) createTexture() {
	text := "Inventory\n--------------------\n"
	text += getJoinedInventoryText(w.inventoryEntries)
	surface, err := w.font.RenderUTF8BlendedWrapped(text, sdl.Color{R: 255, G: 255, B: 255, A: 255}, w.wrapLength)
	if err != nil {
		log.Printf("Error rendering inventory text: %s", err)
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

// Render the InventoryWidget at the given dst rectangle with the given renderer.
func (w *InventoryWidget) Render() {

	r := w.renderer.GetRenderer()
	cr, cg, cb, ca, _ := r.GetDrawColor()
	var bm sdl.BlendMode
	r.GetDrawBlendMode(&bm)
	r.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	r.SetDrawColor(0, 0, 0, 255)
	// r.SetClipRect(w.dst) // Somehow this breaks rendering
	r.FillRect(w.dst)
	if w.drawBorder {
		r.SetDrawColor(255, 255, 255, 255)
		r.DrawRect(w.dst)
	}
	r.SetDrawColor(cr, cg, cb, ca)
	r.SetDrawBlendMode(bm)

	ldst := *w.dst
	ldst.X += 4
	ldst.Y += 4
	ldst.W = w.textW
	ldst.H = w.textH

	r.Copy(w.text, nil, &ldst)

	// r.SetClipRect(nil)
}
