package game

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/utils"
	"github.com/veandco/go-sdl2/sdl"
)

func (g *Game) setupInput() {
	g.commandManager = &commandManager{}
	g.commandManager.RegisterObserver(g)
	g.commandManager.RegisterMouseObserver(g)

	g.commandManager.RegisterCommand(CommandQuit, "quit", sdl.K_ESCAPE, false, false, false, true)
	g.commandManager.RegisterCommand(CommandQuit, "quit", int('q'), false, false, false, true)

	g.commandManager.RegisterCommand(CommandMoveN, "move_n", sdl.K_UP, false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveE, "move_e", sdl.K_RIGHT, false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveS, "move_s", sdl.K_DOWN, false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveW, "move_w", sdl.K_LEFT, false, false, false, true)

	// y    k       u
	// h            l
	// b    j       n
	g.commandManager.RegisterCommand(CommandMoveN, "move_n", int('k'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveE, "move_e", int('l'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveS, "move_s", int('j'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveW, "move_w", int('h'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveNE, "move_ne", int('u'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveSE, "move_se", int('n'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveSW, "move_sw", int('b'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandMoveNW, "move_nw", int('y'), false, false, false, true)

	g.commandManager.RegisterCommand(CommandNextTimeStep, "next_timestep", sdl.K_SPACE, false, false, false, true)
	g.commandManager.RegisterCommand(CommandZoomIn, "zoom_in", int('+'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandZoomOut, "zoom_out", int('-'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandScrollUp, "scroll_up", sdl.K_UP, false, false, true, true)
	g.commandManager.RegisterCommand(CommandScrollLeft, "scroll_left", sdl.K_LEFT, false, false, true, true)
	g.commandManager.RegisterCommand(CommandScrollDown, "scroll_down", sdl.K_DOWN, false, false, true, true)
	g.commandManager.RegisterCommand(CommandScrollRight, "scroll_right", sdl.K_RIGHT, false, false, true, true)

	g.commandManager.RegisterCommand(CommandInteract, "interact", int('g'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandInteract, "interact", sdl.K_RETURN, false, false, false, true)

	g.commandManager.RegisterCommand(CommandSelect1, "select_1", int('1'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect2, "select_2", int('2'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect3, "select_3", int('3'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect4, "select_4", int('4'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect5, "select_5", int('5'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect6, "select_6", int('6'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect7, "select_7", int('7'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect8, "select_8", int('8'), false, false, false, true)
	g.commandManager.RegisterCommand(CommandSelect9, "select_9", int('9'), false, false, false, true)

	if g.debug {
		g.commandManager.RegisterCommand(CommandAltSelect1, "select_map_1", int('1'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect2, "select_map_2", int('2'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect3, "select_map_3", int('3'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect4, "select_map_4", int('4'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect5, "select_map_5", int('5'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect6, "select_map_6", int('6'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect7, "select_map_7", int('7'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect8, "select_map_8", int('8'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandAltSelect9, "select_map_9", int('9'), false, false, true, true)
		g.commandManager.RegisterCommand(CommandDebugReload, "reload", sdl.K_F5, false, false, false, true)
	}
}

// NotifyCommand will be called from commandManager when a registered command is received.
func (g *Game) NotifyCommand(command command) {
	if g.gameState == gameOver {
		switch command {
		case CommandQuit:
			g.gameState = mainMenu
		}
	} else if g.gameState == playersTurn {
		switch command {
		case CommandQuit:
			g.gameState = mainMenu
		case CommandMoveN:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.Y = g.player.Position.Current.Y - 1
			g.nextStep = true
		case CommandMoveNE:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.X = g.player.Position.Current.X + 1
			g.player.TargetPosition.Y = g.player.Position.Current.Y - 1
			g.nextStep = true
		case CommandMoveE:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.X = g.player.Position.Current.X + 1
			g.nextStep = true
		case CommandMoveSE:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.X = g.player.Position.Current.X + 1
			g.player.TargetPosition.Y = g.player.Position.Current.Y + 1
			g.nextStep = true
		case CommandMoveS:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.Y = g.player.Position.Current.Y + 1
			g.nextStep = true
		case CommandMoveSW:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.X = g.player.Position.Current.X - 1
			g.player.TargetPosition.Y = g.player.Position.Current.Y + 1
			g.nextStep = true
		case CommandMoveW:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.X = g.player.Position.Current.X - 1
			g.nextStep = true
		case CommandMoveNW:
			g.movementPath = []utils.Vec2{}
			g.player.TargetPosition = g.player.Position.Current
			g.player.TargetPosition.X = g.player.Position.Current.X - 1
			g.player.TargetPosition.Y = g.player.Position.Current.Y - 1
			g.nextStep = true
		case CommandScrollUp:
			g.renderer.OriginY += 2
		case CommandScrollLeft:
			g.renderer.OriginX += 2
		case CommandScrollDown:
			g.renderer.OriginY -= 2
		case CommandScrollRight:
			g.renderer.OriginX -= 2
		case CommandZoomIn:
			g.renderScale += 0.1
			g.consoleMap.SetOffset(0, int32(float32(g.screenHeight/6)/g.renderScale))
		case CommandZoomOut:
			g.renderScale -= 0.1
			g.consoleMap.SetOffset(0, int32(float32(g.screenHeight/6)/g.renderScale))
		case CommandNextTimeStep:
			g.nextStep = true
		case CommandInteract:
			g.performPlayerAction(components.ActionTypeInteract, 0)
			g.nextStep = true
		case CommandSelect1:
			g.performPlayerAction(components.ActionTypeUseItem, 0)
			g.nextStep = true
		case CommandSelect2:
			g.performPlayerAction(components.ActionTypeUseItem, 1)
			g.nextStep = true
		case CommandSelect3:
			g.performPlayerAction(components.ActionTypeUseItem, 2)
			g.nextStep = true
		case CommandSelect4:
			g.performPlayerAction(components.ActionTypeUseItem, 3)
			g.nextStep = true
		case CommandSelect5:
			g.performPlayerAction(components.ActionTypeUseItem, 4)
			g.nextStep = true
		case CommandSelect6:
			g.performPlayerAction(components.ActionTypeUseItem, 5)
			g.nextStep = true
		case CommandSelect7:
			g.performPlayerAction(components.ActionTypeUseItem, 6)
			g.nextStep = true
		case CommandSelect8:
			g.performPlayerAction(components.ActionTypeUseItem, 7)
			g.nextStep = true
		case CommandSelect9:
			g.performPlayerAction(components.ActionTypeUseItem, 8)
			g.nextStep = true
		case CommandAltSelect1:
			g.selectGameMap(1)
		case CommandAltSelect2:
			g.selectGameMap(2)
		case CommandAltSelect3:
			g.selectGameMap(3)
		case CommandAltSelect4:
			g.selectGameMap(4)
		case CommandAltSelect5:
			g.selectGameMap(5)
		case CommandAltSelect6:
			g.selectGameMap(6)
		case CommandAltSelect7:
			g.selectGameMap(7)
		case CommandAltSelect8:
			g.selectGameMap(8)
		case CommandAltSelect9:
			g.selectGameMap(9)
		case CommandDebugReload:
			g.loadGameMapsFromDirectory("./assets/rooms")
		}
	} else if g.gameState == mainMenu {
		switch command {
		case CommandInteract:
			switch g.mainMenu.Select() {
			case MainMenuActionStartGame:
				g.gameInProgress = true
				g.gameState = playersTurn
				g.player.IsDead = nil
				// TODO: Reset everything and generate new maps when a new game starts
			case MainMenuActionQuit:
				g.quit = true
			}
		case CommandQuit:
			g.quit = true
		case CommandMoveN:
			g.mainMenu.MoveCursor(0, -1)
		case CommandMoveE:
			g.mainMenu.MoveCursor(1, 0)
		case CommandMoveS:
			g.mainMenu.MoveCursor(0, 1)
		case CommandMoveW:
			g.mainMenu.MoveCursor(-1, 0)
		}
	}
}

// NotifyMouseCommand will be called from commandManager when a mouse event is received.
func (g *Game) NotifyMouseCommand(buttonLeft, buttonMiddle, buttonRight bool, x, y int32) {
	if x >= 0 && y >= 0 {
		g.updateMouseTile(int(x), int(y))
	}
	if buttonLeft {
		g.setTargetPosition(g.mouseTileX, g.mouseTileY)
	}
}
