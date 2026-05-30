package monitor

import (
	"github.com/shirou/gopsutil/v4/mem"
)

func GetRamUsed() (float64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return vm.UsedPercent, nil
}

func GetRamTotal() (float64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return octToGo(float64(vm.Total)), nil
}

func GetRamAvailable() (float64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return octToGo(float64(vm.Available)), nil
}

func GetRamCached() (float64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return octToGo(float64(vm.Cached)), nil
}

func GetRamBuffer() (float64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return octToGo(float64(vm.Buffers)), nil
}

func GetSwapUsed() (float64, error) {
	vm, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}

	return octToGo(float64(vm.Used)), nil
}

func GetSwapTotal() (float64, error) {
	vm, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}

	return octToGo(float64(vm.Total)), nil
}

func GetSwapUsedPercent() (float64, error) {
	vm, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}

	return octToGo(vm.UsedPercent), nil
}
