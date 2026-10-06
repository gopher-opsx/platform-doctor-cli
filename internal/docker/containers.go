package docker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

type Container struct {
	Name   string `json:"Names"`
	Image  string `json:"Image"`
	State  string `json:"State"`
	Status string `json:"Status"`
	Ports  string `json:"Ports"`
}

func ListContainers(dir string) ([]Container, error) {
	result, err := runner.Run(
		dir,
		"docker",
		"ps",
		"-a",
		"--format",
		"{{json .}}",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list Docker containers: %w",
			err,
		)
	}

	output := strings.TrimSpace(result.Stdout)

	if output == "" {
		return []Container{}, nil
	}

	lines := strings.Split(output, "\n")

	containers := make([]Container, 0, len(lines))

	for _, line := range lines {
		var container Container

		if err := json.Unmarshal(
			[]byte(line),
			&container,
		); err != nil {
			return nil, fmt.Errorf(
				"parse Docker container information: %w",
				err,
			)
		}

		containers = append(
			containers,
			container,
		)
	}

	return containers, nil
}
