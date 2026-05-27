package widgets

type Widget interface {
	GetBase() BaseWidget
}

type Row []Widget

type BaseWidget struct {
	Id    int    `toml:"id"`
	Type  string `toml:"type"`
	Width int    `toml:"width"`
}

func (b BaseWidget) GetBase() BaseWidget {
	return b
}

type DisplayWidget struct {
	BaseWidget
	Value string `toml:"value"`
}

type ClockWidget struct {
	BaseWidget
	Timezone string `toml:"timezone"`
}
