# Contributing

## Overview

SoloDev welcomes contributors who care about AI-native DevOps, MCP tooling, remediation systems, solo-builder workflows, and making open-source developer infrastructure more useful.

Good contribution areas:
- Backend product surfaces and API contracts
- MCP tools and roadmap completion
- Security and remediation workflows
- Frontend / product UX polish
- Docs, onboarding, and contributor ergonomics

## Getting Started

1. Fork the repository
2. Create a branch from `main`
3. Make your changes (keep commits small and independent)
4. Open a pull request against `main`

## Development Environment

See [Getting Started](Getting-Started) for prerequisites and local setup instructions.

### Build

```bash
make dep && make tools

cd web && yarn install && yarn build && cd ..

make build
```

### Run

```bash
./gitness server .local.env
```

## Coding Standards

### Backend (Go)

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Follow the formatting and style section of Peter Bourgon's [Go: Best Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style)

### Frontend (TypeScript)

- Follow the [Google TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html)

## Pre-Commit Hook

A pre-commit hook checks for required binaries (`grep`, `sed`, `xargs`) and runs Go-specific checks. If issues are found, the commit is halted until they are resolved.

## Linting

The CI pipeline runs separate lint checks for Go and TypeScript. Run locally before submitting a PR:

```bash
# Go
golangci-lint run ./...

# TypeScript
cd web && yarn lint
```

## Pull Request Checklist

- Branch from `main` and rebase if needed before opening
- Commits should be as small as possible while still compiling and passing tests
- Add tests relevant to the fixed bug or new feature
- Update or add documentation if API contracts change

## Dependency Management (Go)

SoloDev uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more).

```bash
# Add or update a dependency
go get example.com/some/module/pkg@vX.Y.Z

# Tidy
GO111MODULE=on go mod tidy
```

Commit both `go.mod` and `go.sum` changes.

## Database Migrations

New database tables require migrations in both backends:
- `app/store/database/migrate/postgres/NNNN_*.{up,down}.sql`
- `app/store/database/migrate/sqlite/NNNN_*.{up,down}.sql`

Migration numbers are sequential. Current SoloDev-specific range: `0102`–`0172`.

## Adding a New Module

To add a new SoloDev product module, follow the existing pattern:

1. **Types** — add structs in `types/<module>.go`
2. **Enums** — add enum types in `types/enum/<module>.go` (if needed)
3. **Store interfaces** — add to `app/store/database.go`
4. **Store implementation** — add `app/store/database/<module>.go`
5. **Migrations** — add SQL migration files
6. **Controller** — add `app/api/controller/<module>/controller.go`
7. **Handlers** — add `app/api/handler/<module>/*.go`
8. **Request helpers** — add `app/api/request/<module>.go` (if needed)
9. **Events** — add `app/events/<module>/` (events.go, reporter.go, reader.go)
10. **Router** — register routes in `app/router/api_modules.go`
11. **Wire** — add to dependency injection setup
12. **MCP tools** — add to `mcp/tools_atomic.go` and update `mcp/tools_compound.go` if needed

## License

Contributions are licensed under Apache License 2.0. All submissions must be your own work or properly attributed.

By contributing, you agree that your contributions may be distributed under the Apache-2.0 license.

> SoloDev is built from a fork of [Gitness by Harness](https://github.com/harness/gitness). Upstream attribution is preserved and derivative-work notices are maintained in [NOTICE](https://github.com/EolaFam1828/SoloDev/blob/main/NOTICE).
