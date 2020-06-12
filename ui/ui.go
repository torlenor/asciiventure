package ui

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
	"github.com/torlenor/asciiventure/renderers"
)

// UI holds all functions and data related to the UI.
type UI struct {
	r *renderers.Renderer

	fontSize int
	font     *ttf.Font

	screenWidth  int
	screenHeight int

	roomRenderPane      sdl.Rect
	characterWindowRect sdl.Rect
	logWindowRect       sdl.Rect
	statusBarRec        sdl.Rect
	mutationsRect       sdl.Rect
	inventoryRect       sdl.Rect

	characterWindow  *TextWidget
	logWindow        *TextWidget
	statusBar        *TextWidget
	mutations        *TextWidget
	inventory        *InventoryWidget
	inventoryEnabled bool
}

// NewUI creates a new UI.
func NewUI(r *renderers.Renderer, font *ttf.Font, fontSize int) *UI {
	ui := &UI{
		font:     font,
		fontSize: fontSize,
		r:        r,
	}

	return ui
}

// SetScreenDimensions sets a new width and height for the current window where the UI is rendered.
// UI will calculate from that how to position the UI elements on the screen, so make sure it is always
// current.
func (ui *UI) SetScreenDimensions(width, height int) {
	ui.screenWidth = width
	ui.screenHeight = height

	ui.roomRenderPane = sdl.Rect{X: int32(ui.screenHeight / 6), Y: 0, W: int32(ui.screenWidth), H: int32(ui.screenHeight - ui.screenHeight/6)}
	ui.characterWindowRect = sdl.Rect{X: 0, Y: 0, W: int32(ui.screenWidth / 2), H: int32(ui.screenHeight / 6)}
	ui.logWindowRect = sdl.Rect{X: int32(ui.screenWidth - ui.screenWidth/2 - 1), Y: 0, W: int32(ui.screenWidth/2 + 1), H: int32(ui.screenHeight / 6)}
	ui.statusBarRec = sdl.Rect{X: 0, Y: int32(ui.screenHeight - ui.fontSize - 16 - 1), W: int32(ui.screenWidth), H: int32(ui.fontSize + 16)}
	ui.mutationsRect = sdl.Rect{X: int32(ui.screenWidth - ui.screenWidth/4), Y: int32(ui.screenHeight/6 - 1), W: int32(ui.screenWidth / 4), H: int32(3*ui.screenHeight/6 + 1)}
	ui.inventoryRect = sdl.Rect{X: int32(ui.screenWidth - ui.screenWidth/4), Y: int32(4*ui.screenHeight/6 - 1), W: int32(ui.screenWidth / 4), H: int32(2*int32(ui.screenHeight/6) - ui.statusBarRec.H + 1)}

	ui.characterWindow = NewTextWidget(ui.r, ui.font, &ui.characterWindowRect, true)
	ui.characterWindow.SetWrapLength(int(ui.characterWindowRect.W - 8))
	ui.logWindow = NewTextWidget(ui.r, ui.font, &ui.logWindowRect, true)
	ui.logWindow.SetWrapLength(int(ui.logWindowRect.W - 8))
	ui.statusBar = NewTextWidget(ui.r, ui.font, &ui.statusBarRec, true)
	ui.statusBar.SetWrapLength(int(ui.statusBarRec.W - 8))
	ui.mutations = NewTextWidget(ui.r, ui.font, &ui.mutationsRect, true)
	ui.mutations.SetWrapLength(int(ui.mutationsRect.W - 8))
	ui.mutations.AddRow("No mutations")
	ui.inventory = NewInventoryWidget(ui.r, ui.font, &ui.inventoryRect, true)
	ui.inventory.SetWrapLength(int(ui.inventoryRect.W - 8))
}

// Render the UI.
func (ui *UI) Render() {
	ui.characterWindow.Render()
	ui.logWindow.Render()
	ui.statusBar.Render()
	ui.mutations.Render()
	if ui.inventoryEnabled {
		ui.inventory.Render()
	}
}

// SetStatusBarText sets a new text in the status bar.
func (ui *UI) SetStatusBarText(text string) {
	if len(text) == 0 {
		ui.statusBar.Clear()
	} else {
		ui.statusBar.SetText([]string{text})
	}
}

// UpdateCharacterPane updates the character infos with the information provided.
func (ui *UI) UpdateCharacterPane(time uint, currentHP, totalHP, vision, power, defense int) {
	ui.characterWindow.SetText([]string{
		fmt.Sprintf("Time: %d", time),
		fmt.Sprintf("HP: %d/%d", currentHP, totalHP),
		fmt.Sprintf("Vision: %d", vision),
		fmt.Sprintf("Power %d", power),
		fmt.Sprintf("Defense %d", defense),
	})
}

// UpdateMutationsPane updates the mutation info with the newly provided list.
func (ui *UI) UpdateMutationsPane(mutations components.Mutations) {
	ui.mutations.Clear()
	if len(mutations) == 0 {
		ui.mutations.AddRow("No mutations")
		return
	}
	ui.mutations.AddRow("Mutations:")
	ui.mutations.AddRow("----------------")
	for _, m := range mutations {
		ui.mutations.AddRow(fmt.Sprintf("%s", m))
	}
}

// UpdateInventoryPane updates the inventory info with the newly provided list.
func (ui *UI) UpdateInventoryPane(inventory *entity.Inventory) {
	ui.inventory.UpdateInventory(inventory)
}

// SetInventoryPaneEnabled shows or hides the inventory.
func (ui *UI) SetInventoryPaneEnabled(enabled bool) {
	ui.inventoryEnabled = enabled
}

// AddLogEntry adds a new entry to the log pane.
func (ui *UI) AddLogEntry(text string) {
	ui.logWindow.AddRow(text)
}
