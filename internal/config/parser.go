package config

import (
	"fmt"
	"pulse/internal/widgets"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Theme    string
	Dashname string

	Widgets []widgets.Widget
}

func Parser(config string) (Config, error) {
	var cfg Config
	_, err := toml.Decode(config, &cfg)
	if err != nil {
		fmt.Println("Error: Failed to parse config file")
		return Config{}, err
	}

}
