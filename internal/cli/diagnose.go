package cli

import (
	"fmt"
	"os"

	dockerinfo "github.com/gopher-opsx/platform-doctor-cli/internal/docker"
	"github.com/gopher-opsx/platform-doctor-cli/internal/profile"
	"github.com/spf13/cobra"
)

var diagnoseCompose bool

var diagnoseCmd = &cobra.Command{
	Use: "diagnose <service>",

	Short: "Collect troubleshooting evidence for a service",

	Args: func(
		cmd *cobra.Command,
		args []string,
	) error {

		if diagnoseCompose {
			if len(args) != 0 {
				return fmt.Errorf(
					"do not provide a service when using --compose",
				)
			}

			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf(
				"provide exactly one service name",
			)
		}

		return nil
	},

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

		services, err :=
			dockerinfo.DiscoverComposeServices(
				currentDir,
			)
		if err != nil {
			return err
		}

		if diagnoseCompose {
			if len(services) == 0 {
				fmt.Fprintln(
					out,
					"No Docker Compose services found.",
				)

				return nil
			}

			fmt.Fprintln(
				out,
				"DOCTOR COMPOSE DIAGNOSIS",
			)

			for _, service := range services {
				fmt.Fprintln(out)

				renderServiceDiagnosis(
					out,
					currentDir,
					service,
				)
			}

			return nil
		}

		requested := args[0]

		service, found :=
			findComposeService(
				services,
				requested,
			)

		if !found {
			return fmt.Errorf(
				"compose service not found: %s",
				requested,
			)
		}

		fmt.Fprintln(
			out,
			"DOCTOR SERVICE DIAGNOSIS",
		)
		fmt.Fprintln(out)

		renderServiceDiagnosis(
			out,
			currentDir,
			service,
		)

		return nil
	},
}

func findComposeService(
	services []dockerinfo.ComposeServiceEvidence,
	requested string,
) (dockerinfo.ComposeServiceEvidence, bool) {

	for _, service := range services {
		if service.Container == requested {
			return service, true
		}
	}

	for _, service := range services {
		if service.Service == requested {
			return service, true
		}
	}

	return dockerinfo.ComposeServiceEvidence{}, false
}

func renderServiceDiagnosis(
	out interface {
		Write([]byte) (int, error)
	},
	currentDir string,
	service dockerinfo.ComposeServiceEvidence,
) {

	fmt.Fprintf(
		out,
		"Project:    %s\n",
		service.Project,
	)

	fmt.Fprintf(
		out,
		"Service:    %s\n",
		service.Service,
	)

	fmt.Fprintf(
		out,
		"Container:  %s\n",
		service.Container,
	)

	fmt.Fprintf(
		out,
		"Image:      %s\n",
		service.Image,
	)

	fmt.Fprintf(
		out,
		"State:      %s\n",
		service.State,
	)

	fmt.Fprintf(
		out,
		"Health:     %s\n",
		service.Health,
	)

	platformService, known :=
		profile.FindService(
			service.Service,
		)

	if !known {
		fmt.Fprintln(out)
		fmt.Fprintln(
			out,
			"DEPENDENCIES",
		)
		fmt.Fprintln(
			out,
			"No dependency profile available.",
		)

		return
	}

	fmt.Fprintln(out)
	fmt.Fprintln(
		out,
		"DEPENDENCIES",
	)

	if len(
		platformService.Dependencies,
	) == 0 {

		fmt.Fprintln(
			out,
			"No configured dependencies.",
		)

		return
	}

	if service.State != "running" {
		fmt.Fprintln(
			out,
			"Dependency probes not applicable (container not running).",
		)

		return
	}

	for _, dependency := range platformService.Dependencies {

		fmt.Fprintf(
			out,
			"\n%s (%s:%d)\n",
			dependency.Name,
			dependency.Host,
			dependency.Port,
		)

		dns :=
			dockerinfo.CheckDNS(
				currentDir,
				service.Container,
				dependency.Host,
			)

		fmt.Fprintf(
			out,
			"  DNS:  %s",
			dns.Status,
		)

		if dns.Detail != "" {
			fmt.Fprintf(
				out,
				" — %s",
				dns.Detail,
			)
		}

		fmt.Fprintln(out)

		tcp :=
			dockerinfo.CheckTCP(
				currentDir,
				service.Container,
				dependency.Host,
				dependency.Port,
			)

		fmt.Fprintf(
			out,
			"  TCP:  %s",
			tcp.Status,
		)

		if tcp.Detail != "" {
			fmt.Fprintf(
				out,
				" — %s",
				tcp.Detail,
			)
		}

		fmt.Fprintln(out)
	}
}

func init() {
	diagnoseCmd.Flags().BoolVar(
		&diagnoseCompose,
		"compose",
		false,
		"collect evidence for discovered Compose services",
	)

	rootCmd.AddCommand(
		diagnoseCmd,
	)
}
