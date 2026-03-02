# Zero-Config Auto-Pipeline

## Overview

The Auto-Pipeline module eliminates the need to manually write `.harness/pipeline.yaml` files. It analyzes a list of repository file paths, detects the technology stack (language, framework, build tool), and generates a sensible default pipeline YAML with lint, test, and build steps — ready to run immediately.

This is designed for a solo developer workflow: push code, get CI instantly, no config ceremony.

## Stack Detection

The detection engine (`app/pipeline/autopipeline/detect.go`) performs lightweight static analysis of file paths — no network calls.

### Language Detection

Languages are detected by file extension:

| Extension(s) | Language |
|-------------|---------|
| `.go` | Go |
| `.py` | Python |
| `.ts`, `.tsx` | TypeScript |
| `.js`, `.jsx` | JavaScript |
| `.rs` | Rust |
| `.java` | Java |
| `.rb` | Ruby |
| `.cs` | C# |
| `.swift` | Swift |

### Build Tool Detection

Build tools are detected by filename:

| File | Build Tool |
|------|-----------|
| `go.mod` | go |
| `package.json` | npm |
| `Cargo.toml` | cargo |
| `pom.xml` | maven |

### Framework Detection

Frameworks are detected by config file presence:

| File | Framework |
|------|---------|
| `next.config.js` | Next.js |
| `vite.config.ts` | Vite |
| `angular.json` | Angular |
| `nuxt.config.ts` | Nuxt.js |

### Additional Detection

| Signal | Detection Method |
|--------|----------------|
| Docker | Presence of `Dockerfile` |
| Existing CI | `.harness/` or `.github/workflows/` directories |
| Tests | `_test.go`, `test_`, `.test.`, `.spec.`, `/tests/`, `/__tests__/` patterns |

## Generated Pipelines by Language

The `GeneratePipelineYAML()` function generates steps based on the primary detected language:

### Go

```yaml
kind: pipeline
type: docker
name: auto-go
steps:
    - step:
        type: run
        name: lint
        spec:
          image: golangci/golangci-lint
          script: |
            golangci-lint run ./...
    - step:
        type: run
        name: test
        spec:
          image: golang
          script: |
            go test -race -coverprofile=coverage.out ./...
    - step:
        type: run
        name: build
        spec:
          image: golang
          script: |
            go build -o /dev/null ./...
    # + docker plugin step if Dockerfile detected
```

### Node / TypeScript

Steps: `npm ci` → `npm run lint || true` → `npm test || true` → `npm run build`
Docker step added if `Dockerfile` is detected.

### Python

Steps: `pip install -r requirements.txt || true` → `ruff check . || true` → `python -m pytest || true`
Docker step added if `Dockerfile` is detected.

### Rust

Steps: `cargo clippy -- -D warnings || true` → `cargo test` → `cargo build --release`
No Docker step is added for Rust.

### Generic (undetected stack)

Echo message step + Docker plugin step if `Dockerfile` is detected.

## API Endpoint

### Generate Auto-Pipeline

**POST** `/api/v1/spaces/{space_ref}/auto-pipeline/generate`

Request:
```json
{
  "files": [
    "go.mod",
    "main.go",
    "cmd/server/main.go",
    "internal/handler/users.go",
    "internal/handler/users_test.go",
    "Dockerfile"
  ]
}
```

Response:
```json
{
  "repo_id": 0,
  "identifier": "auto-pipeline",
  "yaml": "kind: pipeline\ntype: docker\nname: auto-go\n...",
  "detected_stack": {
    "languages": [{"name": "go", "percentage": 100, "primary": true}],
    "frameworks": [],
    "has_tests": true,
    "has_docker": true,
    "has_ci": false,
    "build_tool": "go"
  },
  "generated": 1709283600000
}
```

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/auto_pipeline.go` |
| Detection + generation engine | `app/pipeline/autopipeline/detect.go` |
| Controller | `app/api/controller/autopipeline/controller.go` |
| Handler | `app/api/handler/autopipeline/generate.go` |
| Router registration | `app/router/api_modules.go` (`setupAutoPipeline()`) |

## Extending Detection

To add support for a new language or framework:

1. Add extension matching in `DetectStack()` switch case
2. Add build tool detection if needed
3. Create a new `generateXxxPipeline()` function
4. Add the case to `GeneratePipelineYAML()` switch
