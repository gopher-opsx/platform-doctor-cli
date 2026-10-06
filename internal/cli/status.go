package cli

import (
	"fmt"
	"os"

	dockerinfo "github.com/gopher-opsx/platform-doctor-cli/internal/docker"
	"github.com/spf13/cobra"
)

var statusCompose bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show container or Compose service status",

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		out := cmd.OutOrStdout()

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf(
				"get current directory: %w",
				err,
			)
		}

		if statusCompose {
			services, err :=
				dockerinfo.DiscoverComposeServices(
					currentDir,
				)
			if err != nil {
				return err
			}

			if len(services) == 0 {
				fmt.Fprintln(
					out,
					"No Docker Compose services found.",
				)

				return nil
			}

			fmt.Fprintf(
				out,
				"%-20s %-24s %-12s %-16s %s\n",
				"PROJECT",
				"SERVICE",
				"STATE",
				"HEALTH",
				"CONTAINER",
			)

			for _, service := range services {

				fmt.Fprintf(
					out,
					"%-20s %-24s %-12s %-16s %s\n",
					service.Project,
					service.Service,
					service.State,
					service.Health,
					service.Container,
				)
			}

			return nil
		}

		containers, err :=
			dockerinfo.ListContainers(
				currentDir,
			)
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			fmt.Fprintln(
				out,
				"No Docker containers found.",
			)

			return nil
		}

		fmt.Fprintf(
			out,
			"%-34s %-12s %-30s %s\n",
			"CONTAINER",
			"STATE",
			"STATUS",
			"IMAGE",
		)

		for _, container := range containers {

			fmt.Fprintf(
				out,
				"%-34s %-12s %-30s %s\n",
				container.Name,
				container.State,
				container.Status,
				container.Image,
			)
		}

		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVar(
		&statusCompose,
		"compose",
		false,
		"show Docker Compose services",
	)

	rootCmd.AddCommand(
		statusCmd,
	)
}
