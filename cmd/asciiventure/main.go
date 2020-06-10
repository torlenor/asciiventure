package main

import (
	"flag"

	"github.com/torlenor/asciiventure/game"
)

func main() {
	var (
		windowWidth  = flag.Int("w", 1024, "Window width to use")
		windowHeight = flag.Int("h", 768, "Window height to use")
		f            = flag.Bool("f", false, "Start in fullscreen mode")
	)

	flag.Parse()

	game := &game.Game{}
	game.Setup(*windowWidth, *windowHeight, *f)
	game.GameLoop()
	game.Shutdown()
}
