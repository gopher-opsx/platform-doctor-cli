package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	dockerinfo "github.com/gopher-opsx/platform-doctor-cli/internal/docker"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Docker container status",
	Long: `Status collects the current state of Docker containers.

This command is read-only and makes no changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf(
				"get current directory: %w",
				err,
			)
		}

		containers, err := dockerinfo.ListContainers(currentDir)
		if err != nil {
			return err
		}

		fmt.Fprintln(out, "DOCTOR STATUS")
		fmt.Fprintln(out)

		if len(containers) == 0 {
			fmt.Fprintln(out, "No Docker containers found.")
			return nil
		}

		writer := tabwriter.NewWriter(
			out,
			0,
			4,
			2,
			' ',
			0,
		)

		fmt.Fprintln(
			writer,
			"CONTAINER\tSTATE\tSTATUS\tIMAGE",
		)

		for _, container := range containers {
			fmt.Fprintf(
				writer,
				"%s\t%s\t%s\t%s\n",
				container.Name,
				container.State,
				container.Status,
				container.Image,
			)
		}

		return writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
