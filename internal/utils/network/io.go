package network

import (
	"pulse/internal/utils"
	"time"

	psnet "github.com/shirou/gopsutil/v4/net"
)

type interfaceState struct {
	bytesRecv  uint64
	bytesSent  uint64
	lastUpdate time.Time
}

var netHistory = make(map[string]interfaceState)

func getInterfaceCounters(interName string) (uint64, uint64, error) {
	counters, err := psnet.IOCounters(true)
	if err != nil {
		return 0, 0, err
	}

	if interName == "" || interName == "all" {
		globalCounters, errGlobal := psnet.IOCounters(false)
		if errGlobal != nil || len(globalCounters) == 0 {
			return 0, 0, errGlobal
		}
		return globalCounters[0].BytesRecv, globalCounters[0].BytesSent, nil
	}

	for _, c := range counters {
		if c.Name == interName {
			return c.BytesRecv, c.BytesSent, nil
		}
	}

	return 0, 0, nil
}

func GetNetIoIn(interName string) (float64, error) {
	currentRecv, _, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	state, exists := netHistory[interName]

	if !exists {
		netHistory[interName] = interfaceState{bytesRecv: currentRecv, lastUpdate: now}
		return 0, nil
	}

	duration := now.Sub(state.lastUpdate).Seconds()
	if duration <= 0 {
		return 0, nil
	}

	deltaBytes := currentRecv - state.bytesRecv
	speed := float64(deltaBytes) / duration

	state.bytesRecv = currentRecv
	state.lastUpdate = now
	netHistory[interName] = state

	return utils.OctToMo(speed), nil
}

func GetNetIoOut(interName string) (float64, error) {
	_, currentSent, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	state, exists := netHistory[interName]

	if !exists {
		netHistory[interName] = interfaceState{bytesSent: currentSent, lastUpdate: now}
		return 0, nil
	}

	duration := now.Sub(state.lastUpdate).Seconds()
	if duration <= 0 {
		return 0, nil
	}

	deltaBytes := currentSent - state.bytesSent
	speed := float64(deltaBytes) / duration

	state.bytesSent = currentSent
	state.lastUpdate = now
	netHistory[interName] = state

	return utils.OctToMo(speed), nil
}

func GetNetIoTotal(interName string) (float64, error) {
	inSpeed, err := GetNetIoIn(interName)
	if err != nil {
		return 0, err
	}
	outSpeed, err := GetNetIoOut(interName)
	if err != nil {
		return 0, err
	}
	return inSpeed + outSpeed, nil
}

func GetNetIoSpeed(interName string) (float64, float64, error) {
	currentRecv, currentSent, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, 0, err
	}

	now := time.Now()
	state, exists := netHistory[interName]

	if !exists {
		netHistory[interName] = interfaceState{bytesRecv: currentRecv, bytesSent: currentSent, lastUpdate: now}
		return 0, 0, nil
	}

	duration := now.Sub(state.lastUpdate).Seconds()
	if duration <= 0 {
		return 0, 0, nil
	}

	recvSpeed := float64(currentRecv-state.bytesRecv) / duration
	sentSpeed := float64(currentSent-state.bytesSent) / duration

	state.bytesRecv = currentRecv
	state.bytesSent = currentSent
	state.lastUpdate = now
	netHistory[interName] = state

	return utils.OctToMo(recvSpeed), utils.OctToMo(sentSpeed), nil
}

func GetNetIoBytesIn(interName string) (float64, error) {
	recv, _, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}
	return utils.OctToGo(float64(recv)), nil
}

func GetNetIoBytesOut(interName string) (float64, error) {
	_, sent, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}
	return utils.OctToGo(float64(sent)), nil
}
