package monitor

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/sensors"
)

func GetGlobalCpuPercent() (float64, error) {
	vm, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}

	return vm[0], nil
}

func GetCpuCorePercent(core int) (float64, error) {
	vm, err := cpu.Percent(0, true)
	if err != nil {
		return 0, err
	}

	if core < 0 || len(vm) < core {
		return 0, fmt.Errorf("invalid core index : %i", core)
	}

	return vm[core], nil
}

func GetCpuTemp() (float64, error) {
	temps, err := sensors.SensorsTemperatures()
	if err != nil {
		return 0, err
	}

	for _, temp := range temps {
		if temp.Temperature > 0 {
			return temp.Temperature, nil
		}
	}

	return 0, fmt.Errorf("no temperature found")
}

func GetCpuFrequencyGHz() (float64, error) {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		return 0, err
	}

	value := strings.TrimSpace(string(data))

	freqKHz, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return freqKHz / 1_000_000, nil
}

func readCpuEnergy() (float64, error) {
	data, err := os.ReadFile("/sys/class/powercap/intel-rapl:0/energy_uj")
	if err != nil {
		return 0, err
	}

	value := strings.TrimSpace(string(data))
	energy, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return energy, nil
}

func GetCpuPowerWatts() (float64, error) {
	startEnergy, err := readCpuEnergy()
	if err != nil {
		return 0, err
	}

	startTime := time.Now()
	time.Sleep(200 * time.Millisecond)

	endEnergy, err := readCpuEnergy()
	if err != nil {
		return 0, err
	}

	elapsed := time.Since(startTime).Seconds()
	deltaJoules := (endEnergy - startEnergy) / 1_000_000
	watts := deltaJoules / elapsed

	return watts, nil
}

func GetCpuStates() (user, system, idle, iowait float64, error error) {
	times, err := cpu.Times(false)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	t := times[0]

	total := t.User + t.System + t.Idle + t.Iowait +
		t.Nice + t.Irq + t.Softirq + t.Steal

	if total == 0 {
		return 0, 0, 0, 0, nil
	}

	user = t.User / total * 100
	system = t.System / total * 100
	idle = t.Idle / total * 100
	iowait = t.Iowait / total * 100

	return
}

func GetCpuStateUser() (float64, error) {
	times, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}

	t := times[0]
	total := t.User + t.System + t.Idle + t.Iowait +
		t.Nice + t.Irq + t.Softirq + t.Steal
	if total == 0 {
		return 0, nil
	}

	return t.User / total, nil
}

func GetCpuStateSystem() (float64, error) {
	times, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}

	t := times[0]
	total := t.User + t.System + t.Idle + t.Iowait +
		t.Nice + t.Irq + t.Softirq + t.Steal
	if total == 0 {
		return 0, nil
	}

	return t.System / total, nil
}

func GetCpuStateIdle() (float64, error) {
	times, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}

	t := times[0]
	total := t.User + t.System + t.Idle + t.Iowait +
		t.Nice + t.Irq + t.Softirq + t.Steal
	if total == 0 {
		return 0, nil
	}

	return t.Idle / total, nil
}

func GetCpuStateIowait() (float64, error) {
	times, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}

	t := times[0]
	total := t.User + t.System + t.Idle + t.Iowait +
		t.Nice + t.Irq + t.Softirq + t.Steal
	if total == 0 {
		return 0, nil
	}

	return t.Iowait / total, nil
}
