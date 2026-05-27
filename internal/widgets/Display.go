package widgets

import (
	"pulse/internal/grid"

	"github.com/BurntSushi/toml"
)

type DisplayWidget struct {
	BaseWidget
	Value string `toml:"value"`
}

func init() {
	Register("display", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w DisplayWidget
		err := meta.PrimitiveDecode(prim, &w)
		return w, err
	})
}

func (w DisplayWidget) Render(e *grid.Engine) {
	e.DrawBox(w.X, w.Y, w.Width, w.Height)
	e.DrawText(w.X, w.Y, w.Width, w.Height, w.Value)
}
