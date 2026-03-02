# Coding Standards

Style and architectural rules for SoloDev contributions.

## Backend (Go)

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Follow the formatting and style section of Peter Bourgon's [Go: Best Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style)
- Run `golangci-lint run ./...` before submitting

## Frontend (TypeScript)

- Follow the [Google TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html)
- Run `cd web && yarn lint` before submitting

## Module Pattern

All SoloDev modules follow the same structure:

| Layer | Path Pattern |
|-------|-------------|
| Types | `types/<module>.go` |
| Enums | `types/enum/<module>.go` |
| Store interfaces | `app/store/database.go` |
| Store implementation | `app/store/database/<module>.go` |
| Migrations | `app/store/database/migrate/{postgres,sqlite}/NNNN_*` |
| Controller | `app/api/controller/<module>/controller.go` |
| Handlers | `app/api/handler/<module>/*.go` |
| Request helpers | `app/api/request/<module>.go` |
| Events | `app/events/<module>/` |
| Router registration | `app/router/api_modules.go` |
| Wire setup | Dependency injection in `cmd/gitness/wire.go` |
| MCP tools | `mcp/tools_atomic.go`, `mcp/tools_compound.go` |

## Adding a New Module

1. Define types in `types/<module>.go`
2. Add store interfaces to `app/store/database.go`
3. Implement the store in `app/store/database/<module>.go`
4. Write SQL migrations for PostgreSQL and SQLite
5. Create the controller in `app/api/controller/<module>/`
6. Create handlers in `app/api/handler/<module>/`
7. Add events in `app/events/<module>/`
8. Register routes in `app/router/api_modules.go`
9. Wire into dependency injection
10. Add MCP tools if the module should be agent-accessible

## Documentation Standards

Module documentation pages must include:
- Purpose
- Inputs
- Processing
- Outputs
- Status (Concept, Prototype, Implemented)
- Future Work

Keep descriptions mechanical, not marketing.
