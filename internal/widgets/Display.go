package widgets

import "pulse/internal/grid"

type DisplayWidget struct {
	BaseWidget
	Value string `toml:"value"`
}

func (w DisplayWidget) Render(e *grid.Engine) {
	e.DrawBox(w.X, w.Y, w.Width, w.Height)
	e.DrawText(w.X, w.Y, w.Width, w.Height, w.Value)
}
