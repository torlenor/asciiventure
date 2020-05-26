package main

import "github.com/torlenor/asciiventure/game"

func main() {
	game := &game.Game{}
	game.Setup()
	game.GameLoop()
	game.Shutdown()
}
