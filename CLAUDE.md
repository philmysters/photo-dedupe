# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
# Run all tests with coverage
go test ./... -coverprofile=coverage.out

# Check coverage breakdown
go tool cover -func=coverage.out

# Build
go build ./cmd/main.go

# Run the CLI
go run ./cmd/main.go -in1 <folder1> -in2 <folder2> -out <output> [-config photo_dedupe.yaml] [--dryrun]

# Run a single test
go test ./internal/... -run TestFindPhotoFiles

# Lint (requires golangci-lint installed)
golangci-lint run

# Pre-commit hooks (run on every commit automatically after setup)
pre-commit install
pre-commit run --all-files
```

## Architecture

The tool takes two input photo directories and an output directory, deduplicates across them, and writes results to the output folder. The CLI is in `cmd/main.go`; all logic lives in `internal/`.

**`internal/config_loader.go`** — Loads `photo_dedupe.yaml` (or a custom path via `-config`). The only config field currently is `supported_extensions`, a list of lowercase extension strings (no dots). Defaults to `photo_dedupe.yaml` in the working directory.

**`internal/dedupe.go`** — `FindPhotoFiles(root, exts)` walks a directory tree and returns `[]PhotoFile` for files whose extensions (case-insensitive, dot-stripped) are in the supported set. Deduplication logic is not yet implemented; `cmd/main.go` currently just prints discovered files.

**`testdata/`** — Fixture files for tests. Integration tests should go in `/integration_tests` using `/testdata` as fixtures (not yet created).

## Quality requirements

CI enforces **98% test coverage** — every new function needs tests. The coverage check runs:
```sh
awk '/^total:/ {gsub(/%/,""); if ($3+0 < 98) exit 1}' <(go tool cover -func=coverage.out)
```

Pre-commit runs `gofmt`, `go vet`, and `golangci-lint` on every commit.
