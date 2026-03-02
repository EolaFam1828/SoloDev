# Zero-Config Auto-Pipeline Module

## Overview

The Auto-Pipeline module eliminates the need to manually write `.harness/pipeline.yaml` files. It analyzes repository file paths to detect the technology stack (language, framework, build tool), then generates a sensible default pipeline YAML with lint, test, and build steps — ready to run immediately.

This is designed for a solo developer workflow: push code, get CI instantly, no config ceremony.

## Architecture

### Components

#### 1. Types (`types/auto_pipeline.go`)
- **DetectedStack**: The result of repo analysis — languages, frameworks, build tool, entry point, and booleans for tests/Docker/existing CI
- **LanguageInfo**: Per-language breakdown with percentage and primary flag
- **AutoPipelineConfig**: The generated pipeline output — identifier, YAML string, detected stack, timestamp
- **GeneratePipelineInput**: API request body (optional branch override)

#### 2. Detection Engine (`app/pipeline/autopipeline/detect.go`)

##### `DetectStack(files []string) DetectedStack`
Lightweight static analysis of file paths (no network calls):

| Detection | Method |
|-----------|--------|
| **Language** | File extension matching (.go, .py, .ts/.tsx, .js/.jsx, .rs, .java, .rb, .cs, .swift) |
| **Build tool** | Filename matching (go.mod, package.json, Cargo.toml, pom.xml, etc.) |
| **Framework** | Config files (next.config.js, vite.config.ts, angular.json, nuxt.config.ts) |
| **Docker** | Presence of Dockerfile |
| **Existing CI** | `.harness/` or `.github/workflows/` directories |
| **Tests** | `_test.go`, `test_`, `.test.`, `.spec.`, `/tests/`, `/__tests__/` |

##### `GeneratePipelineYAML(stack DetectedStack) string`
Generates pipeline YAML based on the primary language:

| Language | Pipeline Steps |
|----------|---------------|
| **Go** | golangci-lint → go test -race → go build |
| **Node/TypeScript** | npm ci → npm lint → npm test → npm build |
| **Python** | pip install → ruff lint → pytest |
| **Rust** | cargo clippy → cargo test → cargo build --release |
| **Generic** | Echo message + optional Docker |

All pipelines add a Docker build step if `HasDocker == true`.

#### 3. Controller (`app/api/controller/autopipeline/controller.go`)
Business logic with space-scoped authorization:
- `GenerateAutoConfig()`: Accepts file paths, runs detection, returns generated config

#### 4. Handler (`app/api/handler/autopipeline/generate.go`)
- `HandleGenerate()`: POST handler accepting `{"files": [...]}` and returning the generated pipeline

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
  "yaml": "kind: pipeline\ntype: docker\nname: auto-go\nsteps:\n    - step:\n        type: run\n        name: lint\n        spec:\n          image: golangci/golangci-lint\n          script: |\n            golangci-lint run ./...\n    ...",
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

## Generated Pipeline Example (Go)

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
    - step:
        type: plugin
        name: docker
        spec:
          image: plugins/docker
```

## File Structure

```
├── types/
│   └── auto_pipeline.go
├── app/
│   ├── pipeline/autopipeline/
│   │   └── detect.go               # DetectStack + GeneratePipelineYAML
│   ├── api/
│   │   ├── controller/autopipeline/
│   │   │   └── controller.go
│   │   └── handler/autopipeline/
│   │       └── generate.go
└── app/router/
    └── api_modules.go              # setupAutoPipeline()
```

## Extending Detection

To add support for a new language/framework:

1. Add extension matching in `DetectStack()` switch case
2. Add build tool detection if needed
3. Create a new `generateXxxPipeline()` function
4. Add the case to `GeneratePipelineYAML()` switch
