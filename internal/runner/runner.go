package runner

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Result struct {
	Stdout string
	Stderr string
}

func Run(
	dir string,
	name string,
	args ...string,
) (Result, error) {

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		return result, fmt.Errorf(
			"run %s: %w: %s",
			name,
			err,
			stderr.String(),
		)
	}

	return result, nil
}
