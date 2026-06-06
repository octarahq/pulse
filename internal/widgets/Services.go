package widgets

import (
	"fmt"
	"pulse/internal/grid"
	"pulse/internal/utils/services"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/docker/docker/client"
)

type ServicesWidget struct {
	BaseWidget
	Manager      string   `toml:"manager"`
	Lines        []string `toml:"lines"`
	dockerClient *client.Client
}

func init() {
	Register("services", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w ServicesWidget
		err := meta.PrimitiveDecode(prim, &w)
		return &w, err
	})
}

func (w *ServicesWidget) Render(e *grid.Engine) {
	if w.Title == "" {
		w.Title = "Services monitor"
	}
	e.DrawBoxTitle(w.X, w.Y, w.Width, w.Height, w.Title)

	if len(w.Lines) == 0 {
		e.DrawText(w.X, w.Y+1, w.Width, 1, "No Data available")
		return
	}

	if w.Manager == "docker" && w.dockerClient == nil {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err == nil {
			w.dockerClient = cli
		}
	} else if w.Manager != "docker" && w.dockerClient == nil {
		for _, l := range w.Lines {
			if strings.HasPrefix(l, "docker") {
				cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err == nil {
					w.dockerClient = cli
				}
				break
			}
		}
	}

	var lines []string

	for _, line := range w.Lines {
		rawLine := line
		if w.Manager == "docker" && !strings.Contains(line, ":") {
			line = "docker:" + line
		}

		args := strings.SplitN(line, ":", 2)
		cmd := args[0]
		target := ""
		if len(args) > 1 {
			target = args[1]
		} else {
			target = rawLine
		}

		switch cmd {
		case "docker":
			if w.dockerClient == nil {
				lines = append(lines, formatServiceLine(target, "[ ERR ]", w.Width))
				break
			}
			val, _ := services.GetDockerContainerState(w.dockerClient, target)
			lines = append(lines, formatServiceLine(target, val, w.Width))
			
		case "docker.health":
			if w.dockerClient == nil {
				lines = append(lines, formatServiceLine(target, "[ ERR ]", w.Width))
				break
			}
			val, _ := services.GetDockerContainerHealth(w.dockerClient, target)
			lines = append(lines, formatServiceLine(target, val, w.Width))
			
		case "docker.restart":
			if w.dockerClient == nil {
				lines = append(lines, formatServiceLine(target, "[ ERR ]", w.Width))
				break
			}
			val, _ := services.GetDockerContainerRestartedTime(w.dockerClient, target)
			lines = append(lines, formatServiceLine(target, val, w.Width))
			
		case "docker.total":
			if w.dockerClient == nil {
				lines = append(lines, "Containers: [ ERR ]")
				break
			}
			val1, val2, err := services.GetDockerTotal(w.dockerClient)
			if err != nil {
				lines = append(lines, "Containers: [ ERR ]")
				break
			}
			lines = append(lines, fmt.Sprintf("Containers: %d up / %d down", val1, val2))

		case "systemd":
			val, _ := services.GetSysStatus(target)
			lines = append(lines, formatServiceLine(target, val, w.Width))
			
		case "systemd.cron":
			val, _ := services.GetSysCronStatus(target)
			lines = append(lines, formatServiceLine(target, val, w.Width))
			
		case "proc.name":
			val, _ := services.GetProcessByName(target)
			lines = append(lines, formatServiceLine(target, val, w.Width))
			
		case "proc.pid":
			val, _ := services.GetProcessByPID(target)
			lines = append(lines, formatServiceLine(target, val, w.Width))
			
		default:
			lines = append(lines, formatServiceLine(target, "[ UNKNOWN CMD ]", w.Width))
		}
	}
	
	for i, l := range lines {
		e.DrawTextDontCenter(w.X, w.Y+i+1, w.Width, w.Height, " "+l)
	}
}

func formatServiceLine(name, status string, width int) string {
	availWidth := width - 4
	if availWidth <= 0 {
		return name + " " + status
	}
	
	spacesCount := availWidth - len([]rune(name)) - len([]rune(status))
	if spacesCount < 1 {
		spacesCount = 1
	}
	
	spaces := strings.Repeat(" ", spacesCount)
	return name + spaces + status
}
