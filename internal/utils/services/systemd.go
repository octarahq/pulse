package services

import (
	"os/exec"
	"strings"
)

func GetSysStatus(service string) (string, error) {
	out, err := exec.Command("systemctl", "is-active", service).CombinedOutput()
	status := strings.TrimSpace(string(out))

	if status == "active" {
		return "[ACTIVE]", nil
	}
	if err != nil && status == "" {
		return "[UNKNOWN]", err
	}
	return "[INACTIVE]", nil
}

func GetSysCronStatus(service string) (string, error) {
	unit := service
	if !strings.HasSuffix(unit, ".timer") {
		unit = unit + ".timer"
	}

	out, _ := exec.Command("systemctl", "is-active", unit).CombinedOutput()
	status := strings.TrimSpace(string(out))

	if status == "active" {
		return "[ACTIVE]", nil
	}

	out2, err2 := exec.Command("systemctl", "show", unit, "-p", "LastTriggerUSecRealtime", "-p", "NextElapseUSecRealtime").CombinedOutput()
	if err2 != nil {
		return "[INACTIVE]", nil
	}

	var last, next string
	for _, line := range strings.Split(strings.TrimSpace(string(out2)), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], strings.TrimSpace(parts[1])

		if val == "" || val == "0" || strings.ToLower(val) == "n/a" {
			continue
		}

		if key == "LastTriggerUSecRealtime" {
			last = val
		} else if key == "NextElapseUSecRealtime" {
			next = val
		}
	}

	if next != "" {
		return "[NEXT: " + next + "]", nil
	}
	if last != "" {
		return "[LAST: " + last + "]", nil
	}

	return "[INACTIVE]", nil
}
