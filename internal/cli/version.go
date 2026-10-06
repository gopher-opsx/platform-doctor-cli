package cli

import (
	"fmt"

	"github.com/gopher-opsx/platform-doctor-cli/internal/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Platform Doctor CLI version information",

	Run: func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()

		fmt.Fprintf(
			out,
			"doctor %s\ncommit: %s\nbuilt: %s\n",
			buildinfo.Version,
			buildinfo.Commit,
			buildinfo.Date,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
