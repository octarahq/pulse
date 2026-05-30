package monitor

import (
	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiskUsed(path string) (float64, error) {
	d, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	return octToGo(float64(d.Used)), nil
}

func GetDiskTotal(path string) (float64, error) {
	d, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	return octToGo(float64(d.Total)), nil
}

func GetDiskFree(path string) (float64, error) {
	d, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	return octToGo(float64(d.Free)), nil
}

func GetDiskUsedPercent(path string) (float64, error) {
	d, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	return d.UsedPercent, nil
}

func GetDiskInodesPercent(path string) (float64, error) {
	d, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}

	return octToGo(d.InodesUsedPercent), nil
}

func GetIoReadMo() (float64, error) {
	d, err := disk.IOCounters()
	if err != nil {
		return 0, err
	}

	var totalReadBytes float64
	for _, stat := range d {
		totalReadBytes += float64(stat.ReadBytes)
	}

	return octToMo(totalReadBytes), nil
}

func GetIoWriteMo() (float64, error) {
	d, err := disk.IOCounters()
	if err != nil {
		return 0, err
	}

	var totalWriteBytes float64
	for _, stat := range d {
		totalWriteBytes += float64(stat.WriteBytes)
	}

	return octToMo(totalWriteBytes), nil
}

func GetIoOps() (int, error) {
	d, err := disk.IOCounters()
	if err != nil {
		return 0, err
	}

	var totalOps uint64
	for _, stat := range d {
		totalOps += stat.ReadCount + stat.WriteCount
	}

	return int(totalOps), nil
}
