#!/usr/bin/env bash

set -euo pipefail

mkdir -p bin

go test ./...
go build -o bin/doctor ./cmd/doctor

echo "Built bin/doctor"