package game

import (
	"fmt"
	"time"
)

// GameLoop is a blocking function actually running the game.
func (g *Game) GameLoop() {
	g.createGlyphTexture()
	g.createPlayer()
	g.loadRoomsFromDirectory("./assets/rooms")
	g.selectRoom(1)

	ticker := time.NewTicker(time.Second / 15)

	g.currentRoom.UpdateFoV(playerViewRange, g.player.Position.X, g.player.Position.Y)
	g.updateCharacterWindow()
	g.logWindow.SetText([]string{"Welcome to <Epic Name Here>.", "A small cat takes a stroll and ends up in an epic adventure."})

	for !g.quit {
		start := time.Now()
		g.handleSDLEvents()
		g.markedPath = g.determinePathPlayerMouse()
		if g.gameState != gameOver {
			g.timestep()
		}
		gameLogicUpdateMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0

		start = time.Now()
		g.draw()
		drawUpdateMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0
		start = time.Now()
		<-ticker.C
		spareMs := float32(time.Now().Sub(start).Microseconds()) / 1000.0
		if true {
			fmt.Printf("Game logic duration: %.2f ms, draw duration: %.2f ms, total: %.2f ms, spare: %.2f ms\n", gameLogicUpdateMs, drawUpdateMs, gameLogicUpdateMs+drawUpdateMs, spareMs)
		}
	}

	ticker.Stop()
}
