package network

import (
	"pulse/internal/utils"
	"time"

	psnet "github.com/shirou/gopsutil/v4/net"
)

type interfaceState struct {
	bytesRecv    uint64
	bytesSent    uint64
	lastUpdate   time.Time
	lastSpeedIn  float64
	lastSpeedOut float64
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
	key := interName + "_in"
	state, exists := netHistory[key]

	if !exists {
		netHistory[key] = interfaceState{bytesRecv: currentRecv, lastUpdate: now}
		return 0, nil
	}

	duration := now.Sub(state.lastUpdate).Seconds()
	if duration < 1.0 {
		return state.lastSpeedIn, nil
	}

	deltaBytes := currentRecv - state.bytesRecv
	speed := utils.OctToMo(float64(deltaBytes) / duration)

	state.bytesRecv = currentRecv
	state.lastUpdate = now
	state.lastSpeedIn = speed
	netHistory[key] = state

	return speed, nil
}

func GetNetIoOut(interName string) (float64, error) {
	_, currentSent, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	key := interName + "_out"
	state, exists := netHistory[key]

	if !exists {
		netHistory[key] = interfaceState{bytesSent: currentSent, lastUpdate: now}
		return 0, nil
	}

	duration := now.Sub(state.lastUpdate).Seconds()
	if duration < 1.0 {
		return state.lastSpeedOut, nil
	}

	deltaBytes := currentSent - state.bytesSent
	speed := utils.OctToMo(float64(deltaBytes) / duration)

	state.bytesSent = currentSent
	state.lastUpdate = now
	state.lastSpeedOut = speed
	netHistory[key] = state

	return speed, nil
}

func GetNetIoTotal(interName string) (float64, error) {
	currentRecv, currentSent, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	key := interName + "_total"
	state, exists := netHistory[key]

	if !exists {
		netHistory[key] = interfaceState{bytesRecv: currentRecv, bytesSent: currentSent, lastUpdate: now}
		return 0, nil
	}

	duration := now.Sub(state.lastUpdate).Seconds()
	if duration < 1.0 {
		return state.lastSpeedIn, nil
	}

	deltaBytes := (currentRecv - state.bytesRecv) + (currentSent - state.bytesSent)
	speed := utils.OctToMo(float64(deltaBytes) / duration)

	state.bytesRecv = currentRecv
	state.bytesSent = currentSent
	state.lastUpdate = now
	state.lastSpeedIn = speed
	netHistory[key] = state

	return speed, nil
}

func GetNetIoSpeed(interName string) (float64, float64, error) {
	currentRecv, currentSent, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, 0, err
	}

	now := time.Now()
	key := interName + "_speed"
	state, exists := netHistory[key]

	if !exists {
		netHistory[key] = interfaceState{bytesRecv: currentRecv, bytesSent: currentSent, lastUpdate: now}
		return 0, 0, nil
	}

	duration := now.Sub(state.lastUpdate).Seconds()
	if duration < 1.0 {
		return state.lastSpeedIn, state.lastSpeedOut, nil
	}

	recvSpeed := utils.OctToMo(float64(currentRecv-state.bytesRecv) / duration)
	sentSpeed := utils.OctToMo(float64(currentSent-state.bytesSent) / duration)

	state.bytesRecv = currentRecv
	state.bytesSent = currentSent
	state.lastUpdate = now
	state.lastSpeedIn = recvSpeed
	state.lastSpeedOut = sentSpeed
	netHistory[key] = state

	return recvSpeed, sentSpeed, nil
}

func GetNetIoBytesIn(interName string) (float64, error) {
	recv, _, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}
	return utils.OctToGB(float64(recv)), nil
}

func GetNetIoBytesOut(interName string) (float64, error) {
	_, sent, err := getInterfaceCounters(interName)
	if err != nil {
		return 0, err
	}
	return utils.OctToGB(float64(sent)), nil
}
