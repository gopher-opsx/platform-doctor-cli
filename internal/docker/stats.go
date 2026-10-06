package docker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

type ResourceEvidence struct {
	CPUPercent    string
	MemoryUsage   string
	MemoryLimit   string
	MemoryPercent string
	PIDs          string
}

type statsResponse struct {
	CPUPercent    string `json:"CPUPerc"`
	MemoryUsage   string `json:"MemUsage"`
	MemoryPercent string `json:"MemPerc"`
	PIDs          string `json:"PIDs"`
}

func CollectResources(
	dir string,
	containerName string,
) (*ResourceEvidence, error) {

	result, err := runner.Run(
		dir,
		"docker",
		"stats",
		containerName,
		"--no-stream",
		"--format",
		"{{json .}}",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"collect container resources for %s: %w",
			containerName,
			err,
		)
	}

	var stats statsResponse

	if err := json.Unmarshal(
		[]byte(strings.TrimSpace(result.Stdout)),
		&stats,
	); err != nil {
		return nil, fmt.Errorf(
			"parse Docker stats output: %w",
			err,
		)
	}

	memoryUsage := strings.TrimSpace(stats.MemoryUsage)
	memoryLimit := ""

	parts := strings.SplitN(memoryUsage, "/", 2)

	if len(parts) == 2 {
		memoryUsage = strings.TrimSpace(parts[0])
		memoryLimit = strings.TrimSpace(parts[1])
	}

	return &ResourceEvidence{
		CPUPercent:    stats.CPUPercent,
		MemoryUsage:   memoryUsage,
		MemoryLimit:   memoryLimit,
		MemoryPercent: stats.MemoryPercent,
		PIDs:          stats.PIDs,
	}, nil
}
