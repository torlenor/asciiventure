package game

import (
	"fmt"
	"time"
)

// GameLoop is a blocking function actually running the game.
func (g *Game) GameLoop() {
	ticker := time.NewTicker(time.Second / 15)
	for !g.quit {
		start := time.Now()
		g.handleSDLEvents()
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
