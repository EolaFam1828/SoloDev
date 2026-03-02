# First Pipeline

Run your first CI/CD pipeline with SoloDev.

## Option 1: Manual Pipeline

Create `.harness/pipeline.yaml` in your repository:

```yaml
kind: pipeline
spec:
  stages:
    - type: ci
      spec:
        steps:
          - name: test
            type: run
            spec:
              container: golang:1.21
              script: go test ./...
```

Push the file and the pipeline runs automatically.

## Option 2: Auto-Pipeline

Let SoloDev detect your tech stack and generate a pipeline.

### Via API

```bash
TOKEN="<your-personal-access-token>"

curl -X POST http://localhost:3000/api/v1/spaces/<space>/auto-pipeline/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "files": [
      "go.mod",
      "main.go",
      "internal/handler/users.go",
      "internal/handler/users_test.go",
      "Dockerfile"
    ]
  }'
```

The response includes detected stack information and generated YAML:

```json
{
  "yaml": "kind: pipeline\ntype: docker\nname: auto-go\n...",
  "detected_stack": {
    "languages": [{"name": "go", "percentage": 100, "primary": true}],
    "has_tests": true,
    "has_docker": true,
    "build_tool": "go"
  }
}
```

### Via MCP

```
pipeline_generate(files=["go.mod", "main.go", ...])
```

## Verify

After pushing the pipeline file, check the Pipelines section in the SoloDev dashboard. The pipeline execution should appear with its status (running, success, or failure).

## Supported Stacks

Auto-Pipeline generates pipelines for:
- **Go** — lint (golangci-lint), test, build, optional Docker
- **Node/TypeScript** — install, lint, test, build, optional Docker
- **Python** — install, lint (ruff), test (pytest), optional Docker
- **Rust** — clippy, test, build

See [Auto-Pipeline](../Modules/Core-Platform/Pipelines) for full detection details.
