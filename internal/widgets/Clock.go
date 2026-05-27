package widgets

import (
	"fmt"
	"pulse/internal/grid"
	"time"
)

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
