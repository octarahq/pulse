package widgets

import "strings"

func Display(widget DisplayWidget, width int, height int) []string {
	innerWidth := width - 2

	text := widget.Value
	runes := []rune(text)

	if len(runes) > innerWidth {
		text = string(runes[:innerWidth])
		return []string{text}
	}

	totalSpaces := innerWidth - len(runes)
	spacesLeft := totalSpaces / 2
	spacesRight := totalSpaces - spacesLeft

	centeredText := strings.Repeat(" ", spacesLeft) + text + strings.Repeat(" ", spacesRight)
	return []string{centeredText}
}
