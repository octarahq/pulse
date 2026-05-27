package main

import (
	"fmt"
	"os"
	"pulse/internal/config"
	"pulse/internal/display"
	"time"
)

func main() {
	fmt.Print("\033[H\033[2J")

	path := "config_template.toml"
	file, _ := os.ReadFile(path)
	cfg, _ := config.Parser(string(file))

	var minRefresh int = 0
	for _, row := range cfg.Rows {
		for _, wid := range row {
			minRefresh = min(minRefresh, wid.GetBase().Refresh)
		}
	}

	for {
		display.RenderDisplay(cfg)

		time.Sleep(time.Duration(minRefresh) * time.Second)
	}
}
