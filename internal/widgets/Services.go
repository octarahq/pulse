package widgets

import (
	"fmt"
	"pulse/internal/grid"
	"pulse/internal/utils/services"
	"strings"

	"github.com/BurntSushi/toml"
)

type ServicesWidget struct {
	BaseWidget
	Lines []string `toml:"lines"`
}

func init() {
	Register("services", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w ServicesWidget
		err := meta.PrimitiveDecode(prim, &w)
		return w, err
	})
}

func (w ServicesWidget) Render(e *grid.Engine) {
	if w.Title == "" {
		w.Title = "Services monitor"
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
		//Docker
		case args[0] == "docker":
			if len(args) < 2 {
				lines = append(lines, "No containers are mentioned.")
				break
			}
			val, err := services.GetDockerContainerState(args[1])
			if err != nil {
				lines = append(lines, "(need admin) Cannot get docker stats")
				break
			}

			lines = append(lines, fmt.Sprintf("%s is %s", args[1], val))
		case args[0] == "docker.uptime":
			if len(args) < 2 {
				lines = append(lines, "No containers are mentioned.")
				break
			}
			val, err := services.GetDockerContainerUptime(args[1])
			if err != nil {
				lines = append(lines, "(need admin) Cannot get docker stats")
				break
			}

			lines = append(lines, fmt.Sprintf("%s is %s", args[1], val))
		case args[0] == "docker.restart":
			if len(args) < 2 {
				lines = append(lines, "No containers are mentioned.")
				break
			}
			val, err := services.GetDockerContainerRestartedTime(args[1])
			if err != nil {
				lines = append(lines, "(need admin) Cannot get docker stats")
				break
			}

			lines = append(lines, fmt.Sprintf("%s as restarted %d times", args[1], val))
		case args[0] == "docker.restartauto":
			if len(args) < 2 {
				lines = append(lines, "No containers are mentioned.")
				break
			}
			val, err := services.GetDockerContainerAutoRestartedTime(args[1])
			if err != nil {
				lines = append(lines, "(need admin) Cannot get docker stats")
				break
			}

			lines = append(lines, fmt.Sprintf("%s as auto-restarted %d times", args[1], val))
		case args[0] == "docker.total":
			val1, val2, err := services.GetDockerTotal()
			if err != nil {
				lines = append(lines, "(need admin) Cannot get docker stats")
				break
			}

			lines = append(lines, fmt.Sprintf("Containers: %d up, %d down", val1, val2))

		// SystemD
		case args[0] == "systemd":
			if len(args) < 2 {
				lines = append(lines, "No jobs are mentioned.")
				break
			}
			val, err := services.GetSysStatus(args[1])
			if err != nil {
				lines = append(lines, "Cannot get systemd stats")
				break
			}

			lines = append(lines, fmt.Sprintf("%s is %s", args[1], val))
		case args[0] == "systemd.cron":
			if len(args) < 2 {
				lines = append(lines, "No jobs are mentioned.")
				break
			}
			val, err := services.GetSysCronStatus(args[1])
			if err != nil {
				lines = append(lines, "Cannot get systemd stats")
				break
			}

			lines = append(lines, fmt.Sprintf("%s cron is %s", args[1], val))
		}
	}
	for i, l := range lines {
		e.DrawTextDontCenter(w.X, w.Y+i, w.Width, w.Height, fmt.Sprintf(" %s", l))
	}
}
