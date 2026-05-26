package widgets

type Widget struct {
	BaseWidget
	DisplayWidget
}

type BaseWidget struct {
	Type string
}

type DisplayWidget struct {
	Value string
}
