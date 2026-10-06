package docker

import (
	"fmt"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

func CollectLogs(
	dir string,
	containerName string,
	lines int,
	since string,
) ([]string, error) {

	args := []string{
		"logs",
		"--tail",
		fmt.Sprintf("%d", lines),
	}

	if strings.TrimSpace(since) != "" {
		args = append(
			args,
			"--since",
			since,
		)
	}

	args = append(
		args,
		containerName,
	)

	result, err := runner.Run(
		dir,
		"docker",
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"collect container logs for %s: %w",
			containerName,
			err,
		)
	}

	output := strings.TrimSpace(
		result.Stdout + result.Stderr,
	)

	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}
