# Development Setup

How to set up a local development environment for contributing to SoloDev.

## Prerequisites

| Requirement | Version |
|-------------|---------|
| Go | 1.20+ |
| Node.js | 16+ (latest stable) |
| Yarn | 1.x or 3.x |
| protoc | 3.21.11 |
| Docker | 20.10+ (for pipeline/Gitspace testing) |

## Setup Steps

### 1. Clone the Repository

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
```

### 2. Install Protobuf Tools

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

### 3. Install Dependencies

```bash
make dep && make tools
```

### 4. Build Frontend

```bash
cd web && yarn install && yarn build && cd ..
```

### 5. Build Backend

```bash
make build
```

### 6. Run

```bash
./gitness server .local.env
```

Open [http://localhost:3000](http://localhost:3000).

## Dependency Management

SoloDev uses Go modules:

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

## Linting

```bash
# Go
golangci-lint run ./...

# TypeScript
cd web && yarn lint
```

## Testing

```bash
# Go tests
go test ./...

# MCP server tests
go test ./mcp/... -v

# Frontend
cd web && yarn test
```

## Pre-Commit Hook

A pre-commit hook checks for required binaries and runs Go-specific checks. If issues are found, the commit is halted.
