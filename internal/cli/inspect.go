package cli

import (
	"fmt"
	"os"
	"strings"

	dockerinfo "github.com/gopher-opsx/platform-doctor-cli/internal/docker"
	"github.com/gopher-opsx/platform-doctor-cli/internal/httpcheck"
	"github.com/gopher-opsx/platform-doctor-cli/internal/profile"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <container>",
	Short: "Collect evidence from a Docker container",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf(
				"get current directory: %w",
				err,
			)
		}

		containerName := args[0]

		evidence, err := dockerinfo.InspectContainer(
			currentDir,
			containerName,
		)
		if err != nil {
			return err
		}

		fmt.Fprintln(out, "DOCTOR INSPECTION")
		fmt.Fprintln(out)

		// ------------------------------------------------------------
		// CONTAINER
		// ------------------------------------------------------------

		fmt.Fprintln(out, "CONTAINER")

		fmt.Fprintf(out, "Name:            %s\n", evidence.Name)
		fmt.Fprintf(out, "Image:           %s\n", evidence.Image)
		fmt.Fprintf(out, "Runtime user:    %s\n", evidence.RuntimeUser)
		fmt.Fprintf(out, "State:           %s\n", evidence.State)
		fmt.Fprintf(out, "Health:          %s\n", evidence.Health)
		fmt.Fprintf(out, "Restart count:   %d\n", evidence.RestartCount)
		fmt.Fprintf(out, "Restart policy:  %s\n", evidence.RestartPolicy)
		fmt.Fprintf(out, "Exit code:       %d\n", evidence.ExitCode)
		fmt.Fprintf(out, "OOM killed:      %t\n", evidence.OOMKilled)
		fmt.Fprintf(out, "Started at:      %s\n", evidence.StartedAt)
		fmt.Fprintf(out, "Finished at:     %s\n", evidence.FinishedAt)

		// ------------------------------------------------------------
		// CURRENT RESOURCE USAGE
		// ------------------------------------------------------------

		resources, err := dockerinfo.CollectResources(
			currentDir,
			containerName,
		)
		if err != nil {
			return err
		}

		fmt.Fprintln(out)
		fmt.Fprintln(out, "RESOURCES")

		fmt.Fprintf(
			out,
			"CPU:             %s\n",
			resources.CPUPercent,
		)

		fmt.Fprintf(
			out,
			"Memory usage:    %s\n",
			resources.MemoryUsage,
		)

		fmt.Fprintf(
			out,
			"Memory limit:    %s\n",
			resources.MemoryLimit,
		)

		fmt.Fprintf(
			out,
			"Memory percent:  %s\n",
			resources.MemoryPercent,
		)

		fmt.Fprintf(
			out,
			"PIDs:            %s\n",
			resources.PIDs,
		)

		// ------------------------------------------------------------
		// CONFIGURED RESOURCE LIMITS
		// ------------------------------------------------------------

		fmt.Fprintln(out)
		fmt.Fprintln(out, "RESOURCE LIMITS")

		fmt.Fprintf(
			out,
			"Memory limit:    %s\n",
			dockerinfo.FormatMemoryLimit(
				evidence.MemoryLimit,
			),
		)

		fmt.Fprintf(
			out,
			"CPU limit:       %s\n",
			dockerinfo.FormatCPULimit(
				evidence.NanoCPUs,
			),
		)

		// ------------------------------------------------------------
		// PORTS
		// ------------------------------------------------------------

		fmt.Fprintln(out)
		fmt.Fprintln(out, "PORTS")

		if len(evidence.Ports) == 0 {
			fmt.Fprintln(out, "None")
		} else {
			fmt.Fprintln(
				out,
				strings.Join(evidence.Ports, "\n"),
			)
		}

		// ------------------------------------------------------------
		// NETWORKS
		// ------------------------------------------------------------

		fmt.Fprintln(out)
		fmt.Fprintln(out, "NETWORKS")

		if len(evidence.Networks) == 0 {
			fmt.Fprintln(out, "None")
		} else {
			fmt.Fprintln(
				out,
				strings.Join(evidence.Networks, "\n"),
			)
		}

		// ------------------------------------------------------------
		// MOUNTS
		// ------------------------------------------------------------

		fmt.Fprintln(out)
		fmt.Fprintln(out, "MOUNTS")

		if len(evidence.Mounts) == 0 {
			fmt.Fprintln(out, "None")
		} else {
			for _, mount := range evidence.Mounts {
				fmt.Fprintln(out, mount)
			}
		}

		// ------------------------------------------------------------
		// CONFIGURATION
		// ------------------------------------------------------------

		config, err := dockerinfo.CollectConfig(
			currentDir,
			containerName,
		)
		if err != nil {
			return err
		}

		fmt.Fprintln(out)
		fmt.Fprintln(out, "CONFIGURATION")

		if len(config) == 0 {
			fmt.Fprintln(
				out,
				"No environment configuration found.",
			)
		} else {
			for _, entry := range config {
				fmt.Fprintln(out, entry)
			}
		}

		// ------------------------------------------------------------
		// APPLICATION HEALTH / READINESS
		// ------------------------------------------------------------

		if service, ok := findPlatformLabService(
			evidence.Name,
		); ok {

			fmt.Fprintln(out)
			fmt.Fprintln(out, "APPLICATION")

			if evidence.State != "running" {
				fmt.Fprintln(
					out,
					"Health:          not applicable (container not running)",
				)

				fmt.Fprintln(
					out,
					"Readiness:       not applicable (container not running)",
				)
			} else {
				healthURL := fmt.Sprintf(
					"http://localhost:%d%s",
					service.Port,
					service.HealthPath,
				)

				readyURL := fmt.Sprintf(
					"http://localhost:%d%s",
					service.Port,
					service.ReadyPath,
				)

				healthResult := httpcheck.Probe(
					healthURL,
				)

				readyResult := httpcheck.Probe(
					readyURL,
				)

				fmt.Fprintf(
					out,
					"Health:          %s\n",
					httpcheck.Format(
						healthResult,
					),
				)

				fmt.Fprintf(
					out,
					"Readiness:       %s\n",
					httpcheck.Format(
						readyResult,
					),
				)
			}
		}

		// ------------------------------------------------------------
		// RECENT LOGS
		// ------------------------------------------------------------

		logs, err := dockerinfo.CollectLogs(
			currentDir,
			containerName,
			50,
			"15m",
		)
		if err != nil {
			return err
		}

		fmt.Fprintln(out)
		fmt.Fprintln(
			out,
			"RECENT LOGS — LAST 15 MINUTES",
		)

		if len(logs) == 0 {
			fmt.Fprintln(
				out,
				"No recent logs.",
			)
		} else {
			for _, line := range logs {
				fmt.Fprintln(out, line)
			}
		}

		return nil
	},
}

func findPlatformLabService(
	containerName string,
) (profile.Service, bool) {

	for _, service := range profile.PlatformLabServices {
		if service.Container == containerName {
			return service, true
		}
	}

	return profile.Service{}, false
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
