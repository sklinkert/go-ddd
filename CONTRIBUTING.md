# Contributing to go-ddd

Thanks for your interest in improving this template!

## Getting started

1. Fork and clone the repository.
2. Install Go 1.26+, Docker (for integration tests), and sqlc:
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```
3. Run the checks locally before opening a PR:
   ```bash
   make test    # full suite incl. testcontainers (needs Docker)
   make lint
   ```

## What makes a good contribution

- **Small and focused.** One pattern, fix, or improvement per PR.
- **Tests included.** Every behavior change needs a test. Integration tests use testcontainers.
- **Consistent style.** Run `make fmt`. Note the deliberate house conventions: `Id` instead of `ID` (see `.golangci.yml` for rationale), sparse comments, snake_case JSON.
- **Schema changes** need an `up` *and* `down` migration plus `sqlc generate`.

## Proposing bigger changes

Open an issue first for new patterns (e.g. new bounded contexts, event bus integrations) so we can discuss whether it fits the template's teaching scope. The goal is a template that stays small enough to read in an afternoon.

## Questions

Open an issue — questions about applying DDD in Go are welcome and often turn into documentation improvements.
