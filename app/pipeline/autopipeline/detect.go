// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package autopipeline

import (
	"fmt"
	"strings"

	"github.com/EolaFam1828/SoloDev/types"
)

// DetectStack analyzes a list of repo file paths and returns the detected technology stack.
// This is a lightweight static analysis that avoids any network calls.
func DetectStack(files []string) types.DetectedStack {
	stack := types.DetectedStack{}

	langHits := map[string]int{}
	total := 0

	for _, f := range files {
		lower := strings.ToLower(f)

		// Language detection
		switch {
		case strings.HasSuffix(lower, ".go"):
			langHits["go"]++
			total++
		case strings.HasSuffix(lower, ".py"):
			langHits["python"]++
			total++
		case strings.HasSuffix(lower, ".ts") || strings.HasSuffix(lower, ".tsx"):
			langHits["typescript"]++
			total++
		case strings.HasSuffix(lower, ".js") || strings.HasSuffix(lower, ".jsx"):
			langHits["javascript"]++
			total++
		case strings.HasSuffix(lower, ".rs"):
			langHits["rust"]++
			total++
		case strings.HasSuffix(lower, ".java"):
			langHits["java"]++
			total++
		case strings.HasSuffix(lower, ".rb"):
			langHits["ruby"]++
			total++
		case strings.HasSuffix(lower, ".cs"):
			langHits["csharp"]++
			total++
		case strings.HasSuffix(lower, ".swift"):
			langHits["swift"]++
			total++
		}

		// Framework / build tool detection
		base := lower
		if idx := strings.LastIndex(lower, "/"); idx >= 0 {
			base = lower[idx+1:]
		}

		switch base {
		case "go.mod":
			stack.BuildTool = "go"
		case "package.json":
			stack.BuildTool = "npm"
		case "requirements.txt", "pyproject.toml", "setup.py":
			stack.BuildTool = "pip"
		case "cargo.toml":
			stack.BuildTool = "cargo"
		case "pom.xml":
			stack.BuildTool = "maven"
		case "build.gradle", "build.gradle.kts":
			stack.BuildTool = "gradle"
		case "gemfile":
			stack.BuildTool = "bundler"
		case "makefile":
			if stack.BuildTool == "" {
				stack.BuildTool = "make"
			}
		}

		// Framework detection
		switch base {
		case "next.config.js", "next.config.mjs", "next.config.ts":
			stack.Frameworks = appendUnique(stack.Frameworks, "nextjs")
		case "vite.config.ts", "vite.config.js":
			stack.Frameworks = appendUnique(stack.Frameworks, "vite")
		case "angular.json":
			stack.Frameworks = appendUnique(stack.Frameworks, "angular")
		case "nuxt.config.ts", "nuxt.config.js":
			stack.Frameworks = appendUnique(stack.Frameworks, "nuxt")
		}

		// Docker
		if base == "dockerfile" || strings.HasPrefix(base, "dockerfile.") {
			stack.HasDocker = true
		}

		// Existing CI
		if strings.Contains(lower, ".harness/") || strings.Contains(lower, ".github/workflows/") {
			stack.HasCI = true
		}

		// Tests
		if strings.Contains(lower, "_test.go") ||
			strings.Contains(lower, "test_") ||
			strings.Contains(lower, ".test.") ||
			strings.Contains(lower, ".spec.") ||
			strings.Contains(lower, "/tests/") ||
			strings.Contains(lower, "/__tests__/") {
			stack.HasTests = true
		}
	}

	// Build language list
	maxLang := ""
	maxCount := 0
	for lang, count := range langHits {
		pct := 0.0
		if total > 0 {
			pct = float64(count) / float64(total) * 100
		}
		info := types.LanguageInfo{
			Name:       lang,
			Percentage: pct,
		}
		if count > maxCount {
			maxCount = count
			maxLang = lang
		}
		stack.Languages = append(stack.Languages, info)
	}

	// Mark primary language
	for i := range stack.Languages {
		if stack.Languages[i].Name == maxLang {
			stack.Languages[i].Primary = true
		}
	}

	return stack
}

// GeneratePipelineYAML generates a Harness pipeline YAML for the detected stack.
func GeneratePipelineYAML(stack types.DetectedStack) string {
	primary := ""
	for _, l := range stack.Languages {
		if l.Primary {
			primary = l.Name
			break
		}
	}

	switch primary {
	case "go":
		return generateGoPipeline(stack)
	case "typescript", "javascript":
		return generateNodePipeline(stack)
	case "python":
		return generatePythonPipeline(stack)
	case "rust":
		return generateRustPipeline(stack)
	default:
		return generateGenericPipeline(stack)
	}
}

func generateGoPipeline(stack types.DetectedStack) string {
	steps := []string{
		step("lint", "golangci/golangci-lint", "golangci-lint run ./..."),
		step("test", "golang", "go test -race -coverprofile=coverage.out ./..."),
		step("build", "golang", "go build -o /dev/null ./..."),
	}
	if stack.HasDocker {
		steps = append(steps, step("docker", "plugins/docker", ""))
	}
	return wrapPipeline("auto-go", steps)
}

func generateNodePipeline(stack types.DetectedStack) string {
	install := "npm ci"
	steps := []string{
		step("install", "node", install),
		step("lint", "node", "npm run lint || true"),
		step("test", "node", "npm test || true"),
		step("build", "node", "npm run build"),
	}
	if stack.HasDocker {
		steps = append(steps, step("docker", "plugins/docker", ""))
	}
	return wrapPipeline("auto-node", steps)
}

func generatePythonPipeline(stack types.DetectedStack) string {
	steps := []string{
		step("install", "python", "pip install -r requirements.txt || true"),
		step("lint", "python", "pip install ruff && ruff check . || true"),
		step("test", "python", "python -m pytest || true"),
	}
	if stack.HasDocker {
		steps = append(steps, step("docker", "plugins/docker", ""))
	}
	return wrapPipeline("auto-python", steps)
}

func generateRustPipeline(stack types.DetectedStack) string {
	steps := []string{
		step("lint", "rust", "cargo clippy -- -D warnings || true"),
		step("test", "rust", "cargo test"),
		step("build", "rust", "cargo build --release"),
	}
	return wrapPipeline("auto-rust", steps)
}

func generateGenericPipeline(stack types.DetectedStack) string {
	steps := []string{
		step("echo", "alpine", "echo 'No auto-detection matched. Configure your pipeline manually.'"),
	}
	if stack.HasDocker {
		steps = append(steps, step("docker", "plugins/docker", ""))
	}
	return wrapPipeline("auto-generic", steps)
}

func step(name, image, commands string) string {
	if commands == "" {
		return fmt.Sprintf(`    - step:
        type: plugin
        name: %s
        spec:
          image: %s`, name, image)
	}
	return fmt.Sprintf(`    - step:
        type: run
        name: %s
        spec:
          image: %s
          script: |
            %s`, name, image, commands)
}

func wrapPipeline(name string, steps []string) string {
	return fmt.Sprintf(`kind: pipeline
type: docker
name: %s
steps:
%s
`, name, strings.Join(steps, "\n"))
}

func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}
