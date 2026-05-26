package config

import "pulse/internal/widgets"

type Config struct {
	Theme    string
	Dashname string

	Widgets []widgets.BaseWidget
}

func Parser(config string) {

}
