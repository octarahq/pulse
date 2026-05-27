package widgets

import (
	"fmt"
	"pulse/internal/chars"
	"pulse/internal/grid"
)

func DisplayError(e *grid.Engine, w BaseWidget, value string) {
	e.DrawBox(w.X, w.Y, w.Width, w.Height)
	e.DrawText(w.X, w.Y, w.Width, w.Height, value)
}

func DisplayGauge(e *grid.Engine, w BaseWidget, value float64, max float64, unit string) {
	e.DrawBox(w.X, w.Y, w.Width, w.Height)

	textLabel := fmt.Sprintf(" %.2f/%.2f %s", value, max, unit)
	textRunes := []rune(textLabel)
	textLen := len(textRunes)

	totalInnerWidth := w.Width - 2
	availWidth := totalInnerWidth - textLen

	if availWidth <= 0 {
		return
	}

	ratio := float64(value) / float64(max)
	if ratio > 1.0 {
		ratio = 1.0
	}
	if ratio < 0.0 {
		ratio = 0.0
	}

	totalGaugeBlocks := float64(availWidth) * ratio
	fullBlocksCount := int(totalGaugeBlocks)

	fraction := totalGaugeBlocks - float64(fullBlocksCount)
	fractionIdx := int(fraction * 4.0)

	var gaugeRunes []rune

	for i := 0; i < fullBlocksCount; i++ {
		gaugeRunes = append(gaugeRunes, []rune(chars.ProgressChars[4])[0])
	}

	if len(gaugeRunes) < availWidth && fractionIdx > 0 {
		gaugeRunes = append(gaugeRunes, []rune(chars.ProgressChars[fractionIdx])[0])
	}

	for len(gaugeRunes) < availWidth {
		gaugeRunes = append(gaugeRunes, ' ')
	}

	centerY := w.Y + (w.Height / 2)

	startXText := w.X + w.Width - 1 - textLen
	for i, r := range textRunes {
		x := startXText + i
		if x > w.X && x < w.X+w.Width-1 && centerY >= 0 && centerY < e.Height {
			e.Buffer[centerY][x] = r
		}
	}

	for i := 0; i < len(gaugeRunes); i++ {
		runeToPlace := gaugeRunes[len(gaugeRunes)-1-i]
		x := startXText - 1 - i

		if x > w.X && x < w.X+w.Width-1 && centerY >= 0 && centerY < e.Height {
			e.Buffer[centerY][x] = runeToPlace
		}
	}
}

func DisplayGraph(e *grid.Engine, w BaseWidget, stats []int, maxint, unit string) {
	e.DrawBox(w.X, w.Y, w.Width, w.Height)

	textLabel := fmt.Sprintf(" %s %s ", maxint, unit)
	textRunes := []rune(textLabel)
	textLen := len([]rune(textLabel))

	totalInnerWidth := w.Width - 2
	availWidth := totalInnerWidth - textLen

	if availWidth <= 0 {
		return
	}

	var graphRunes []rune
	for i := 0; i < len(stats)-1 && len(graphRunes) < availWidth; i += 2 {
		hLeft := stats[i]
		hRight := stats[i+1]

		if hLeft < 0 {
			hLeft = 0
		} else if hLeft > 4 {
			hLeft = 4
		}
		if hRight < 0 {
			hRight = 0
		} else if hRight > 4 {
			hRight = 4
		}

		char := chars.BrailleGrid[hRight][hLeft]
		graphRunes = append(graphRunes, char)
	}

	if len(graphRunes) > availWidth {
		graphRunes = graphRunes[len(graphRunes)-availWidth:]
	}

	centerY := w.Y + (w.Height / 2)

	startXText := w.X + w.Width - 1 - textLen
	for i, r := range textRunes {
		x := startXText + i
		if x > w.X && x < w.X+w.Width-1 && centerY >= 0 && centerY < e.Height {
			e.Buffer[centerY][x] = r
		}
	}

	for i := 0; i < len(graphRunes); i++ {
		runeToPlace := graphRunes[len(graphRunes)-1-i]

		x := startXText - 1 - i

		if x > w.X && x < w.X+w.Width-1 && centerY >= 0 && centerY < e.Height {
			e.Buffer[centerY][x] = runeToPlace
		}
	}
}
