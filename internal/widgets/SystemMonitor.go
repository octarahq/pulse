package widgets

import (
	"fmt"
	"pulse/internal/grid"
	"pulse/internal/utils/monitor"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type MonitorSysWidget struct {
	BaseWidget
	Format  string `toml:"format"`
	History map[string]int
	Lines   []string
}

func init() {
	Register("sys_monitor", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w MonitorSysWidget
		err := meta.PrimitiveDecode(prim, &w)
		if err != nil {
			return nil, err
		}
		w.History = map[string]int{}

		return &w, nil
	})
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

			lines = append(lines, fmt.Sprintf("CPU : %.2f%%", value))
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

			lines = append(lines, fmt.Sprintf("CPU Core %d : %.2f%%", coreID, value))

		case args[0] == "cpu.temp":
			value, err := monitor.GetCpuTemp()
			if err != nil {
				lines = append(lines, "Cannot get cpu temp")
				break
			}

			lines = append(lines, fmt.Sprintf("CPU : %d°", int(value)))
		case args[0] == "cpu.freq":
			value, err := monitor.GetCpuFrequencyGHz()
			if err != nil {
				lines = append(lines, "Cannot get cpu freq")
				break
			}

			lines = append(lines, fmt.Sprintf("CPU : %.1f GHz", value))
		case args[0] == "cpu.power":
			value, err := monitor.GetCpuPowerWatts()
			if err != nil {
				lines = append(lines, "Cannot get cpu power usage (need admin)")
				break
			}

			lines = append(lines, fmt.Sprintf("CPU : %.1f Watt", value))

		case args[0] == "cpu.user":
			value, err := monitor.GetCpuStateUser()
			if err != nil {
				lines = append(lines, "Cannot get cpu user")
				break
			}

			lines = append(lines, fmt.Sprintf("CPU : %.2f%%", value))

		case args[0] == "cpu.system":
			value, err := monitor.GetCpuStateSystem()
			if err != nil {
				lines = append(lines, "Cannot get cpu system")
				break
			}

			lines = append(lines, fmt.Sprintf("CPU : %.2f%%", value))

		case args[0] == "cpu.idle":
			value, err := monitor.GetCpuStateIdle()
			if err != nil {
				lines = append(lines, "Cannot get cpu idle")
				break
			}

			lines = append(lines, fmt.Sprintf("CPU : %.2f%%", value))

		case args[0] == "cpu.iowait":
			value, err := monitor.GetCpuStateIowait()
			if err != nil {
				lines = append(lines, "Cannot get cpu iowait")
				break
			}

			lines = append(lines, fmt.Sprintf("CPU : %.2f%%", value))

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

			lines = append(lines, fmt.Sprintf("RAM : %.2f/%.2fGo", value, total))
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

			lines = append(lines, fmt.Sprintf("RAM : %.2f%%", value/total*100))
		case args[0] == "ram.available":
			value, err := monitor.GetRamAvailable()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}

			lines = append(lines, fmt.Sprintf("RAM : %.2fGo", value))
		case args[0] == "ram.cached":
			value, err := monitor.GetRamCached()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}

			lines = append(lines, fmt.Sprintf("RAM : %.2fGo", value))
		case args[0] == "ram.buffers":
			value, err := monitor.GetRamBuffer()
			if err != nil {
				lines = append(lines, "Cannot get ram stats")
				break
			}

			lines = append(lines, fmt.Sprintf("RAM : %.2fGo", value))

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

			lines = append(lines, fmt.Sprintf("Swap : %.2f/%.2fGo", value, total))

		case args[0] == "swap.percent":
			value, err := monitor.GetSwapUsedPercent()
			if err != nil {
				lines = append(lines, "Cannot get swap stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Swap : %.2f%%", value))

		// SYS
		case args[0] == "sys.procs":
			value, err := monitor.GetSysProcs()
			if err != nil {
				lines = append(lines, "Cannot get sys procs stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Procs : %d", value))
		case args[0] == "sys.procs.running":
			value, err := monitor.GetSysProcsRunning()
			if err != nil {
				lines = append(lines, "Cannot get sys procs stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Procs (running) : %d", value))
		case args[0] == "sys.procs.blocked":
			value, err := monitor.GetSysProcsBlocked()
			if err != nil {
				lines = append(lines, "Cannot get sys procs stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Procs (blocked): %d", value))
		case args[0] == "sys.threads":
			value, err := monitor.GetSysThreads()
			if err != nil {
				lines = append(lines, "Cannot get sys thread stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Threads : %d", value))
		case args[0] == "sys.users":
			value, err := monitor.GetSysUsers()
			if err != nil {
				lines = append(lines, "Cannot get sys users stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Users : %d", value))
		case args[0] == "sys.uptime":
			value, err := monitor.GetSysUptime()
			if err != nil {
				lines = append(lines, "Cannot get sys uptime stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Sys : %s", value))
		case args[0] == "sys.ostype":
			value, err := monitor.GetSysOsType()
			if err != nil {
				lines = append(lines, "Cannot get sys stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Sys : %s", value))

		// DISKs
		case strings.Contains(args[0], "disk.percent"):
			path := args[1]
			value, err := monitor.GetDiskUsedPercent(path)
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Disk %s : %.2f%%", path, value))
		case strings.Contains(args[0], "disk.inodes"):
			path := args[1]
			value, err := monitor.GetDiskInodesPercent(path)
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Disk %s : %.2f%%", path, value))
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

			lines = append(lines, fmt.Sprintf("Disk %s : %.2f/%.2fGo", path, used, total))
		case strings.Contains(args[0], "io.read"):
			value, err := monitor.GetIoReadMo()
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Io Read : %.2fMo/s", value))
		case strings.Contains(args[0], "io.write"):
			value, err := monitor.GetIoWriteMo()
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Io Write : %.2fMo/s", value))
		case strings.Contains(args[0], "io.ops"):
			value, err := monitor.GetIoOps()
			if err != nil {
				lines = append(lines, "Cannot get disk stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Io Ops : %dops/s", value))

		}
	}

	for i, l := range lines {
		e.DrawTextDontCenter(w.X, w.Y+i, w.Width, w.Height, fmt.Sprintf(" %s", l))
	}
}
