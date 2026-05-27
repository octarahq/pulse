package display

import (
	"fmt"
	"pulse/internal/config"
	"pulse/internal/widgets"
)

var ConsoleWidth = 120

func RenderDisplay(cfg config.Config) {
	fmt.Print("\033[H")

	wid := cfg.Rows[0][0]

	if dw, ok := wid.(widgets.DisplayWidget); ok {
		display := widgets.Display(dw, 10, 3)
		Draw(DrawBox(10, 3, true, display))
	} else {
		fmt.Errorf("Cannot render, unknown widget")
	}
}
