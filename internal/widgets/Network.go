package widgets

import (
	"fmt"
	"math"
	"pulse/internal/grid"
	"pulse/internal/utils/network"
	"strings"

	"github.com/BurntSushi/toml"
)

type NetworkWidget struct {
	BaseWidget
	Format  string `toml:"format"`
	History map[string][]int
	Lines   []string
}

func init() {
	Register("network", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w NetworkWidget
		err := meta.PrimitiveDecode(prim, &w)
		if err != nil {
			return nil, err
		}
		w.History = make(map[string][]int)

		return &w, nil
	})
}

func (w *NetworkWidget) Render(e *grid.Engine) {
	if w.Title == "" {
		w.Title = "Network monitor"
	}
	e.DrawBoxTitle(w.X, w.Y, w.Width, w.Height, w.Title)

	if len(w.Lines) == 0 {
		e.DrawText(w.X, w.Y+1, w.Width, 1, "No Data available")
		return
	}

	var lines []string

	for _, line := range w.Lines {
		args := strings.Split(line, ":")

		switch {
		case strings.HasPrefix(args[0], "net.in"):
			if len(args) < 2 { continue }
			iface := args[1]
			val, _ := network.GetNetIoIn(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s in:", shortIface(iface)), float64(val), 100.0, fmt.Sprintf(" %.2f MB/s", val), w.Width))

		case strings.HasPrefix(args[0], "net.out"):
			if len(args) < 2 { continue }
			iface := args[1]
			val, _ := network.GetNetIoOut(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s out:", shortIface(iface)), float64(val), 100.0, fmt.Sprintf(" %.2f MB/s", val), w.Width))

		case strings.HasPrefix(args[0], "net.total"):
			if len(args) < 2 { continue }
			iface := args[1]
			val, _ := network.GetNetIoTotal(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s tot:", shortIface(iface)), float64(val), 100.0, fmt.Sprintf(" %.2f MB/s", val), w.Width))

		case strings.HasPrefix(args[0], "net.speed"):
			if len(args) < 2 { continue }
			iface := args[1]
			in, out, _ := network.GetNetIoSpeed(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s:", shortIface(iface)), 0, 0, fmt.Sprintf(" ▼ %.1f MB/s ▲ %.1f MB/s", in, out), w.Width))

		case strings.HasPrefix(args[0], "net.bytes.in"):
			if len(args) < 2 { continue }
			iface := args[1]
			val, _ := network.GetNetIoBytesIn(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s in:", shortIface(iface)), float64(val), 1000.0, fmt.Sprintf(" %.2f GB", val), w.Width))

		case strings.HasPrefix(args[0], "net.bytes.out"):
			if len(args) < 2 { continue }
			iface := args[1]
			val, _ := network.GetNetIoBytesOut(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s out:", shortIface(iface)), float64(val), 1000.0, fmt.Sprintf(" %.2f GB", val), w.Width))

		case strings.HasPrefix(args[0], "net.status"):
			if len(args) < 2 { continue }
			iface := args[1]
			status, _ := network.GetNetStatus(iface)
			statusStr := "[ DOWN ]"
			if isUp, ok := status[iface]; ok && isUp {
				statusStr = "[ UP ]"
			}
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s:", shortIface(iface)), 0, 0, fmt.Sprintf(" %s", statusStr), w.Width))

		case strings.HasPrefix(args[0], "net.errors"):
			if len(args) < 2 { continue }
			iface := args[1]
			val, _ := network.GetNetStatusError(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s err:", shortIface(iface)), float64(val), 1000.0, fmt.Sprintf(" %d", val), w.Width))

		case strings.HasPrefix(args[0], "net.dropped"):
			if len(args) < 2 { continue }
			iface := args[1]
			val, _ := network.GetNetStatusDroped(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s drp:", shortIface(iface)), float64(val), 1000.0, fmt.Sprintf(" %d", val), w.Width))

		case args[0] == "net.wifi.signal":
			val, _ := network.GetNetStatusWifi()
			lines = append(lines, w.formatLine(line, "Wi-Fi:", float64(val), 100.0, fmt.Sprintf(" %d%%", val), w.Width))

		case strings.HasPrefix(args[0], "ip.local"):
			if len(args) < 2 { continue }
			iface := args[1]
			val := network.GetNetIpLocal(iface)
			lines = append(lines, w.formatLine(line, fmt.Sprintf("%s IP:", shortIface(iface)), 0, 0, fmt.Sprintf(" %s", val), w.Width))

		case args[0] == "ip.public":
			val := network.GetNetIpPublic()
			lines = append(lines, w.formatLine(line, "Pub IP:", 0, 0, fmt.Sprintf(" %s", val), w.Width))

		case args[0] == "ip.gateway":
			val := network.GetNetIpGateway()
			lines = append(lines, w.formatLine(line, "Gate:", 0, 0, fmt.Sprintf(" %s", val), w.Width))

		case args[0] == "net.connections":
			val := network.GetNetConnOpen()
			lines = append(lines, w.formatLine(line, "Conns:", float64(val), 1000.0, fmt.Sprintf(" %d", val), w.Width))

		case args[0] == "net.tcp":
			val := network.GetNetConnTcp()
			lines = append(lines, w.formatLine(line, "TCP:", float64(val), 1000.0, fmt.Sprintf(" %d", val), w.Width))

		case args[0] == "net.udp":
			val := network.GetNetConnUdp()
			lines = append(lines, w.formatLine(line, "UDP:", float64(val), 1000.0, fmt.Sprintf(" %d", val), w.Width))

		case args[0] == "net.ssh":
			val := network.GetNetConnSsh()
			lines = append(lines, w.formatLine(line, "SSH:", float64(val), 100.0, fmt.Sprintf(" %d", val), w.Width))

		case args[0] == "net.established":
			val := network.GetNetConnEstablished()
			lines = append(lines, w.formatLine(line, "Estab:", float64(val), 1000.0, fmt.Sprintf(" %d", val), w.Width))
		}
	}
	for i, l := range lines {
		e.DrawTextDontCenter(w.X, w.Y+i, w.Width, w.Height, fmt.Sprintf(" %s", l))
	}
}

func (w *NetworkWidget) formatLine(line, prefix string, value, max float64, textLabel string, width int) string {
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
	if availWidth <= 0 {
		return prefix + textLabel
	}

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

func shortIface(name string) string {
	if strings.HasPrefix(name, "enx") || strings.HasPrefix(name, "eth") || strings.HasPrefix(name, "en") {
		return "eth"
	}
	if strings.HasPrefix(name, "wl") {
		return "wifi"
	}
	if name == "tailscale0" {
		return "ts0"
	}
	if name == "all" {
		return "all"
	}
	if len(name) > 4 {
		return name[:4]
	}
	return name
}
