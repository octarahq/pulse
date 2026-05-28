package grid

import (
	"strings"
)

type Engine struct {
	Width  int
	Height int
	Buffer [][]rune
}

func NewEngine(width, height int) *Engine {
	buf := make([][]rune, height)
	for i := range buf {
		buf[i] = make([]rune, width)
		for j := range buf[i] {
			buf[i][j] = ' '
		}
	}
	return &Engine{
		Width:  width,
		Height: height,
		Buffer: buf,
	}
}

func (e *Engine) Box(widgetX, widgetY, widgetWidth, widgetHeight int, title string) {
	if widgetWidth == -1 {
		widgetWidth = e.Width - widgetX
	}
	if widgetHeight == -1 {
		widgetHeight = e.Height - widgetY
	}

	for i := 0; i < widgetWidth; i++ {
		for j := 0; j < widgetHeight; j++ {
			x := widgetX + i
			y := widgetY + j

			if x < 0 || x >= e.Width || y < 0 || y >= e.Height {
				continue
			}

			var char rune
			switch {
			case i == 0 && j == 0:
				char = '┌'
			case i == widgetWidth-1 && j == 0:
				char = '┐'
			case i == 0 && j == widgetHeight-1:
				char = '└'
			case i == widgetWidth-1 && j == widgetHeight-1:
				char = '┘'
			case i == 0 || i == widgetWidth-1:
				char = '│'
			case j == 0:
				if title != "" {
					titleText := " " + title + " "
					titleStart := 2
					titleEnd := titleStart + len([]rune(titleText))

					if i >= titleStart && i < titleEnd {
						char = []rune(titleText)[i-titleStart]
					} else {
						char = '─'
					}
				} else {
					char = '─'
				}

			case j == widgetHeight-1:
				char = '─'
			default:
				continue
			}

			e.Buffer[y][x] = e.merge(e.Buffer[y][x], char)
		}
	}
}

func (e *Engine) DrawBox(widgetX, widgetY, widgetWidth, widgetHeight int) {
	e.Box(widgetX, widgetY, widgetWidth, widgetHeight, "")
}

func (e *Engine) DrawBoxTitle(widgetX, widgetY, widgetWidth, widgetHeight int, title string) {
	e.Box(widgetX, widgetY, widgetWidth, widgetHeight, title)
}

func (e *Engine) DrawText(widgetX, widgetY, widgetWidth, widgetHeight int, text string) {
	if widgetWidth == -1 {
		widgetWidth = e.Width - widgetX
	}
	if widgetHeight == -1 {
		widgetHeight = e.Height - widgetY
	}

	runes := []rune(text)

	availWidth := widgetWidth - 2
	if availWidth <= 0 || widgetHeight <= 2 {
		return
	}

	if len(runes) > availWidth {
		runes = runes[:availWidth]
	}
	textLen := len(runes)

	startX := widgetX + 1 + (availWidth-textLen)/2
	centerY := widgetY + widgetHeight/2

	for i, r := range runes {
		x := startX + i
		if x > widgetX && x < widgetX+widgetWidth-1 && centerY > widgetY && centerY < widgetY+widgetHeight-1 {
			if x >= 0 && x < e.Width && centerY >= 0 && centerY < e.Height {
				e.Buffer[centerY][x] = r
			}
		}
	}
}

func (e *Engine) merge(old, new rune) rune {
	if old == ' ' || old == new {
		return new
	}

	pair := func(r1, r2 rune) bool {
		return (old == r1 && new == r2) || (old == r2 && new == r1)
	}

	switch {
	case pair('┐', '┌'):
		return '┬'
	case pair('┘', '└'):
		return '┴'
	case pair('│', '─'):
		return '┼'
	case pair('┌', '└'):
		return '├'
	case pair('┐', '┘'):
		return '┤'
	case pair('├', '┤'), pair('┬', '┴'):
		return '┼'
	}

	return old
}

func (e *Engine) Flush() string {
	var builder strings.Builder
	builder.WriteString("\033[H")

	for y := 0; y < e.Height; y++ {
		builder.WriteString(string(e.Buffer[y]))
		if y < e.Height-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
