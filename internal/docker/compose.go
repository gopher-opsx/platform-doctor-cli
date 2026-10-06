package docker

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

type ComposeServiceEvidence struct {
	Project   string
	Service   string
	Container string
	Image     string
	State     string
	Health    string
}

type composeInspectResponse struct {
	Name string `json:"Name"`

	Config struct {
		Image  string            `json:"Image"`
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`

	State struct {
		Status string `json:"Status"`

		Health *struct {
			Status string `json:"Status"`
		} `json:"Health"`
	} `json:"State"`
}

func DiscoverComposeServices(
	dir string,
) ([]ComposeServiceEvidence, error) {

	result, err := runner.Run(
		dir,
		"docker",
		"ps",
		"-a",
		"--filter",
		"label=com.docker.compose.service",
		"--format",
		"{{.ID}}",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"discover compose containers: %w",
			err,
		)
	}

	rawIDs := strings.TrimSpace(result.Stdout)

	if rawIDs == "" {
		return []ComposeServiceEvidence{}, nil
	}

	ids := strings.Fields(rawIDs)

	args := []string{"inspect"}
	args = append(args, ids...)

	inspectResult, err := runner.Run(
		dir,
		"docker",
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"inspect compose containers: %w",
			err,
		)
	}

	var containers []composeInspectResponse

	if err := json.Unmarshal(
		[]byte(inspectResult.Stdout),
		&containers,
	); err != nil {
		return nil, fmt.Errorf(
			"parse compose container metadata: %w",
			err,
		)
	}

	services := make(
		[]ComposeServiceEvidence,
		0,
		len(containers),
	)

	for _, container := range containers {
		project := container.Config.Labels["com.docker.compose.project"]
		service := container.Config.Labels["com.docker.compose.service"]

		if service == "" {
			continue
		}

		health := "not configured"

		if container.State.Status != "running" {
			health = "not applicable"
		} else if container.State.Health != nil {
			health = container.State.Health.Status
		}

		services = append(
			services,
			ComposeServiceEvidence{
				Project: project,
				Service: service,
				Container: strings.TrimPrefix(
					container.Name,
					"/",
				),
				Image:  container.Config.Image,
				State:  container.State.Status,
				Health: health,
			},
		)
	}

	sort.Slice(
		services,
		func(i, j int) bool {
			if services[i].Project == services[j].Project {
				return services[i].Service < services[j].Service
			}

			return services[i].Project < services[j].Project
		},
	)

	return services, nil
}
