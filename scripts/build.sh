#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p bin
go build -ldflags="-s -w" -o bin/share ./cmd/share
echo "Built bin/share"
