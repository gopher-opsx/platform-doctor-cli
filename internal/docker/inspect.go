package docker

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

type InspectEvidence struct {
	Name          string
	Image         string
	RuntimeUser   string
	State         string
	Health        string
	RestartPolicy string
	RestartCount  int
	ExitCode      int
	OOMKilled     bool
	StartedAt     string
	FinishedAt    string

	MemoryLimit int64
	NanoCPUs    int64

	Ports    []string
	Networks []string
	Mounts   []string
}

type inspectResponse struct {
	Name string `json:"Name"`

	State struct {
		Status     string `json:"Status"`
		ExitCode   int    `json:"ExitCode"`
		OOMKilled  bool   `json:"OOMKilled"`
		StartedAt  string `json:"StartedAt"`
		FinishedAt string `json:"FinishedAt"`

		Health *struct {
			Status string `json:"Status"`
		} `json:"Health"`
	} `json:"State"`

	RestartCount int `json:"RestartCount"`

	Config struct {
		Image string `json:"Image"`
		User  string `json:"User"`
	} `json:"Config"`

	HostConfig struct {
		Memory   int64 `json:"Memory"`
		NanoCPUs int64 `json:"NanoCpus"`

		RestartPolicy struct {
			Name string `json:"Name"`
		} `json:"RestartPolicy"`
	} `json:"HostConfig"`

	NetworkSettings struct {
		Ports map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"Ports"`

		Networks map[string]struct {
			IPAddress string `json:"IPAddress"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`

	Mounts []struct {
		Type        string `json:"Type"`
		Name        string `json:"Name"`
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		RW          bool   `json:"RW"`
	} `json:"Mounts"`
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

	runtimeUser := strings.TrimSpace(raw.Config.User)

	if runtimeUser == "" {
		runtimeUser = "default"
	}

	restartPolicy := strings.TrimSpace(
		raw.HostConfig.RestartPolicy.Name,
	)

	if restartPolicy == "" {
		restartPolicy = "no"
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

	sort.Strings(ports)

	var networks []string

	for name, network := range raw.NetworkSettings.Networks {
		ip := network.IPAddress

		if ip == "" {
			ip = "-"
		}

		networks = append(
			networks,
			fmt.Sprintf(
				"%s (%s)",
				name,
				ip,
			),
		)
	}

	sort.Strings(networks)

	var mounts []string

	for _, mount := range raw.Mounts {
		source := mount.Source

		if mount.Type == "volume" && mount.Name != "" {
			source = mount.Name
		}

		mode := "ro"

		if mount.RW {
			mode = "rw"
		}

		mounts = append(
			mounts,
			fmt.Sprintf(
				"%s: %s -> %s (%s)",
				mount.Type,
				source,
				mount.Destination,
				mode,
			),
		)
	}

	sort.Strings(mounts)

	return &InspectEvidence{
		Name:          strings.TrimPrefix(raw.Name, "/"),
		Image:         raw.Config.Image,
		RuntimeUser:   runtimeUser,
		State:         raw.State.Status,
		Health:        health,
		RestartPolicy: restartPolicy,
		RestartCount:  raw.RestartCount,
		ExitCode:      raw.State.ExitCode,
		OOMKilled:     raw.State.OOMKilled,
		StartedAt:     raw.State.StartedAt,
		FinishedAt:    raw.State.FinishedAt,
		MemoryLimit:   raw.HostConfig.Memory,
		NanoCPUs:      raw.HostConfig.NanoCPUs,
		Ports:         ports,
		Networks:      networks,
		Mounts:        mounts,
	}, nil
}
