package main

import (
	"flag"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/torlenor/asciiventure/console"
	"github.com/torlenor/asciiventure/renderers"
	"github.com/torlenor/asciiventure/utils"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	var (
		windowWidth  = flag.Int("w", 1280, "Window width to use")
		windowHeight = flag.Int("h", 800, "Window height to use")
	)

	flag.Parse()
	err := ttf.Init()
	if err != nil {
		log.Fatalf("Failed to initialize ttf: %s", err)
	}

	err = sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		log.Fatalf("Failed to initialize sdl: %s", err)
	}

	window, err := sdl.CreateWindow("Console Test", sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED, int32(*windowWidth), int32(*windowHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s", err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	if runtime.GOOS == "windows" {
		sdl.SetHint(sdl.HINT_RENDER_DRIVER, "opengl")
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s", err)
	}

	myRenderer := renderers.NewRenderer(renderer)

	// tileset := console.NewFontTileset(myRenderer, "./assets/textures/terminal10x10_gs_tc.png")
	// tileset := console.NewFontTileset(myRenderer, "./assets/textures/symbols64x64.png")
	tileset, _ := console.NewFontTileset(myRenderer, "./assets/textures/courier12x12_aa_tc.png")
	mconsole := console.NewMatrixConsole(myRenderer, tileset, int32(tileset.GetCharWidth()*80), int32(tileset.GetCharHeight()*50), 80, 50)

	mconsole.SetOffset(100, 100)

	arbitraryTextLines := []string{"", " Hello fellow programmer!", "", " This is a console test"}

	ticker := time.NewTicker(time.Second / 2)
	quit := false
	for !quit {
		mconsole.Clear()
		x := int32(0)
		y := int32(0)
		y++
		for _, line := range arbitraryTextLines {
			for _, c := range line {
				fc := utils.ColorRGBA{R: uint8(rand.Intn(255)), G: 0, B: 0, A: 255}
				mconsole.PutCharColor(x, y, strings.ToUpper(string(c)), fc, utils.ColorRGBA{})
				x++
			}
			x = 0
			y++
		}

		mconsole.Border(utils.ColorRGBA{R: 255, A: 255}, utils.ColorRGBA{})

		myRenderer.SetDrawColor(255, 255, 255, 255)
		myRenderer.Clear()
		mconsole.Render()
		myRenderer.Present()
		<-ticker.C
	}

	ticker.Stop()
}
