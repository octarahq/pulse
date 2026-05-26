package main

import (
	"os"
	"pulse/internal/config"
)

func main() {
	path := "config_template.toml"
	file, _ := os.ReadFile(path)
	config.Parser(string(file))
}
