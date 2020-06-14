package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

// commandType is the type of command.
type command int

// List of commandTypes.
const (
	CommandUnknown command = iota
	CommandQuit
	CommandMoveN
	CommandMoveNE
	CommandMoveE
	CommandMoveSE
	CommandMoveS
	CommandMoveSW
	CommandMoveW
	CommandMoveNW
	CommandZoomOut
	CommandZoomIn
	CommandScrollLeft
	CommandScrollRight
	CommandScrollUp
	CommandScrollDown
	CommandNextTimeStep
	CommandInteract
	CommandSelect1
	CommandSelect2
	CommandSelect3
	CommandSelect4
	CommandSelect5
	CommandSelect6
	CommandSelect7
	CommandSelect8
	CommandSelect9
	CommandSelect0
	CommandAltSelect1
	CommandAltSelect2
	CommandAltSelect3
	CommandAltSelect4
	CommandAltSelect5
	CommandAltSelect6
	CommandAltSelect7
	CommandAltSelect8
	CommandAltSelect9
	CommandAltSelect0
	CommandDebugReload
)

type commandObserver interface {
	NotifyCommand(command)
}

type mouseObserver interface {
	NotifyMouseCommand(buttonLeft, buttonMiddle, buttonRight bool, x, y int32)
}

type registeredCommand struct {
	command command
	name    string
	key     int
	shift   bool
	ctrl    bool
	alt     bool
	pressed bool
}

type commandManager struct {
	registeredCommands []registeredCommand
	observers          []commandObserver
	mouseObservers     []mouseObserver
}

func (c *commandManager) RegisterCommand(command command, name string, key int, shift, ctrl, alt, pressed bool) {
	c.registeredCommands = append(c.registeredCommands, registeredCommand{
		command: command,
		name:    name,
		key:     key,
		shift:   shift,
		ctrl:    ctrl,
		alt:     alt,
		pressed: pressed,
	})
}

func (c *commandManager) RegisterObserver(observer commandObserver) {
	c.observers = append(c.observers, observer)
}

func (c *commandManager) RegisterMouseObserver(observer mouseObserver) {
	c.mouseObservers = append(c.mouseObservers, observer)
}

func (c *commandManager) DispatchMouseCommand(event sdl.Event) {
	var buttonLeftPressed bool
	var buttonRightPressed bool
	var buttonMiddlePressed bool

	switch t := event.(type) {
	case *sdl.MouseMotionEvent:
		if int32(t.State&sdl.ButtonLMask()) > 0 {
			buttonLeftPressed = true
		}
		if int32(t.State&sdl.ButtonRMask()) > 0 {
			buttonRightPressed = true
		}
		if int32(t.State&sdl.ButtonMMask()) > 0 {
			buttonMiddlePressed = true
		}
		for _, observer := range c.mouseObservers {
			observer.NotifyMouseCommand(buttonLeftPressed, buttonMiddlePressed, buttonRightPressed, t.X, t.Y)
		}
	case *sdl.MouseButtonEvent:
		if t.State == sdl.PRESSED {
			switch t.Button {
			case sdl.BUTTON_LEFT:
				buttonLeftPressed = true
			case sdl.BUTTON_MIDDLE:
				buttonMiddlePressed = true
			case sdl.BUTTON_RIGHT:
				buttonRightPressed = true
			}
		}
		for _, observer := range c.mouseObservers {
			observer.NotifyMouseCommand(buttonLeftPressed, buttonMiddlePressed, buttonRightPressed, -1, -1)
		}
	}
}

// DispatchCommand will dispatch the command to its observers if it is registered.
func (c *commandManager) DispatchCommand(event sdl.Event) {
	key := -1
	shift := false
	ctrl := false
	alt := false
	pressed := false

	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		switch t.Keysym.Mod {
		case sdl.KMOD_LSHIFT:
			fallthrough
		case sdl.KMOD_RSHIFT:
			fallthrough
		case sdl.KMOD_SHIFT:
			shift = true
		case sdl.KMOD_LCTRL:
			fallthrough
		case sdl.KMOD_RCTRL:
			fallthrough
		case sdl.KMOD_CTRL:
			ctrl = true
		case sdl.KMOD_RALT:
			fallthrough
		case sdl.KMOD_LALT:
			fallthrough
		case sdl.KMOD_ALT:
			alt = true
		}

		if t.State == sdl.PRESSED {
			pressed = true
		}

		key = int(t.Keysym.Sym)
	}

	for _, registeredCommand := range c.registeredCommands {
		if alt == registeredCommand.alt &&
			ctrl == registeredCommand.ctrl &&
			shift == registeredCommand.shift &&
			pressed == registeredCommand.pressed &&
			key == registeredCommand.key {
			for _, observer := range c.observers {
				observer.NotifyCommand(registeredCommand.command)
			}
		}
	}

}
