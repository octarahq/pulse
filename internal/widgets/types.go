package widgets

type BaseWidget struct {
	Type string
}

type DisplayWidjet struct {
	*BaseWidget
	Value string
}
