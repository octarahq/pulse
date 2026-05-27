package widgets

import "pulse/internal/grid"

type Widget interface {
	Render(e *grid.Engine)
}

type BaseWidget struct {
	ID     int    `toml:"id"`
	Type   string `toml:"type"`
	X      int    `toml:"x"`
	Y      int    `toml:"y"`
	Width  int    `toml:"width"`
	Height int    `toml:"height"`
}
