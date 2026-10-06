package report

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Render(
	report *Report,
	format string,
) (string, error) {

	switch strings.ToLower(format) {

	case "text":
		return renderText(report), nil

	case "json":
		data, err := json.MarshalIndent(
			report,
			"",
			"  ",
		)
		if err != nil {
			return "", fmt.Errorf(
				"render json report: %w",
				err,
			)
		}

		return string(data), nil

	case "markdown", "md":
		return renderMarkdown(report), nil

	default:
		return "", fmt.Errorf(
			"unsupported format %q; use text, json, or markdown",
			format,
		)
	}
}

func renderText(
	report *Report,
) string {

	var b strings.Builder

	fmt.Fprintln(
		&b,
		"PLATFORM DOCTOR REPORT",
	)

	fmt.Fprintf(
		&b,
		"Generated: %s\n",
		report.GeneratedAt,
	)

	for _, service := range report.Services {

		fmt.Fprintln(&b)
		fmt.Fprintln(
			&b,
			"========================================",
		)

		fmt.Fprintf(
			&b,
			"Service:        %s\n",
			service.Service,
		)

		fmt.Fprintf(
			&b,
			"Container:      %s\n",
			service.Container,
		)

		fmt.Fprintf(
			&b,
			"State:          %s\n",
			service.State,
		)

		fmt.Fprintf(
			&b,
			"Health:         %s\n",
			service.Health,
		)

		fmt.Fprintf(
			&b,
			"Image:          %s\n",
			service.Image,
		)

		fmt.Fprintf(
			&b,
			"Runtime user:   %s\n",
			service.RuntimeUser,
		)

		fmt.Fprintf(
			&b,
			"Restart count:  %d\n",
			service.RestartCount,
		)

		fmt.Fprintf(
			&b,
			"Restart policy: %s\n",
			service.RestartPolicy,
		)

		fmt.Fprintf(
			&b,
			"Exit code:      %d\n",
			service.ExitCode,
		)

		fmt.Fprintf(
			&b,
			"OOM killed:     %t\n",
			service.OOMKilled,
		)

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "RESOURCES")

		fmt.Fprintf(
			&b,
			"CPU:            %s\n",
			service.Resources.CPUPercent,
		)

		fmt.Fprintf(
			&b,
			"Memory:         %s\n",
			service.Resources.MemoryUsage,
		)

		fmt.Fprintf(
			&b,
			"Memory %%:       %s\n",
			service.Resources.MemoryPercent,
		)

		fmt.Fprintf(
			&b,
			"PIDs:           %s\n",
			service.Resources.PIDs,
		)

		fmt.Fprintf(
			&b,
			"Memory limit:   %s\n",
			service.ConfiguredMemoryLimit,
		)

		fmt.Fprintf(
			&b,
			"CPU limit:      %s\n",
			service.ConfiguredCPULimit,
		)

		if service.ApplicationHealth != "" {
			fmt.Fprintln(&b)
			fmt.Fprintln(&b, "APPLICATION")

			fmt.Fprintf(
				&b,
				"Health:        %s\n",
				service.ApplicationHealth,
			)

			fmt.Fprintf(
				&b,
				"Readiness:     %s\n",
				service.ApplicationReadiness,
			)
		}

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "DEPENDENCIES")

		if len(service.Dependencies) == 0 {
			fmt.Fprintln(&b, "None")
		} else {
			for _, dependency := range service.Dependencies {

				fmt.Fprintf(
					&b,
					"%s (%s:%d)\n",
					dependency.Name,
					dependency.Host,
					dependency.Port,
				)

				fmt.Fprintf(
					&b,
					"  DNS: %s — %s\n",
					dependency.DNS.Status,
					dependency.DNS.Detail,
				)

				fmt.Fprintf(
					&b,
					"  TCP: %s — %s\n",
					dependency.TCP.Status,
					dependency.TCP.Detail,
				)
			}
		}

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "PORTS")

		printLines(&b, service.Ports)

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "NETWORKS")

		printLines(&b, service.Networks)

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "MOUNTS")

		printLines(&b, service.Mounts)

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "CONFIGURATION")

		printLines(&b, service.Config)

		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "RECENT LOGS")

		printLines(&b, service.Logs)
	}

	return b.String()
}

func renderMarkdown(
	report *Report,
) string {

	var b strings.Builder

	fmt.Fprintln(
		&b,
		"# Platform Doctor Report",
	)

	fmt.Fprintln(&b)

	fmt.Fprintf(
		&b,
		"Generated: `%s`\n",
		report.GeneratedAt,
	)

	for _, service := range report.Services {

		fmt.Fprintln(&b)

		fmt.Fprintf(
			&b,
			"## %s\n\n",
			service.Service,
		)

		fmt.Fprintf(
			&b,
			"- Container: `%s`\n",
			service.Container,
		)

		fmt.Fprintf(
			&b,
			"- Image: `%s`\n",
			service.Image,
		)

		fmt.Fprintf(
			&b,
			"- State: `%s`\n",
			service.State,
		)

		fmt.Fprintf(
			&b,
			"- Health: `%s`\n",
			service.Health,
		)

		fmt.Fprintf(
			&b,
			"- Runtime user: `%s`\n",
			service.RuntimeUser,
		)

		fmt.Fprintf(
			&b,
			"- Restart count: `%d`\n",
			service.RestartCount,
		)

		fmt.Fprintf(
			&b,
			"- Exit code: `%d`\n",
			service.ExitCode,
		)

		fmt.Fprintf(
			&b,
			"- OOM killed: `%t`\n",
			service.OOMKilled,
		)

		if service.ApplicationHealth != "" {
			fmt.Fprintln(&b)
			fmt.Fprintln(
				&b,
				"### Application",
			)

			fmt.Fprintln(&b)

			fmt.Fprintf(
				&b,
				"- Health: `%s`\n",
				service.ApplicationHealth,
			)

			fmt.Fprintf(
				&b,
				"- Readiness: `%s`\n",
				service.ApplicationReadiness,
			)
		}

		fmt.Fprintln(&b)
		fmt.Fprintln(
			&b,
			"### Dependencies",
		)

		fmt.Fprintln(&b)

		if len(service.Dependencies) == 0 {
			fmt.Fprintln(&b, "None.")
		} else {
			for _, dependency := range service.Dependencies {

				fmt.Fprintf(
					&b,
					"- **%s** `%s:%d` — DNS: `%s`, TCP: `%s`\n",
					dependency.Name,
					dependency.Host,
					dependency.Port,
					dependency.DNS.Status,
					dependency.TCP.Status,
				)
			}
		}

		renderMarkdownSection(
			&b,
			"Ports",
			service.Ports,
		)

		renderMarkdownSection(
			&b,
			"Networks",
			service.Networks,
		)

		renderMarkdownSection(
			&b,
			"Mounts",
			service.Mounts,
		)

		renderMarkdownSection(
			&b,
			"Configuration",
			service.Config,
		)

		renderMarkdownCodeSection(
			&b,
			"Recent Logs",
			service.Logs,
		)
	}

	return b.String()
}

func printLines(
	b *strings.Builder,
	values []string,
) {

	if len(values) == 0 {
		fmt.Fprintln(b, "None")
		return
	}

	for _, value := range values {
		fmt.Fprintln(b, value)
	}
}

func renderMarkdownSection(
	b *strings.Builder,
	title string,
	values []string,
) {

	fmt.Fprintln(b)

	fmt.Fprintf(
		b,
		"### %s\n\n",
		title,
	)

	if len(values) == 0 {
		fmt.Fprintln(b, "None.")
		return
	}

	for _, value := range values {
		fmt.Fprintf(
			b,
			"- `%s`\n",
			value,
		)
	}
}

func renderMarkdownCodeSection(
	b *strings.Builder,
	title string,
	values []string,
) {

	fmt.Fprintln(b)

	fmt.Fprintf(
		b,
		"### %s\n\n",
		title,
	)

	if len(values) == 0 {
		fmt.Fprintln(b, "No recent logs.")
		return
	}

	fmt.Fprintln(b, "```text")

	for _, value := range values {
		fmt.Fprintln(b, value)
	}

	fmt.Fprintln(b, "```")
}
