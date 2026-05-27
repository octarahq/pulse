package config

import (
	"fmt"
	"pulse/internal/widgets"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Theme    string        `toml:"theme"`
	Dashname string        `toml:"dashname"`
	Rows     []widgets.Row `toml:"rows"`
}

type rowIntermediate struct {
	Widgets []toml.Primitive `toml:"widgets"`
}

type tomlIntermediate struct {
	Theme    string            `toml:"theme"`
	Dashname string            `toml:"dashname"`
	Rows     []rowIntermediate `toml:"rows"`
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
		Rows:     make([]widgets.Row, 0, len(inter.Rows)),
	}

	for _, rowInter := range inter.Rows {
		row := make(widgets.Row, 0, len(rowInter.Widgets))

		for _, prim := range rowInter.Widgets {
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
				row = append(row, w)

			case "clock":
				var w widgets.ClockWidget
				if err := meta.PrimitiveDecode(prim, &w); err != nil {
					return Config{}, fmt.Errorf("fail parsing clock widget: %w", err)
				}
				row = append(row, w)

			default:
				return Config{}, fmt.Errorf("unknown widget type: %s", base.Type)
			}
		}

		cfg.Rows = append(cfg.Rows, row)
	}

	return cfg, nil
}
