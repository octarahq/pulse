package main

import (
	"fmt"
	"os"
	"pulse/internal/config"
	"pulse/internal/grid"
	_ "pulse/internal/widgets"
	"time"
)

func main() {
	fmt.Print("\033[H\033[2J")

	path := "config.toml"
	file, _ := os.ReadFile(path)
	cfg, err := config.Parser(string(file))
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return
	}

	engine := grid.NewEngine(140, 20)

	for {
		for i := range engine.Buffer {
			for j := range engine.Buffer[i] {
				engine.Buffer[i][j] = ' '
			}
		}

		for _, w := range cfg.Widgets {
			w.Render(engine)
		}

		fmt.Print(engine.Flush())

		time.Sleep(3 * time.Second)
	}
}
