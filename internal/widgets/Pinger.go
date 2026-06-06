package widgets

import (
	"pulse/internal/grid"
	"pulse/internal/utils/network"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type PingerWidget struct {
	BaseWidget
	Targets []string `toml:"targets"`
	results map[string]string
	mutex   sync.Mutex
	lastRun time.Time
}

func init() {
	Register("pinger", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w PingerWidget
		err := meta.PrimitiveDecode(prim, &w)
		w.results = make(map[string]string)
		return &w, err
	})
}

func (w *PingerWidget) Render(e *grid.Engine) {
	if w.Title == "" {
		w.Title = "Pinger"
	}
	e.DrawBoxTitle(w.X, w.Y, w.Width, w.Height, w.Title)

	if len(w.Targets) == 0 {
		e.DrawText(w.X, w.Y+1, w.Width, 1, "No targets defined")
		return
	}

	w.mutex.Lock()
	if time.Since(w.lastRun) > 5*time.Second {
		w.lastRun = time.Now()
		go w.pingAll()
	}
	w.mutex.Unlock()

	for i, targetLine := range w.Targets {
		parts := strings.SplitN(targetLine, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		target := strings.TrimSpace(parts[1])

		w.mutex.Lock()
		res, ok := w.results[target]
		w.mutex.Unlock()

		if !ok {
			res = "Pinging..."
		}

		line := formatServiceLine(name, res, w.Width)
		e.DrawTextDontCenter(w.X, w.Y+i+1, w.Width, w.Height, " "+line)
	}
}

func (w *PingerWidget) pingAll() {
	for _, targetLine := range w.Targets {
		parts := strings.SplitN(targetLine, "=", 2)
		if len(parts) != 2 {
			continue
		}
		target := strings.TrimSpace(parts[1])

		res, err := network.PingTarget(target)

		w.mutex.Lock()
		if err != nil {
			w.results[target] = "[ DOWN ]"
		} else {
			w.results[target] = res
		}
		w.mutex.Unlock()
	}
}
