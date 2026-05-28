package widgets

import (
	"pulse/internal/grid"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/shirou/gopsutil/v4/mem"
)

type RamSysWidget struct {
	BaseWidget
	Format  string `toml:"format"`
	History []int
}

func init() {
	Register("sys_ram", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w RamSysWidget
		err := meta.PrimitiveDecode(prim, &w)
		if err != nil {
			return nil, err
		}
		w.History = make([]int, 40)

		return &w, nil
	})
}

func (w *RamSysWidget) Render(e *grid.Engine) {
	textLen := 8
	neededStats := (w.Width - 2 - textLen) * 2

	if len(w.History) < neededStats {
		w.History = make([]int, neededStats)
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		DisplayError(e, w.BaseWidget, "Cannot get ram usage")
		return
	}

	percentUsed := vm.UsedPercent

	currentHeight := int((percentUsed / 100.0) * 4.0)
	if currentHeight > 4 {
		currentHeight = 4
	}

	for i := 0; i < len(w.History)-1; i++ {
		w.History[i] = w.History[i+1]
	}
	w.History[len(w.History)-1] = currentHeight

	usedGiB := float64(vm.Used) / (1024 * 1024 * 1024)
	usedStr := strconv.FormatFloat(usedGiB, 'f', 2, 64)

	totalGiB := float64(vm.Total) / (1024 * 1024 * 1024)

	renderType := "gauge"
	if w.Format == "graph" {
		renderType = "graph"
	}

	switch renderType {
	case "graph":
		DisplayGraph(e, w.BaseWidget, w.History, usedStr, "Gi")
	case "gauge":
		DisplayGauge(e, w.BaseWidget, usedGiB, totalGiB, "Gi")
	}
}
