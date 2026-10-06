package docker

import "fmt"

func FormatMemoryLimit(bytes int64) string {
	if bytes <= 0 {
		return "unlimited"
	}

	const (
		kiB = 1024
		miB = 1024 * kiB
		giB = 1024 * miB
	)

	switch {
	case bytes >= giB:
		return fmt.Sprintf(
			"%.2f GiB",
			float64(bytes)/float64(giB),
		)

	case bytes >= miB:
		return fmt.Sprintf(
			"%.2f MiB",
			float64(bytes)/float64(miB),
		)

	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func FormatCPULimit(nanoCPUs int64) string {
	if nanoCPUs <= 0 {
		return "unlimited"
	}

	cpus := float64(nanoCPUs) / 1_000_000_000

	return fmt.Sprintf("%.2f CPUs", cpus)
}
