package widgets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pulse/internal/grid"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type ApiWidget struct {
	BaseWidget
	URL       string            `toml:"url"`
	Find      string            `toml:"find"`
	Variables map[string]string `toml:"variables"`
	Label     string            `toml:"label"`
	value     string
	values    map[string]string
	mutex     sync.Mutex
	lastRun   time.Time
}

func init() {
	Register("api", func(prim toml.Primitive, meta *toml.MetaData) (Widget, error) {
		var w ApiWidget
		err := meta.PrimitiveDecode(prim, &w)
		w.values = make(map[string]string)
		return &w, err
	})
}

func (w *ApiWidget) Render(e *grid.Engine) {
	if w.Title == "" {
		w.Title = "API Tracker"
	}
	e.DrawBoxTitle(w.X, w.Y, w.Width, w.Height, w.Title)

	if w.URL == "" || (w.Find == "" && len(w.Variables) == 0) {
		e.DrawText(w.X, w.Y+1, w.Width, 1, "Missing URL or Find/Variables")
		return
	}

	w.mutex.Lock()
	if time.Since(w.lastRun) > 60*time.Second {
		w.lastRun = time.Now()
		go w.fetchAPI()
	}
	val := w.value
	vals := make(map[string]string)
	for k, v := range w.values {
		vals[k] = v
	}
	w.mutex.Unlock()

	displayLabel := w.Label
	if val == "" && len(vals) == 0 {
		displayLabel = "Loading..."
	} else {
		if strings.Contains(displayLabel, "{value}") {
			if val == "" {
				val = "Loading..."
			}
			displayLabel = strings.Replace(displayLabel, "{value}", val, -1)
		}
		
		for k, v := range vals {
			placeholder := "{" + k + "}"
			if v == "" {
				v = "Loading..."
			}
			displayLabel = strings.Replace(displayLabel, placeholder, v, -1)
		}

		if !strings.Contains(w.Label, "{") && w.Find != "" {
			displayLabel = fmt.Sprintf("%s %s", w.Label, val)
		}
	}

	e.DrawText(w.X, w.Y+1, w.Width, w.Height-2, displayLabel)
}

func (w *ApiWidget) fetchAPI() {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(w.URL)
	if err != nil {
		w.mutex.Lock()
		w.value = "Err: HTTP"
		for k := range w.Variables {
			w.values[k] = "Err: HTTP"
		}
		w.mutex.Unlock()
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.mutex.Lock()
		w.value = "Err: Read"
		for k := range w.Variables {
			w.values[k] = "Err: Read"
		}
		w.mutex.Unlock()
		return
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.mutex.Lock()
		w.value = "Err: JSON"
		for k := range w.Variables {
			w.values[k] = "Err: JSON"
		}
		w.mutex.Unlock()
		return
	}

	res := ""
	if w.Find != "" {
		res, err = extractValue(data, w.Find)
		if err != nil {
			res = "Err: Path"
		}
	}

	newValues := make(map[string]string)
	for k, path := range w.Variables {
		varRes, errVar := extractValue(data, path)
		if errVar != nil {
			newValues[k] = "Err: Path"
		} else {
			newValues[k] = varRes
		}
	}

	w.mutex.Lock()
	if w.Find != "" {
		w.value = res
	}
	w.values = newValues
	w.mutex.Unlock()
}

func extractValue(data interface{}, path string) (string, error) {
	keys := strings.Split(path, ".")
	var current interface{} = data

	for _, key := range keys {
		if strings.HasPrefix(key, "[") && strings.HasSuffix(key, "]") {
			indexStr := key[1 : len(key)-1]
			idx, err := strconv.Atoi(indexStr)
			if err != nil {
				return "", fmt.Errorf("invalid index: %s", key)
			}
			arr, ok := current.([]interface{})
			if !ok || idx < 0 || idx >= len(arr) {
				return "", fmt.Errorf("index out of bounds or not an array: %s", key)
			}
			current = arr[idx]
			continue
		}

		m, ok := current.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("not an object at key: %s", key)
		}
		val, exists := m[key]
		if !exists {
			return "", fmt.Errorf("key not found: %s", key)
		}
		current = val
	}

	return fmt.Sprintf("%v", current), nil
}
