package docker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

type InspectEvidence struct {
	Name         string
	Image        string
	State        string
	Health       string
	RestartCount int
	ExitCode     int
	OOMKilled    bool
	Ports        []string
	Networks     []string
}

type inspectResponse struct {
	Name  string `json:"Name"`
	State struct {
		Status    string `json:"Status"`
		ExitCode  int    `json:"ExitCode"`
		OOMKilled bool   `json:"OOMKilled"`
		Health    *struct {
			Status string `json:"Status"`
		} `json:"Health"`
	} `json:"State"`

	RestartCount int `json:"RestartCount"`

	Config struct {
		Image string `json:"Image"`
	} `json:"Config"`

	NetworkSettings struct {
		Ports map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"Ports"`

		Networks map[string]struct {
			IPAddress string `json:"IPAddress"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
}

func InspectContainer(
	dir string,
	containerName string,
) (*InspectEvidence, error) {

	result, err := runner.Run(
		dir,
		"docker",
		"inspect",
		containerName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"inspect container %s: %w",
			containerName,
			err,
		)
	}

	var responses []inspectResponse

	if err := json.Unmarshal(
		[]byte(result.Stdout),
		&responses,
	); err != nil {
		return nil, fmt.Errorf(
			"parse container inspect output: %w",
			err,
		)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf(
			"container not found: %s",
			containerName,
		)
	}

	raw := responses[0]

	health := "not configured"

	if raw.State.Status != "running" {
		health = "not applicable"
	} else if raw.State.Health != nil {
		health = raw.State.Health.Status
	}

	var ports []string

	for containerPort, bindings := range raw.NetworkSettings.Ports {
		if len(bindings) == 0 {
			ports = append(
				ports,
				containerPort,
			)
			continue
		}

		for _, binding := range bindings {
			ports = append(
				ports,
				fmt.Sprintf(
					"%s -> %s:%s",
					containerPort,
					binding.HostIP,
					binding.HostPort,
				),
			)
		}
	}

	var networks []string

	for name, network := range raw.NetworkSettings.Networks {
		networks = append(
			networks,
			fmt.Sprintf(
				"%s (%s)",
				name,
				network.IPAddress,
			),
		)
	}

	return &InspectEvidence{
		Name:         strings.TrimPrefix(raw.Name, "/"),
		Image:        raw.Config.Image,
		State:        raw.State.Status,
		Health:       health,
		RestartCount: raw.RestartCount,
		ExitCode:     raw.State.ExitCode,
		OOMKilled:    raw.State.OOMKilled,
		Ports:        ports,
		Networks:     networks,
	}, nil
}
