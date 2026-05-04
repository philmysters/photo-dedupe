# photo-dedupe

A public Go utility by philmysters for deduplication and date-organization of photo files (JPEG, RAW, HEIC, and more; see `photo_dedupe.yaml`).

MIT Licensed.

This project enforces:
- Pre-commit hooks (formatting, vet, multi-linter)
- Automated CI builds for Go (lint, vet, and **test coverage ≥98%**) on every PR
- Dynamic code coverage badge via Codecov
- Pre-commit status badge

[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/philmysters/photo-dedupe/main.svg)](https://results.pre-commit.ci/latest/github/philmysters/photo-dedupe/main)

## Quickstart

```sh
# Clone and run pre-commit setup
brew install pre-commit     # Or: pip install pre-commit
pre-commit install

# Run the CLI (to be implemented)
go run ./cmd/main.go
```

See CONTRIBUTING.md for more details.

## Quality and Coverage

PRs must:
- Pass all pre-commit checks (lint, format, vet)
- Achieve and maintain 98% code coverage (enforced by CI)
- Pass all automated tests (unit and integration)

[![codecov](https://codecov.io/gh/philmysters/photo-dedupe/branch/main/graph/badge.svg)](https://codecov.io/gh/philmysters/photo-dedupe)
