package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Evidence collection assistant for containerized platforms",
	Long: `Doctor is a read-only troubleshooting evidence collector.

Doctor collects evidence from containers and Docker Compose environments.

Doctor does not modify the environment.
Doctor does not restart services.
Doctor does not fix incidents.

It collects evidence.
You diagnose the incident.`,
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}
