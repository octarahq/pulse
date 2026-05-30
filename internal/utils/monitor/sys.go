package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v4/process"
)

func GetSysProcs() (int, error) {
	procs, err := process.Processes()
	if err != nil {
		return 0, err
	}
	return len(procs), nil
}

func GetSysProcsWithState(state string) (int, error) {
	procs, err := process.Processes()
	if err != nil {
		return 0, err
	}

	running := 0
	for _, p := range procs {
		statuses, err := p.Status()
		if err != nil {
			continue
		}
		for _, s := range statuses {
			if s == state {
				running++
				break
			}
		}
	}

	return running, nil
}

func GetSysProcsRunning() (int, error) {
	return GetSysProcsWithState("running")
}

func GetSysProcsBlocked() (int, error) {
	return GetSysProcsWithState("blocked")
}

func GetSysThreads() (int, error) {
	procs, err := process.Processes()
	if err != nil {
		return 0, err
	}

	total := 0
	for _, p := range procs {
		threads, err := p.NumThreads()
		if err != nil {
			continue
		}
		total += int(threads)
	}

	return total, nil
}

func GetSysUsers() (int, error) {
	users, err := host.Users()
	if err != nil {
		return 0, err
	}

	return len(users), nil
}

func GetSysUptime() (string, error) {
	u, err := host.Uptime()
	if err != nil {
		return "", err
	}

	uptime := time.Duration(u) * time.Second
	if uptime >= 24*time.Hour {
		return fmt.Sprintf("%dd", int(uptime.Hours()/24)), nil
	}
	if uptime >= time.Hour {
		return fmt.Sprintf("%dh", int(uptime.Hours())), nil
	}
	if uptime >= time.Minute {
		return fmt.Sprintf("%dm", int(uptime.Minutes())), nil
	}
	return fmt.Sprintf("%ds", int(uptime.Seconds())), nil
}

func GetSysOsType() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}

	distro := info.Platform
	if distro == "" {
		distro = info.OS
	}
	if distro == "" {
		distro = "unknown"
	}

	version := info.PlatformVersion
	if version == "" {
		version = "unknown"
	}

	kernel := info.KernelVersion
	if kernel == "" {
		kernel = "unknown"
	}

	return fmt.Sprintf("%s %s / %s", distro, version, kernel), nil
}
