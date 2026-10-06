package cli

import (
	"fmt"
	"os"

	"github.com/gopher-opsx/platform-doctor-cli/internal/runner"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify Doctor prerequisites",
	Long: `Verify checks that the local environment is ready for Doctor.

It checks Docker and Docker Compose availability.

This command is read-only and makes no changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		fmt.Fprintln(out, "Verifying Doctor environment...")
		fmt.Fprintln(out)

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}

		fmt.Fprintf(out, "✓ Working directory: %s\n", currentDir)

		if _, err := runner.Run(
			currentDir,
			"docker",
			"version",
		); err != nil {
			return fmt.Errorf(
				"Docker is not available or the Docker daemon is not running: %w",
				err,
			)
		}

		fmt.Fprintln(out, "✓ Docker available")

		if _, err := runner.Run(
			currentDir,
			"docker",
			"compose",
			"version",
		); err != nil {
			return fmt.Errorf(
				"Docker Compose is not available: %w",
				err,
			)
		}

		fmt.Fprintln(out, "✓ Docker Compose available")

		fmt.Fprintln(out)
		fmt.Fprintln(out, "Doctor environment verified.")
		fmt.Fprintln(out, "No changes were made.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
