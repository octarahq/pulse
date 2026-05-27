package main

import (
	"fmt"
	"os"
	"pulse/internal/config"
)

func main() {
	path := "config_template.toml"
	file, _ := os.ReadFile(path)
	cfg, _ := config.Parser(string(file))
	fmt.Printf("%+v\n", cfg)
}
