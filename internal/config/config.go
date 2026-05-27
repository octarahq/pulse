package config

import (
	"fmt"
	"pulse/internal/widgets"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Theme    string           `toml:"theme"`
	Dashname string           `toml:"dashname"`
	Widgets  []widgets.Widget `toml:"-"`
}

type tomlIntermediate struct {
	Theme    string           `toml:"theme"`
	Dashname string           `toml:"dashname"`
	Widgets  []toml.Primitive `toml:"widgets"`
}

func Parser(tomlContent string) (Config, error) {
	var inter tomlIntermediate
	meta, err := toml.Decode(tomlContent, &inter)
	if err != nil {
		return Config{}, fmt.Errorf("fail parsing global toml: %w", err)
	}

	cfg := Config{
		Theme:    inter.Theme,
		Dashname: inter.Dashname,
		Widgets:  make([]widgets.Widget, 0, len(inter.Widgets)),
	}

	for _, prim := range inter.Widgets {
		var base widgets.BaseWidget
		if err := meta.PrimitiveDecode(prim, &base); err != nil {
			return Config{}, fmt.Errorf("fail parsing base widget: %w", err)
		}

		switch base.Type {
		case "display":
			var w widgets.DisplayWidget
			if err := meta.PrimitiveDecode(prim, &w); err != nil {
				return Config{}, fmt.Errorf("fail parsing display widget: %w", err)
			}
			cfg.Widgets = append(cfg.Widgets, w)

		case "clock":
			var w widgets.ClockWidget
			if err := meta.PrimitiveDecode(prim, &w); err != nil {
				return Config{}, fmt.Errorf("fail parsing clock widget: %w", err)
			}
			cfg.Widgets = append(cfg.Widgets, w)

		default:
			return Config{}, fmt.Errorf("unknown widget type: %s", base.Type)
		}
	}

	return cfg, nil
}
