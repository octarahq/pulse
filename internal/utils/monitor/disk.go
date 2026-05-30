package monitor

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

var (
	ioMu         sync.Mutex
	lastIoTime   time.Time
	lastIoRead   float64
	lastIoWrite  float64
	lastIoOps    uint64
	ioReadSpeed  float64
	ioWriteSpeed float64
	ioOpsSpeed   int
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

func updateIoStats() error {
	ioMu.Lock()
	defer ioMu.Unlock()

	now := time.Now()
	if now.Sub(lastIoTime) < 500*time.Millisecond {
		return nil // already updated recently
	}

	d, err := disk.IOCounters()
	if err != nil {
		return err
	}

	var totalReadBytes, totalWriteBytes float64
	var totalOps uint64
	for _, stat := range d {
		totalReadBytes += float64(stat.ReadBytes)
		totalWriteBytes += float64(stat.WriteBytes)
		totalOps += stat.ReadCount + stat.WriteCount
	}

	if !lastIoTime.IsZero() {
		elapsedSec := now.Sub(lastIoTime).Seconds()
		if elapsedSec > 0 {
			ioReadSpeed = octToMo((totalReadBytes - lastIoRead) / elapsedSec)
			ioWriteSpeed = octToMo((totalWriteBytes - lastIoWrite) / elapsedSec)
			ioOpsSpeed = int(float64(totalOps - lastIoOps) / elapsedSec)
		}
	} else {
		ioReadSpeed = 0
		ioWriteSpeed = 0
		ioOpsSpeed = 0
	}

	lastIoRead = totalReadBytes
	lastIoWrite = totalWriteBytes
	lastIoOps = totalOps
	lastIoTime = now

	return nil
}

func GetIoReadMo() (float64, error) {
	if err := updateIoStats(); err != nil {
		return 0, err
	}
	return ioReadSpeed, nil
}

func GetIoWriteMo() (float64, error) {
	if err := updateIoStats(); err != nil {
		return 0, err
	}
	return ioWriteSpeed, nil
}

func GetIoOps() (int, error) {
	if err := updateIoStats(); err != nil {
		return 0, err
	}
	return ioOpsSpeed, nil
}
