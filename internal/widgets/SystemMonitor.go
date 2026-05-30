package widgets

import (
	"fmt"
	"math"
	"pulse/internal/grid"
	"pulse/internal/utils/monitor"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type MonitorSysWidget struct {
	BaseWidget
	Format  string `toml:"format"`
	History map[string][]int
	Lines   []string
}

func init() {
	Register("sys_monitor", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w MonitorSysWidget
		err := meta.PrimitiveDecode(prim, &w)
		if err != nil {
			return nil, err
		}
		w.History = map[string][]int{}

		return &w, nil
	})
}


func (w *MonitorSysWidget) formatLine(line, prefix string, value, max float64, textLabel string, width int) string {
	args := strings.Split(line, ":")
	opt := ""
	if len(args) > 1 {
		lastArg := args[len(args)-1]
		if lastArg == "gauge" || lastArg == "graph" {
			opt = lastArg
		}
	}

	if opt == "" {
		return prefix + textLabel
	}

	availWidth := width - 3
	if opt == "gauge" {
		return MakeGauge(prefix, value, max, textLabel, availWidth)
	} else if opt == "graph" {
		key := line
		h := 0
		if max > 0 {
			ratio := value / max
			if ratio > 1.0 {
				ratio = 1.0
			} else if ratio < 0.0 {
				ratio = 0.0
			}
			
			if ratio == 0.0 {
				h = 1
			} else {
				h = 1 + int(math.Ceil(ratio*3.0))
			}
		}
		if h > 4 {
			h = 4
		}
		if h < 0 {
			h = 0
		}
		w.History[key] = append(w.History[key], h)

		graphAvailWidth := availWidth - len([]rune(prefix)) - len([]rune(textLabel)) - 1
		if graphAvailWidth > 0 && len(w.History[key]) > graphAvailWidth*2 {
			w.History[key] = w.History[key][len(w.History[key])-graphAvailWidth*2:]
		}

		return MakeGraph(prefix, w.History[key], textLabel, availWidth)
	}

	return prefix + textLabel
}

func (w *MonitorSysWidget) Render(e *grid.Engine) {

	if w.Title == "" {
		w.Title = "System monitor"
	}
	e.DrawBoxTitle(w.X, w.Y, w.Width, w.Height, w.Title)

	if len(w.Lines) == 0 {
		e.DrawText(w.X, w.Y, w.Width, w.Height, "No Data available")
		return
	}

	var lines []string

	for _, line := range w.Lines {
		args := strings.Split(line, ":")

		switch {
		// CPU
		case args[0] == "cpu":
			value, err := monitor.GetGlobalCpuPercent()
			if err != nil {
				lines = append(lines, "Cannot get cpu usage")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU Usage :", float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))
		case strings.HasPrefix(args[0], "cpu.core"):
			coreStr := strings.TrimPrefix(args[0], "cpu.core")

			coreID, err := strconv.Atoi(coreStr)
			if err != nil {
				lines = append(lines, "Invalid core id")
				break
			}

			value, err := monitor.GetCpuCorePercent(coreID)
			if err != nil {
				lines = append(lines, "Cannot get cpu core usage")
				break
			}

			lines = append(lines, w.formatLine(line, fmt.Sprintf("Core %d Usage :", coreID), float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))

		case args[0] == "cpu.temp":
			value, err := monitor.GetCpuTemp()
			if err != nil {
				lines = append(lines, "Cannot get cpu temp")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU Temp :", float64(value), 100.0, fmt.Sprintf(" %d°", int(value)), w.Width))
		case args[0] == "cpu.freq":
			value, err := monitor.GetCpuFrequencyGHz()
			if err != nil {
				lines = append(lines, "Cannot get cpu freq")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU Freq :", float64(value), 5.0, fmt.Sprintf(" %.1f GHz", value), w.Width))
		case args[0] == "cpu.power":
			value, err := monitor.GetCpuPowerWatts()
			if err != nil {
				lines = append(lines, "Cannot get cpu power usage (need admin)")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU Power :", float64(value), 150.0, fmt.Sprintf(" %.1f Watt", value), w.Width))

		case args[0] == "cpu.user":
			value, err := monitor.GetCpuStateUser()
			if err != nil {
				lines = append(lines, "Cannot get cpu user")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU User :", float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))

		case args[0] == "cpu.system":
			value, err := monitor.GetCpuStateSystem()
			if err != nil {
				lines = append(lines, "Cannot get cpu system")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU System :", float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))

		case args[0] == "cpu.idle":
			value, err := monitor.GetCpuStateIdle()
			if err != nil {
				lines = append(lines, "Cannot get cpu idle")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU Idle :", float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))

		case args[0] == "cpu.iowait":
			value, err := monitor.GetCpuStateIowait()
			if err != nil {
				lines = append(lines, "Cannot get cpu iowait")
				break
			}

			lines = append(lines, w.formatLine(line, "CPU IOWait :", float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))

		// RAM
		case args[0] == "ram":
			value, err := monitor.GetRamAvailable()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}
			total, err := monitor.GetRamTotal()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}

			lines = append(lines, w.formatLine(line, "RAM Usage :", float64(value), float64(total), fmt.Sprintf(" %.2f/%.2f Go", value, total), w.Width))
		case args[0] == "ram.percent":
			value, err := monitor.GetRamAvailable()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}
			total, err := monitor.GetRamTotal()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}

			if total == 0 {
				lines = append(lines, "Cannot get ram stats")
				break
			}

			lines = append(lines, w.formatLine(line, "RAM Usage % :", float64(value/total*100), 100.0, fmt.Sprintf(" %.2f%%", value/total*100), w.Width))
		case args[0] == "ram.available":
			value, err := monitor.GetRamAvailable()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}
			total, _ := monitor.GetRamTotal()

			lines = append(lines, w.formatLine(line, "RAM Avail :", float64(value), float64(total), fmt.Sprintf(" %.2f Go", value), w.Width))
		case args[0] == "ram.cached":
			value, err := monitor.GetRamCached()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}
			total, _ := monitor.GetRamTotal()

			lines = append(lines, w.formatLine(line, "RAM Cached :", float64(value), float64(total), fmt.Sprintf(" %.2f Go", value), w.Width))
		case args[0] == "ram.buffers":
			value, err := monitor.GetRamBuffer()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}
			total, _ := monitor.GetRamTotal()

			lines = append(lines, w.formatLine(line, "RAM Buffers :", float64(value), float64(total), fmt.Sprintf(" %.2f Go", value), w.Width))

		// SWAP
		case args[0] == "swap":
			value, err := monitor.GetSwapUsed()
			if err != nil {
				lines = append(lines, "Cannot get swap stats")
				break
			}
			total, err := monitor.GetSwapTotal()
			if err != nil {
				lines = append(lines, "Cannot get swap stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Swap Usage :", float64(value), float64(total), fmt.Sprintf(" %.2f/%.2f Go", value, total), w.Width))

		case args[0] == "swap.percent":
			value, err := monitor.GetSwapUsedPercent()
			if err != nil {
				lines = append(lines, "Cannot get swap stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Swap Usage % :", float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))

		// SYS
		case args[0] == "sys.procs":
			value, err := monitor.GetSysProcs()
			if err != nil {
				lines = append(lines, "Cannot get sys procs stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Processes :", float64(value), 2000.0, fmt.Sprintf(" %d", value), w.Width))
		case args[0] == "sys.procs.running":
			value, err := monitor.GetSysProcsRunning()
			if err != nil {
				lines = append(lines, "Cannot get sys procs stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Procs Run :", float64(value), 50.0, fmt.Sprintf(" %d", value), w.Width))
		case args[0] == "sys.procs.blocked":
			value, err := monitor.GetSysProcsBlocked()
			if err != nil {
				lines = append(lines, "Cannot get sys procs stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Procs Block :", float64(value), 50.0, fmt.Sprintf(" %d", value), w.Width))
		case args[0] == "sys.threads":
			value, err := monitor.GetSysThreads()
			if err != nil {
				lines = append(lines, "Cannot get sys thread stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Threads :", float64(value), 5000.0, fmt.Sprintf(" %d", value), w.Width))
		case args[0] == "sys.users":
			value, err := monitor.GetSysUsers()
			if err != nil {
				lines = append(lines, "Cannot get sys users stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Users :", float64(value), 10.0, fmt.Sprintf(" %d", value), w.Width))
		case args[0] == "sys.uptime":
			value, err := monitor.GetSysUptime()
			if err != nil {
				lines = append(lines, "Cannot get sys uptime stats")
				break
			}

			lines = append(lines, w.formatLine(line, "Uptime :", 0.0, 0.0, fmt.Sprintf(" %s", value), w.Width))
		case args[0] == "sys.ostype":
			value, err := monitor.GetSysOsType()
			if err != nil {
				lines = append(lines, "Cannot get sys stats")
				break
			}

			lines = append(lines, w.formatLine(line, "OS Type :", 0.0, 0.0, fmt.Sprintf(" %s", value), w.Width))

		// DISKs
		case strings.Contains(args[0], "disk.percent"):
			path := args[1]
			value, err := monitor.GetDiskUsedPercent(path)
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, w.formatLine(line, fmt.Sprintf("Disk %s Usage :", path), float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))
		case strings.Contains(args[0], "disk.inodes"):
			path := args[1]
			value, err := monitor.GetDiskInodesPercent(path)
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, w.formatLine(line, fmt.Sprintf("Disk %s Inodes :", path), float64(value), 100.0, fmt.Sprintf(" %.2f%%", value), w.Width))
		case strings.Contains(args[0], "disk"):
			path := args[1]
			used, err := monitor.GetDiskUsed(path)
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			total, err := monitor.GetDiskTotal(path)
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, w.formatLine(line, fmt.Sprintf("Disk %s Used :", path), float64(used), float64(total), fmt.Sprintf(" %.2f/%.2f Go", used, total), w.Width))
		case strings.Contains(args[0], "io.read"):
			value, err := monitor.GetIoReadMo()
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, w.formatLine(line, "IO Read :", float64(value), 500.0, fmt.Sprintf(" %.2f MB/s", value), w.Width))
		case strings.Contains(args[0], "io.write"):
			value, err := monitor.GetIoWriteMo()
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, w.formatLine(line, "IO Write :", float64(value), 500.0, fmt.Sprintf(" %.2f MB/s", value), w.Width))
		case strings.Contains(args[0], "io.ops"):
			value, err := monitor.GetIoOps()
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, w.formatLine(line, "IO Ops :", float64(value), 5000.0, fmt.Sprintf(" %d ops/s", value), w.Width))

		}
	}

	for i, l := range lines {
		e.DrawTextDontCenter(w.X, w.Y+i, w.Width, w.Height, fmt.Sprintf(" %s", l))
	}
}
