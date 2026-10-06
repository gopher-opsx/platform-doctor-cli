package cli

import (
	"fmt"
	"os"
	"strings"

	reportpkg "github.com/gopher-opsx/platform-doctor-cli/internal/report"
	"github.com/spf13/cobra"
)

var (
	reportCompose bool
	reportFormat  string
	reportOutput  string
	reportLogs    int
	reportSince   string
)

var reportCmd = &cobra.Command{
	Use:   "report [service]",
	Short: "Generate a troubleshooting evidence report",

	Args: func(
		cmd *cobra.Command,
		args []string,
	) error {

		if reportCompose {
			if len(args) != 0 {
				return fmt.Errorf(
					"do not provide a service when using --compose",
				)
			}

			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf(
				"provide one service or use --compose",
			)
		}

		return nil
	},

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf(
				"get current directory: %w",
				err,
			)
		}

		target := ""

		if len(args) == 1 {
			target = args[0]
		}

		evidenceReport, err := reportpkg.Collect(
			currentDir,
			target,
			reportCompose,
			reportLogs,
			reportSince,
		)
		if err != nil {
			return err
		}

		rendered, err := reportpkg.Render(
			evidenceReport,
			reportFormat,
		)
		if err != nil {
			return err
		}

		if strings.TrimSpace(reportOutput) != "" {

			if err := os.WriteFile(
				reportOutput,
				[]byte(rendered),
				0644,
			); err != nil {
				return fmt.Errorf(
					"write report: %w",
					err,
				)
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"Report written to %s\n",
				reportOutput,
			)

			return nil
		}

		fmt.Fprint(
			cmd.OutOrStdout(),
			rendered,
		)

		return nil
	},
}

func init() {

	reportCmd.Flags().BoolVar(
		&reportCompose,
		"compose",
		false,
		"report all discovered Compose services",
	)

	reportCmd.Flags().StringVar(
		&reportFormat,
		"format",
		"text",
		"report format: text, json, or markdown",
	)

	reportCmd.Flags().StringVar(
		&reportOutput,
		"output",
		"",
		"write report to a file",
	)

	reportCmd.Flags().IntVar(
		&reportLogs,
		"logs",
		50,
		"number of recent log lines",
	)

	reportCmd.Flags().StringVar(
		&reportSince,
		"since",
		"15m",
		"log time window",
	)

	rootCmd.AddCommand(reportCmd)
}
