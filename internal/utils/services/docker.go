package services

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

func GetDockerContainerState(cli *client.Client, containerName string) (string, error) {
	ctx := context.Background()
	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return "[ NOT FOUND ]", nil
	}

	if inspect.State.Running {
		return "[ RUNNING ]", nil
	} else if inspect.State.Dead || (inspect.State.ExitCode != 0 && inspect.State.ExitCode != 137 && inspect.State.ExitCode != 143) {
		return "[ CRASHED ]", nil
	}
	return "[ EXITED ]", nil
}

func GetDockerContainerHealth(cli *client.Client, containerName string) (string, error) {
	ctx := context.Background()
	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return "[ NOT FOUND ]", nil
	}

	if inspect.State.Health != nil {
		status := strings.ToUpper(inspect.State.Health.Status)
		if status == "HEALTHY" {
			return "[ HEALTHY ]", nil
		}
		return "[ " + status + " ]", nil
	}
	return "[ NO HEALTHCHECK ]", nil
}

func GetDockerContainerRestartedTime(cli *client.Client, containerName string) (string, error) {
	ctx := context.Background()
	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return "[ NOT FOUND ]", nil
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

	totalRestarts := state.ManualRestarts + inspect.RestartCount
	if inspect.State.Restarting {
		return fmt.Sprintf("Restarts: %d [ LOOPING ]", totalRestarts), nil
	}
	return fmt.Sprintf("Restarts: %d", totalRestarts), nil
}

func GetDockerTotal(cli *client.Client) (int, int, error) {
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return 0, 0, err
	}

	statsUp, statsDown := 0, 0
	for _, c := range containers {
		if c.State == "running" || c.State == "restarting" {
			statsUp++
		} else {
			statsDown++
		}
	}

	return statsUp, statsDown, nil
}
