package services

import (
	"context"
	"sync"

	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v4/docker"
)

var (
	containerStats      = make(map[string]*containerState)
	containerStatsMutex sync.Mutex
)

type containerState struct {
	ManualRestarts   int
	LastStartedAt    string
	LastAutoRestarts int
}

func GetDockerContainerState(container string) (string, error) {
	do, err := docker.GetDockerStat()
	if err != nil {
		return "", err
	}

	var cn *docker.CgroupDockerStat
	for _, c := range do {
		if c.Name == container {
			cn = &c
			break
		}
	}

	if cn == nil {
		return "Container not found", nil
	}

	state := cn.Running
	if state {
		return "Running", nil
	} else {
		return "Exited", nil
	}
}

func GetDockerContainerUptime(container string) (string, error) {
	do, err := docker.GetDockerStat()
	if err != nil {
		return "", err
	}

	var cn *docker.CgroupDockerStat
	for _, c := range do {
		if c.Name == container {
			cn = &c
			break
		}
	}

	if cn == nil {
		return "Container not found", nil
	}

	return cn.Status, nil
}

func GetDockerContainerRestartedTime(containerName string) (int, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return 0, err
	}
	defer cli.Close()

	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return 0, err
	}

	containerStatsMutex.Lock()
	defer containerStatsMutex.Unlock()

	state, exists := containerStats[containerName]
	if !exists {
		state = &containerState{
			ManualRestarts:   0,
			LastStartedAt:    inspect.State.StartedAt,
			LastAutoRestarts: inspect.RestartCount,
		}
		containerStats[containerName] = state
	} else {
		if inspect.State.StartedAt != state.LastStartedAt {
			if inspect.RestartCount <= state.LastAutoRestarts {
				state.ManualRestarts++
			}
			state.LastStartedAt = inspect.State.StartedAt
		}
		state.LastAutoRestarts = inspect.RestartCount
	}

	return state.ManualRestarts + inspect.RestartCount, nil
}

func GetDockerContainerAutoRestartedTime(containerName string) (int, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return 0, err
	}
	defer cli.Close()

	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return 0, err
	}

	return inspect.RestartCount, nil
}

func GetDockerTotal() (int, int, error) {
	do, err := docker.GetDockerStat()
	if err != nil {
		return 0, 0, err
	}

	statsUp, statsDown := 0, 0
	for _, cn := range do {
		if cn.Running {
			statsUp++
		} else {
			statsDown++
		}
	}

	return statsUp, statsDown, nil
}
