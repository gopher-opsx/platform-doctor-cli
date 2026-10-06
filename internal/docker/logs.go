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
) ([]string, error) {

	result, err := runner.Run(
		dir,
		"docker",
		"logs",
		"--tail",
		fmt.Sprintf("%d", lines),
		containerName,
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
