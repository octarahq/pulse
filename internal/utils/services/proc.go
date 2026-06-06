package services

import (
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func GetProcessByName(name string) (string, error) {
	procs, err := process.Processes()
	if err != nil {
		return "[ MISSING ]", err
	}

	for _, p := range procs {
		pName, err := p.Name()
		if err == nil && strings.Contains(strings.ToLower(pName), strings.ToLower(name)) {
			return "[ RUNNING ]", nil
		}
	}
	return "[ MISSING ]", nil
}

func GetProcessByPID(pidStr string) (string, error) {
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		return "[ MISSING ]", err
	}

	exists, err := process.PidExists(int32(pid))
	if err != nil || !exists {
		return "[ MISSING ]", err
	}

	return "[ RUNNING ]", nil
}
