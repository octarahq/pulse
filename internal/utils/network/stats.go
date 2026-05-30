package network

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	psnet "github.com/shirou/gopsutil/v4/net"
)

func GetNetStatus(interName string) (map[string]bool, error) {
	ifaces, err := psnet.Interfaces()
	if err != nil {
		return nil, err
	}
	res := make(map[string]bool)
	isUp := func(flags []string) bool {
		for _, f := range flags {
			if f == "up" || f == "UP" {
				return true
			}
		}
		return false
	}
	if interName == "all" {
		for _, iface := range ifaces {
			res[iface.Name] = isUp(iface.Flags)
		}
		return res, nil
	}
	for _, iface := range ifaces {
		if iface.Name == interName {
			res[iface.Name] = isUp(iface.Flags)
			return res, nil
		}
	}
	return res, nil
}

func GetNetStatusError(interName string) (int, error) {
	stats, err := psnet.IOCounters(true)
	if err != nil {
		return 0, err
	}

	for _, stat := range stats {
		if interName == "all" || stat.Name == interName {
			return int(stat.Errin + stat.Errout), nil
		}
	}
	return 0, nil
}

func GetNetStatusDroped(interName string) (int, error) {
	stats, err := psnet.IOCounters(true)
	if err != nil {
		return 0, err
	}

	for _, stat := range stats {
		if interName == "all" || stat.Name == interName {
			return int(stat.Dropin + stat.Dropout), nil
		}
	}
	return 0, nil
}

func GetNetStatusWifi() (int, error) {
	f, err := os.Open("/proc/net/wireless")
	if err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for i := 0; scanner.Scan(); i++ {
			if i < 2 {
				continue
			}
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				levelStr := strings.TrimSuffix(fields[3], ".")
				levelStr = strings.TrimSpace(levelStr)
				if levelStr != "" {
					if v, err := strconv.ParseFloat(levelStr, 64); err == nil {
						return int(v), nil
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			_ = err
		}
	}

	if out, err := exec.Command("iwconfig").Output(); err == nil {
		re := regexp.MustCompile(`Signal level[:=]\s*(-?\d+)\s*d?B?m?`)
		if m := re.FindSubmatch(out); len(m) >= 2 {
			if v, err := strconv.Atoi(string(m[1])); err == nil {
				return v, nil
			}
		}
	}

	return 0, errors.New("wifi signal not found")
}
