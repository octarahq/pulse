package widgets

import (
	"pulse/internal/grid"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/shirou/gopsutil/v4/cpu"
)

type CpuSysWidget struct {
	BaseWidget
	Format  string `toml:"format"`
	History []int
}

func init() {
	Register("sys_cpu", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w CpuSysWidget
		err := meta.PrimitiveDecode(prim, &w)
		if err != nil {
			return nil, err
		}
		w.History = make([]int, 40)

		return &w, nil
	})
}

func (w *CpuSysWidget) Render(e *grid.Engine) {
	textLen := 8
	neededStats := (w.Width - 2 - textLen) * 2

	if len(w.History) < neededStats {
		w.History = make([]int, neededStats)
	}

	vm, err := cpu.Percent(0, false)
	if err != nil {
		DisplayError(e, w.BaseWidget, "Cannot get cpu usage")
		return
	}

	percentUsed := vm[0]

	var currentHeight int
	switch {
	case percentUsed <= 0.2:
		currentHeight = 0
	case percentUsed <= 2.0:
		currentHeight = 1
	case percentUsed <= 10.0:
		currentHeight = 2
	case percentUsed <= 50.0:
		currentHeight = 3
	default:
		currentHeight = 4
	}

	for i := 0; i < len(w.History)-1; i++ {
		w.History[i] = w.History[i+1]
	}
	w.History[len(w.History)-1] = currentHeight

	renderType := "gauge"
	if w.Format == "graph" {
		renderType = "graph"
	}

	switch renderType {
	case "graph":
		DisplayGraph(e, w.BaseWidget, w.History, strconv.FormatFloat(percentUsed, 'f', 1, 64), "%")
	case "gauge":
		DisplayGauge(e, w.BaseWidget, percentUsed, 100, "%")
	}
}
