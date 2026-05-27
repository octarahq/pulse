package widgets

import (
	"pulse/internal/grid"

	"github.com/BurntSushi/toml"
)

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

type Factory func(prim toml.Primitive, meta *toml.MetaData) (Widget, error)

var Registry = make(map[string]Factory)

func Register(widgetType string, factory Factory) {
	Registry[widgetType] = factory
}
