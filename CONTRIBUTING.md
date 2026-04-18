# Contributing to photo-dedupe

Thank you for your interest in contributing!

## Workflow
- Fork the repo and create a feature branch for your changes.
- Run `pre-commit install` after cloning; all code must pass formatting (`gofmt`), vet, and `golangci-lint`.
- All PRs must pass automated CI checks (lint, tests, and 98%+ coverage) before merge.
- Open a Pull Request (PR) against the `main` branch; describe your change, reference issues, and include tests if possible.

## Commit Messages
- Use clear, descriptive commit messages.
- Reference issues when appropriate: e.g. `fix: handle EXIF for DNG (closes #42)`

## Code Style
- Go code must be formatted (`gofmt`).
- No style warnings or errors from `golangci-lint`.

## Tests and Coverage
- All code must be tested.
- PRs must result in **at least 98% overall test coverage** (enforced by CI, see README).
- All test suites must be non-empty and add meaningful coverage.
- Run locally:

  ```sh
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out
  ```

## Integration tests
- End-to-end scenarios should be in `/integration_tests` and use `/testdata` as fixtures.
