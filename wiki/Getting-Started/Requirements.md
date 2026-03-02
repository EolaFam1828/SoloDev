# Requirements

Prerequisites for running SoloDev.

## Docker Path (Recommended)

| Requirement | Minimum Version |
|-------------|----------------|
| Docker Engine | 20.10+ |
| Docker Compose | V2 |
| RAM | 4 GB allocated to Docker |
| Disk | 2 GB free for images and data |

This is the simplest path. No other dependencies are needed.

## Source Build Path

| Requirement | Minimum Version |
|-------------|----------------|
| Go | 1.20+ |
| Node.js | 16+ (latest stable recommended) |
| Yarn | 1.x or 3.x |
| protoc | 3.21.11 |
| RAM | 4 GB |
| Disk | 3 GB (Go modules, node_modules, build artifacts) |

### Additional Go Tools

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

## Operating Systems

| OS | Status |
|----|--------|
| Linux (x86_64) | Fully supported |
| macOS (Intel) | Supported |
| macOS (Apple Silicon) | Supported (cross-compilation may be slower) |
| Windows (WSL2) | Supported via Docker |

## Optional: LLM Provider

For AI remediation features, one of the following is required:

| Provider | Requirement |
|----------|------------|
| Anthropic | API key |
| OpenAI | API key |
| Google Gemini | API key |
| Ollama | Local installation (no API key needed) |

SoloDev functions without an LLM provider — remediation tasks will be created but not processed.

## Optional: PostgreSQL

For production deployment:

| Requirement | Minimum Version |
|-------------|----------------|
| PostgreSQL | 14+ |

SQLite is used by default for local development.
