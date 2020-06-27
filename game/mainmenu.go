package game

import (
	"github.com/torlenor/asciiventure/console"
	"github.com/torlenor/asciiventure/utils"
)

const (
	maxOptions = 4
)

var (
	fc = utils.ColorRGBA{R: 255, G: 255, B: 255, A: 255}
	bc = utils.ColorRGBA{}

	fcSelected = utils.ColorRGBA{R: 0, G: 0, B: 0, A: 255}
	bcSelected = utils.ColorRGBA{R: 255, G: 255, B: 255, A: 255}
)

// TODO: Move MainMenu out of the 'game' package

// MainMenu represents the main menu of the game
type MainMenu struct {
	selectedOption int32
}

// Render the main menu on the provided console
func (g *MainMenu) Render(console *console.MatrixConsole, gameInProgress bool) {
	logo := []string{
		`,-~w`,
		`  ,                                          -^,>^1`,
		` /      C              .             ..  z ;` + "`" + `,*    $`,
		`       C            j  .  '*w;.+ '~"` + "`" + `       ]=     ]`,
		` y    'y                ` + "`" + `\                   t|    jp`,
		`  Y,     ` + "`" + `"^^^^^^ -. L     \                   ` + "`" + `   AH`,
		`    "%w              Y    Q-                        b`,
		`         '"""""""~,       ` + "`" + `                         '`,
		`                   "                                .`,
		`                     w,J                 a"2Y`,
		`                     ]L ,     4M2mV      Mw5A`,
		`                     M  ` + "`" + `v     *w<L                -`,
		`                    j      .              >*""` + "`" + `` + "`" + ` ,+`,
		`                    |       ^Lz~:^,  *C/` + "`" + `  ~<JYo^ ]`,
		`                    !       ` + "`" + `  >F*<--a+9%m*` + "`" + `` + "`" + `` + "`" + ` *  ]`,
		`                     L   !                        ]`,
		`                     '    V\                      ]`,
		`                      \    yL                    ,M`,
		`                       Y    Y          p         M`,
		`                        Y    \         ##L      ;`,
		`                         ` + "`" + `y   V` + "`" + `       1K      ,`,
		`                           "p  Y@       0     ,`,
		`                            "m  ""      'p   ]`,
		`                             [   ` + "`" + `g      1p ]C`,
		`                              ,   ?H      "VA`,
		`                              ` + "`" + `""` + "`" + ` H        !L`,
		`                                   Yp     G LH`,
		`                                     ***f*""` + "`",
		"",
		"                     Lili's Quest",
		"A game featuring a little cat, monsters and mutations.",
		"",
	}
	options := []string{
		"                      New_Game",
		"                      Load_Game",
		"                      Options",
		"                      Quit",
	}

	console.Clear()
	x := int32(0)
	y := int32(0)
	y++
	for _, line := range logo {
		for _, c := range line {
			console.PutCharColor(x, y, string(c), fc, bc)

			x++
		}
		x = 0
		y++
	}

	x = 0
	y++
	for n, line := range options {
		for _, c := range line {
			if c == ' ' {
				x++
				continue
			}
			if int32(n) == g.selectedOption {
				console.PutCharColor(x, y, string(c), fcSelected, bcSelected)
			} else {
				console.PutCharColor(x, y, string(c), fc, bc)
			}
			x++
		}
		x = 0
		y++
	}
}

// MoveCursor moves the cursor of the currently selected item.
func (g *MainMenu) MoveCursor(dx, dy int32) {
	g.selectedOption += dy
	if g.selectedOption >= maxOptions {
		g.selectedOption = 0
	}
	if g.selectedOption < 0 {
		g.selectedOption = maxOptions - 1
	}
}

// Select selects the currently activated cursor element.
func (g *MainMenu) Select() MainMenuActionType {
	switch g.selectedOption {
	case 0:
		return MainMenuActionStartGame
	case 1:
		return MainMenuActionLoadGame
	case 2:
		return MainMenuActionOptions
	case 3:
		return MainMenuActionQuit
	default:
		return MainMenuActionUnknown
	}
}
