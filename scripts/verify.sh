#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "==> go build ./..."
./bin/go build ./...

echo "==> go vet ./..."
./bin/go vet ./...

echo "==> go tool mygo vet ."
./bin/go tool mygo vet .

echo "==> astryx check --dir . --vet"
./bin/go run ./cmd/astryx check --dir . --vet

echo "==> go test -count=1 ./..."
./bin/go test -count=1 ./...
