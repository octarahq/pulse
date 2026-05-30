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

func MakeGauge(prefix string, value float64, max float64, textLabel string, width int) string {
	textRunes := []rune(textLabel)
	prefixRunes := []rune(prefix)

	availWidth := width - len(textRunes) - len(prefixRunes) - 1 // 1 for space

	if availWidth <= 0 {
		return prefix + textLabel
	}

	ratio := 0.0
	if max > 0 {
		ratio = float64(value) / float64(max)
		if ratio > 1.0 {
			ratio = 1.0
		}
		if ratio < 0.0 {
			ratio = 0.0
		}
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

	return fmt.Sprintf("%s%s %s", prefix, string(gaugeRunes), textLabel)
}

func MakeGraph(prefix string, stats []int, textLabel string, width int) string {
	textRunes := []rune(textLabel)
	prefixRunes := []rune(prefix)

	availWidth := width - len(textRunes) - len(prefixRunes) - 1

	if availWidth <= 0 {
		return prefix + textLabel
	}

	var graphRunes []rune
	for i := 0; i < len(stats)-1 && len(graphRunes) < availWidth; i += 2 {
		hLeft := stats[i]
		hRight := stats[i+1]

		if hLeft <= 0 {
			hLeft = 1
		} else if hLeft > 4 {
			hLeft = 4
		}
		if hRight <= 0 {
			hRight = 1
		} else if hRight > 4 {
			hRight = 4
		}

		char := chars.BrailleGrid[hLeft][hRight]
		graphRunes = append(graphRunes, char)
	}

	for len(graphRunes) < availWidth {
		graphRunes = append([]rune{chars.BrailleGrid[1][1]}, graphRunes...)
	}
	if len(graphRunes) > availWidth {
		graphRunes = graphRunes[len(graphRunes)-availWidth:]
	}

	return fmt.Sprintf("%s%s %s", prefix, string(graphRunes), textLabel)
}
