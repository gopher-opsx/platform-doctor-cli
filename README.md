# Platform Doctor CLI

Platform Doctor CLI is a read-only troubleshooting evidence collector for Docker and Docker Compose environments.

It is designed to help engineers investigate containerized systems systematically:

> Understand the system → observe the symptom → collect evidence → form a hypothesis → test → identify the root cause → fix → verify.

Doctor does **not** automatically fix incidents or declare a root cause.

**It collects evidence. You diagnose the incident.**

---

## Why Platform Doctor?

Container troubleshooting often turns into a sequence of random commands:

```bash
docker ps
docker logs
docker inspect
docker stats
```

Each command is useful, but the evidence is scattered.

Platform Doctor brings that evidence together into a consistent troubleshooting workflow.

It can collect information about:

- Container state
- Docker health
- Restart count
- Restart policy
- Exit code
- OOM termination
- Image identity
- Runtime user
- CPU usage
- Memory usage
- PID count
- CPU and memory limits
- Ports
- Docker networks
- Mounts and volumes
- Environment configuration
- Secret-redacted configuration
- Recent bounded logs
- Application health
- Application readiness
- Docker Compose services
- Service dependencies
- DNS resolution
- TCP connectivity
- Troubleshooting reports

---

## Safety First

Platform Doctor is intentionally read-only.

Doctor does **not**:

- Restart containers
- Stop containers
- Remove containers
- Modify Docker Compose configuration
- Edit environment variables
- Repair dependencies
- Install packages
- Change application configuration
- Automatically declare a root cause

The purpose of Doctor is evidence collection.

The engineer remains responsible for interpreting that evidence.

---

## Commands

### Verify the Environment

```bash
doctor verify
```

Checks whether Docker and Docker Compose are available.

---

### Show Container Status

```bash
doctor status
```

Shows Docker container state.

For Docker Compose services:

```bash
doctor status --compose
```

---

### Inspect a Container

```bash
doctor inspect platform-lab-catalog-service
```

Collects detailed container evidence including:

- Container state
- Image
- Runtime user
- Docker health
- Restart information
- Exit information
- Resource usage
- Resource limits
- Ports
- Networks
- Mounts
- Safe configuration
- Application health/readiness
- Recent logs

---

### Diagnose a Service

```bash
doctor diagnose catalog-service
```

Doctor discovers the Docker Compose service and collects dependency evidence.

Example:

```text
DOCTOR SERVICE DIAGNOSIS

Project:    platform-lab
Service:    catalog-service
Container:  platform-lab-catalog-service
State:      running
Health:     healthy

DEPENDENCIES

PostgreSQL (postgres:5432)
  DNS: resolved
  TCP: reachable
```

Doctor reports evidence.

It does not automatically conclude whether PostgreSQL, networking, configuration, or another component is the root cause.

---

### Diagnose the Compose Environment

```bash
doctor diagnose --compose
```

Collects evidence across discovered Docker Compose services.

---

## Reports

### Text Report

```bash
doctor report catalog-service
```

### JSON Report

```bash
doctor report catalog-service --format json
```

### Markdown Report

```bash
doctor report catalog-service --format markdown
```

### Save a Report

```bash
doctor report catalog-service \
  --format markdown \
  --output catalog-evidence.md
```

### Generate a Platform-Wide Report

```bash
doctor report --compose \
  --format markdown \
  --output platform-evidence.md
```

### Control Log Collection

```bash
doctor report catalog-service \
  --logs 100 \
  --since 30m
```

---

## Platform Lab Integration

Platform Doctor is designed as a general Docker and Docker Compose troubleshooting utility.

It also includes a Platform Lab profile used by the Platform Lab training environment.

Known Platform Lab relationships include:

```text
catalog-service
    └── PostgreSQL

cart-service
    └── Redis

order-service
    ├── PostgreSQL
    └── Kafka

inventory-service
    ├── PostgreSQL
    └── Kafka

payment-service
    ├── PostgreSQL
    └── Kafka

notification-service
    ├── PostgreSQL
    └── Kafka
```

This allows Doctor to collect dependency evidence from the perspective of the affected application container.

Examples include:

- DNS resolution
- TCP reachability
- Container state
- Service health
- Application readiness

---

## Example Investigation

An incident might begin with:

```text
Catalog page is unavailable.
```

Start by understanding the overall platform state:

```bash
doctor status --compose
```

Then investigate the affected service:

```bash
doctor diagnose catalog-service
```

Collect deeper container evidence:

```bash
doctor inspect platform-lab-catalog-service
```

Generate an evidence artifact:

```bash
doctor report catalog-service \
  --format markdown \
  --output catalog-evidence.md
```

The intended troubleshooting process is:

```text
Symptom
   ↓
Understand request path
   ↓
Collect evidence
   ↓
Form hypothesis
   ↓
Test hypothesis
   ↓
Root cause
   ↓
Fix
   ↓
Verify
```

---

## Lab CLI and Doctor CLI

Platform Lab uses two separate tools with different responsibilities.

### Platform Lab CLI

The Lab CLI introduces controlled failures.

```text
Healthy Platform Lab
       ↓
Lab CLI
       ↓
Controlled failure
```

### Platform Doctor CLI

Doctor collects troubleshooting evidence from the environment.

```text
Controlled failure
       ↓
Doctor CLI
       ↓
Evidence
       ↓
Engineer investigation
       ↓
Root cause
       ↓
Recovery
```

Doctor does not inject failures.

Lab CLI does not diagnose them.

This separation keeps troubleshooting exercises predictable and realistic.

---

## Build from Source

### Requirements

- Go
- Docker
- Docker Compose

Clone the repository:

```bash
git clone https://github.com/gopher-opsx/platform-doctor-cli.git
cd platform-doctor-cli
```

Download dependencies:

```bash
go mod tidy
```

Run tests:

```bash
go test ./...
```

Build:

```bash
go build -o bin/doctor ./cmd/doctor
```

Check the binary:

```bash
bin/doctor version
```

On Windows:

```powershell
bin/doctor.exe version
```

---

## Build Scripts

### Linux / macOS

```bash
./scripts/build.sh
```

### Windows PowerShell

```powershell
./scripts/build.ps1
```

---

## Release Builds

Platform Doctor uses GoReleaser.

Create a local snapshot build:

```bash
goreleaser release --snapshot --clean
```

Release artifacts are generated for:

| Operating System | Architecture |
|---|---|
| Linux | amd64 |
| Linux | arm64 |
| Windows | amd64 |
| Windows | arm64 |
| macOS | amd64 |
| macOS | arm64 |

Checksums are also generated for release artifacts.

---

## Version

Check the installed version:

```bash
doctor version
```

Example:

```text
doctor v1.0.1
commit: abc1234
built: 2026-10-07T00:00:00Z
```

---

## Design Principles

Platform Doctor follows a few important principles.

### Read-Only

Evidence collection should not alter the system being investigated.

### Evidence Before Conclusions

Doctor presents what the system is reporting.

It does not automatically claim a root cause.

### Bounded Evidence

Logs are limited by line count and time window so reports remain focused.

### Safe Configuration Collection

Sensitive environment values are redacted before being displayed or written into reports.

### Reusable Core

Docker and Docker Compose evidence collection remains generic.

Platform-specific knowledge is kept in profiles.

---

## Troubleshooting Philosophy

Production troubleshooting should begin with the system and the symptom — not with a random command.

The investigation should move deliberately:

```text
What is the symptom?
        ↓
Which component is involved?
        ↓
What evidence do we need?
        ↓
What does the evidence show?
        ↓
What hypothesis can we form?
        ↓
How can we test it?
```

Platform Doctor exists to make that evidence collection:

- Repeatable
- Structured
- Safe
- Understandable
- Reusable

It intentionally separates **evidence collection** from **engineering judgment**.

---

## Security

Platform Doctor may inspect runtime configuration and container metadata.

Sensitive environment values are redacted before being presented.

Do not use Doctor as a replacement for proper secrets management.

For security issues, see:

```text
SECURITY.md
```

---

## Contributing

Contributions are welcome.

Before contributing, see:

```text
CONTRIBUTING.md
CODE_OF_CONDUCT.md
```

---

## Changelog

Release history is documented in:

```text
CHANGELOG.md
```

---

## License

Platform Doctor CLI is licensed under the Apache License 2.0.

See:

```text
LICENSE
```

for details.