package widgets

import (
	"fmt"
	"pulse/internal/grid"
	"time"
)

type CalendarWidget struct {
	BaseWidget
	Format   string `toml:"format"` // "compact" | "full"
	Timezone string `toml:"timezone"`
}

func (w CalendarWidget) Render(e *grid.Engine) {
	timezone := "Local"
	if w.Timezone != "" {
		timezone = w.Timezone
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		DisplayError(e, w.BaseWidget, "Unknown Timezone")
		return
	}
	now := time.Now().In(loc)

	if w.Format == "compact" {
		date := fmt.Sprintf("%d %s %d", now.Day(), now.Month(), now.Year())
		e.DrawBox(w.X, w.Y, w.Width, w.Height)
		e.DrawText(w.X, w.Y, w.Width, w.Height, date)
	}

	e.DrawBox(w.X, w.Y, w.Width, w.Height)

	monthTitle := fmt.Sprintf(" %s %d ", now.Month(), now.Year())
	e.DrawText(w.X, w.Y, w.Width, 3, monthTitle)

	weeksDays := "Mo Tu We Th Fr Sa Su"
	e.DrawText(w.X, w.Y+2, w.Width, 3, weeksDays)

	currentDay := now.Day()

	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	startOffset := int(firstDayOfMonth.Weekday()) - 1

	if startOffset < 0 {
		startOffset = 6
	}

	totalDays := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, loc).Day()
	lineY := w.Y + 4
	dayCounter := 1

	for row := 0; row < 6; row++ {
		var lineText string

		for col := 0; col < 7; col++ {
			if (row == 0 && col < startOffset) || dayCounter > totalDays {
				lineText += "   "
			} else {
				if dayCounter == currentDay {
					lineText += fmt.Sprintf("[%2d]", dayCounter)
				} else {
					lineText += fmt.Sprintf(" %2d ", dayCounter)
				}
				dayCounter++
			}
		}

		if stringsTrimmed := len(lineText) > 0; stringsTrimmed {
			e.DrawText(w.X, lineY, w.Width, 3, lineText)
			lineY++
		}

		if dayCounter > totalDays {
			break
		}
	}
}
