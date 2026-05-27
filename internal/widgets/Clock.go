package widgets

import (
	"fmt"
	"pulse/internal/grid"
	"time"

	"github.com/BurntSushi/toml"
)

func init() {
	Register("clock", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w ClockWidget
		err := meta.PrimitiveDecode(prim, &w)
		return w, err
	})
}

type ClockWidget struct {
	BaseWidget
	Timezone string `toml:"timezone"`
}

func (w ClockWidget) Render(e *grid.Engine) {
	var value string

	loc, err := time.LoadLocation(w.Timezone)
	if err != nil {
		value = "Unknown timezone"
	}
	now := time.Now().In(loc)
	value = fmt.Sprintf("%d:%d", now.Hour(), now.Minute())

	e.DrawBox(w.X, w.Y, w.Width, w.Height)
	e.DrawText(w.X, w.Y, w.Width, w.Height, value)
}
