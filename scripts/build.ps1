$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path bin | Out-Null

go test ./...
go build -o bin/doctor.exe ./cmd/doctor

Write-Host "Built bin/doctor.exe"