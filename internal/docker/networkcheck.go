package docker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
)

type NetworkProbe struct {
	Status string
	Detail string
}

func CheckDNS(
	dir string,
	container string,
	host string,
) NetworkProbe {

	script := `
if command -v getent >/dev/null 2>&1; then
	getent hosts "$1"
	exit $?
fi

if command -v nslookup >/dev/null 2>&1; then
	nslookup "$1"
	exit $?
fi

exit 127
`

	result, err := runner.Run(
		dir,
		"docker",
		"exec",
		container,
		"sh",
		"-c",
		script,
		"doctor",
		host,
	)

	if err == nil {
		return NetworkProbe{
			Status: "resolved",
			Detail: firstOutputLine(
				result.Stdout,
			),
		}
	}

	errorText := strings.ToLower(
		err.Error(),
	)

	if strings.Contains(
		errorText,
		"exit status 127",
	) ||
		strings.Contains(
			errorText,
			"executable file not found",
		) {

		return NetworkProbe{
			Status: "unavailable",
			Detail: "DNS utility not available in container",
		}
	}

	return NetworkProbe{
		Status: "unresolved",
		Detail: host,
	}
}

func CheckTCP(
	dir string,
	container string,
	host string,
	port int,
) NetworkProbe {

	script := `
if command -v nc >/dev/null 2>&1; then
	nc -z -w 2 "$1" "$2"
	exit $?
fi

exit 127
`

	result, err := runner.Run(
		dir,
		"docker",
		"exec",
		container,
		"sh",
		"-c",
		script,
		"doctor",
		host,
		strconv.Itoa(port),
	)

	if err == nil {
		return NetworkProbe{
			Status: "reachable",
			Detail: fmt.Sprintf(
				"%s:%d",
				host,
				port,
			),
		}
	}

	_ = result

	errorText := strings.ToLower(
		err.Error(),
	)

	if strings.Contains(
		errorText,
		"exit status 127",
	) ||
		strings.Contains(
			errorText,
			"executable file not found",
		) {

		return NetworkProbe{
			Status: "unavailable",
			Detail: "TCP probe utility not available in container",
		}
	}

	return NetworkProbe{
		Status: "unreachable",
		Detail: fmt.Sprintf(
			"%s:%d",
			host,
			port,
		),
	}
}

func firstOutputLine(
	value string,
) string {

	value = strings.TrimSpace(value)

	if value == "" {
		return "-"
	}

	lines := strings.Split(
		value,
		"\n",
	)

	return strings.TrimSpace(
		lines[0],
	)
}
