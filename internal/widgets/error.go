package widgets

import (
	"pulse/internal/grid"
)

func DisplayError(e *grid.Engine, w BaseWidget, value string) {
	e.DrawBox(w.X, w.Y, w.Width, w.Height)
	e.DrawText(w.X, w.Y, w.Width, w.Height, value)
}
