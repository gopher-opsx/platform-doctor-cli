package cli

import (
	"fmt"
	"os"
	"strings"

	dockerinfo "github.com/gopher-opsx/platform-doctor-cli/internal/docker"
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

		fmt.Fprintf(out, "Container:      %s\n", evidence.Name)
		fmt.Fprintf(out, "Image:          %s\n", evidence.Image)
		fmt.Fprintf(out, "State:          %s\n", evidence.State)
		fmt.Fprintf(out, "Health:         %s\n", evidence.Health)
		fmt.Fprintf(out, "Restart count:  %d\n", evidence.RestartCount)
		fmt.Fprintf(out, "Exit code:      %d\n", evidence.ExitCode)
		fmt.Fprintf(out, "OOM killed:     %t\n", evidence.OOMKilled)

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
			"CPU:            %s\n",
			resources.CPUPercent,
		)

		fmt.Fprintf(
			out,
			"Memory usage:   %s\n",
			resources.MemoryUsage,
		)

		fmt.Fprintf(
			out,
			"Memory limit:   %s\n",
			resources.MemoryLimit,
		)

		fmt.Fprintf(
			out,
			"Memory percent: %s\n",
			resources.MemoryPercent,
		)

		fmt.Fprintf(
			out,
			"PIDs:           %s\n",
			resources.PIDs,
		)

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
			fmt.Fprintln(out, "No environment configuration found.")
		} else {
			for _, entry := range config {
				fmt.Fprintln(out, entry)
			}
		}

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
		fmt.Fprintln(out, "RECENT LOGS — LAST 15 MINUTES")

		if len(logs) == 0 {
			fmt.Fprintln(out, "No recent logs.")
		} else {
			for _, line := range logs {
				fmt.Fprintln(out, line)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
