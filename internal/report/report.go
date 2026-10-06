package report

import (
	"fmt"
	"time"

	dockerinfo "github.com/gopher-opsx/platform-doctor-cli/internal/docker"
	"github.com/gopher-opsx/platform-doctor-cli/internal/httpcheck"
	"github.com/gopher-opsx/platform-doctor-cli/internal/profile"
)

type DependencyEvidence struct {
	Name string                  `json:"name"`
	Host string                  `json:"host"`
	Port int                     `json:"port"`
	DNS  dockerinfo.NetworkProbe `json:"dns"`
	TCP  dockerinfo.NetworkProbe `json:"tcp"`
}

type ResourceEvidence struct {
	CPUPercent    string `json:"cpu_percent"`
	MemoryUsage   string `json:"memory_usage"`
	MemoryLimit   string `json:"memory_limit"`
	MemoryPercent string `json:"memory_percent"`
	PIDs          string `json:"pids"`
}

type ServiceReport struct {
	Project   string `json:"project"`
	Service   string `json:"service"`
	Container string `json:"container"`
	Image     string `json:"image"`

	RuntimeUser   string `json:"runtime_user"`
	State         string `json:"state"`
	Health        string `json:"health"`
	RestartPolicy string `json:"restart_policy"`
	RestartCount  int    `json:"restart_count"`
	ExitCode      int    `json:"exit_code"`
	OOMKilled     bool   `json:"oom_killed"`
	StartedAt     string `json:"started_at"`
	FinishedAt    string `json:"finished_at"`

	ConfiguredMemoryLimit string `json:"configured_memory_limit"`
	ConfiguredCPULimit    string `json:"configured_cpu_limit"`

	Resources ResourceEvidence `json:"resources"`

	Ports    []string `json:"ports"`
	Networks []string `json:"networks"`
	Mounts   []string `json:"mounts"`
	Config   []string `json:"config"`

	ApplicationHealth    string `json:"application_health,omitempty"`
	ApplicationReadiness string `json:"application_readiness,omitempty"`

	Dependencies []DependencyEvidence `json:"dependencies"`
	Logs         []string             `json:"logs"`
}

type Report struct {
	GeneratedAt string          `json:"generated_at"`
	Services    []ServiceReport `json:"services"`
}

func Collect(
	dir string,
	target string,
	composeMode bool,
	logLines int,
	since string,
) (*Report, error) {

	services, err := dockerinfo.DiscoverComposeServices(dir)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf(
			"no Docker Compose services found",
		)
	}

	selected := services

	if !composeMode {
		selected = nil

		for _, service := range services {
			if service.Service == target ||
				service.Container == target {

				selected = append(
					selected,
					service,
				)

				break
			}
		}

		if len(selected) == 0 {
			return nil, fmt.Errorf(
				"compose service not found: %s",
				target,
			)
		}
	}

	result := &Report{
		GeneratedAt: time.Now().
			UTC().
			Format(time.RFC3339),
	}

	for _, composeService := range selected {

		serviceReport, err := collectService(
			dir,
			composeService,
			logLines,
			since,
		)
		if err != nil {
			return nil, err
		}

		result.Services = append(
			result.Services,
			serviceReport,
		)
	}

	return result, nil
}

func collectService(
	dir string,
	composeService dockerinfo.ComposeServiceEvidence,
	logLines int,
	since string,
) (ServiceReport, error) {

	inspection, err := dockerinfo.InspectContainer(
		dir,
		composeService.Container,
	)
	if err != nil {
		return ServiceReport{}, err
	}

	resources, err := dockerinfo.CollectResources(
		dir,
		composeService.Container,
	)
	if err != nil {
		return ServiceReport{}, err
	}

	config, err := dockerinfo.CollectConfig(
		dir,
		composeService.Container,
	)
	if err != nil {
		return ServiceReport{}, err
	}

	logs, err := dockerinfo.CollectLogs(
		dir,
		composeService.Container,
		logLines,
		since,
	)
	if err != nil {
		return ServiceReport{}, err
	}

	serviceReport := ServiceReport{
		Project:   composeService.Project,
		Service:   composeService.Service,
		Container: composeService.Container,
		Image:     inspection.Image,

		RuntimeUser:   inspection.RuntimeUser,
		State:         inspection.State,
		Health:        inspection.Health,
		RestartPolicy: inspection.RestartPolicy,
		RestartCount:  inspection.RestartCount,
		ExitCode:      inspection.ExitCode,
		OOMKilled:     inspection.OOMKilled,
		StartedAt:     inspection.StartedAt,
		FinishedAt:    inspection.FinishedAt,

		ConfiguredMemoryLimit: dockerinfo.FormatMemoryLimit(
			inspection.MemoryLimit,
		),

		ConfiguredCPULimit: dockerinfo.FormatCPULimit(
			inspection.NanoCPUs,
		),

		Resources: ResourceEvidence{
			CPUPercent:    resources.CPUPercent,
			MemoryUsage:   resources.MemoryUsage,
			MemoryLimit:   resources.MemoryLimit,
			MemoryPercent: resources.MemoryPercent,
			PIDs:          resources.PIDs,
		},

		Ports:    inspection.Ports,
		Networks: inspection.Networks,
		Mounts:   inspection.Mounts,
		Config:   config,
		Logs:     logs,
	}

	platformService, known := profile.FindService(
		composeService.Service,
	)

	if !known {
		return serviceReport, nil
	}

	if composeService.State != "running" {
		serviceReport.ApplicationHealth =
			"not applicable (container not running)"

		serviceReport.ApplicationReadiness =
			"not applicable (container not running)"
	} else if platformService.Port > 0 {

		healthURL := fmt.Sprintf(
			"http://localhost:%d%s",
			platformService.Port,
			platformService.HealthPath,
		)

		readyURL := fmt.Sprintf(
			"http://localhost:%d%s",
			platformService.Port,
			platformService.ReadyPath,
		)

		serviceReport.ApplicationHealth =
			httpcheck.Format(
				httpcheck.Probe(healthURL),
			)

		serviceReport.ApplicationReadiness =
			httpcheck.Format(
				httpcheck.Probe(readyURL),
			)
	}

	for _, dependency := range platformService.Dependencies {

		item := DependencyEvidence{
			Name: dependency.Name,
			Host: dependency.Host,
			Port: dependency.Port,
		}

		if composeService.State != "running" {
			item.DNS = dockerinfo.NetworkProbe{
				Status: "not applicable",
				Detail: "container not running",
			}

			item.TCP = dockerinfo.NetworkProbe{
				Status: "not applicable",
				Detail: "container not running",
			}
		} else {
			item.DNS = dockerinfo.CheckDNS(
				dir,
				composeService.Container,
				dependency.Host,
			)

			item.TCP = dockerinfo.CheckTCP(
				dir,
				composeService.Container,
				dependency.Host,
				dependency.Port,
			)
		}

		serviceReport.Dependencies = append(
			serviceReport.Dependencies,
			item,
		)
	}

	return serviceReport, nil
}
