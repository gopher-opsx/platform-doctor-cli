package docker

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/redact"
	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

func CollectConfig(
	dir string,
	containerName string,
) ([]string, error) {

	result, err := runner.Run(
		dir,
		"docker",
		"inspect",
		"--format",
		"{{json .Config.Env}}",
		containerName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"collect container configuration for %s: %w",
			containerName,
			err,
		)
	}

	var env []string

	if err := json.Unmarshal(
		[]byte(strings.TrimSpace(result.Stdout)),
		&env,
	); err != nil {
		return nil, fmt.Errorf(
			"parse container environment: %w",
			err,
		)
	}

	config := make([]string, 0, len(env))

	for _, entry := range env {
		parts := strings.SplitN(entry, "=", 2)

		key := parts[0]
		value := ""

		if len(parts) == 2 {
			value = parts[1]
		}

		value = redact.Value(key, value)

		config = append(
			config,
			fmt.Sprintf("%s=%s", key, value),
		)
	}

	sort.Strings(config)

	return config, nil
}
