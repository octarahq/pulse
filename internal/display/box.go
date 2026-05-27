package display

import (
	"strings"
)

func DrawBox(width int, height int, hasBorder bool, content []string) []string {
	lines := make([]string, height)

	topLeft, topRight := "┌", "┐"
	bottomLeft, bottomRight := "└", "┘"
	horiz, vert := "─", "│"

	innerWidth := width
	if hasBorder {
		innerWidth = width - 2
	}

	for y := 0; y < height; y++ {
		var line strings.Builder

		if hasBorder {
			if y == 0 {
				line.WriteString(topLeft)
				line.WriteString(strings.Repeat(horiz, innerWidth))
				line.WriteString(topRight)
			} else if y == height-1 {
				line.WriteString(bottomLeft)
				line.WriteString(strings.Repeat(horiz, innerWidth))
				line.WriteString(bottomRight)
			} else {
				line.WriteString(vert)

				contentIdx := y - 1
				text := ""
				if contentIdx < len(content) {
					text = content[contentIdx]
				}

				runes := []rune(text)
				if len(runes) > innerWidth {
					line.WriteString(string(runes[:innerWidth]))
				} else {
					line.WriteString(text)
					line.WriteString(strings.Repeat(" ", innerWidth-len(runes)))
				}

				line.WriteString(vert)
			}
		} else {
			text := ""
			if y < len(content) {
				text = content[y]
			}
			runes := []rune(text)
			if len(runes) > width {
				line.WriteString(string(runes[:width]))
			} else {
				line.WriteString(text)
				line.WriteString(strings.Repeat(" ", width-len(runes)))
			}
		}

		lines[y] = line.String()
	}

	return lines
}
