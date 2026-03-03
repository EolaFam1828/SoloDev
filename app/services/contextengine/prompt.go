// Copyright 2023 Harness, Inc.
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

package contextengine

import (
	"fmt"
	"strings"
)

const systemPrompt = `You are an expert software engineer specializing in automated code remediation.
Your task is to analyze the provided context fragments and generate a precise fix as a unified diff patch.

Each context fragment is labeled with its provenance (where it came from). Use this to gauge reliability:
- "error_log": Direct error output from CI/runtime — highest signal.
- "git_fetch": Source code fetched from the repository at the specified commit — exact state.
- "user_input": Context supplied by the operator — treat as advisory.
- "security_finding": Structured data from a security scanner — high signal for vulnerability fixes.
- "vector_search": Related code retrieved by similarity search — useful for understanding patterns and dependencies.

RULES:
1. Output EXACTLY one unified diff block wrapped in triple backticks with the "diff" language tag.
2. The diff must be directly applicable via "patch -p1".
3. Include a CONFIDENCE line at the very end: "CONFIDENCE: 0.X" where X is your confidence (0.0 to 1.0).
4. Keep changes minimal — fix only what is broken.
5. Do not add unrelated improvements or refactoring.
6. If you cannot determine a fix, output an empty diff block and CONFIDENCE: 0.0.

EXAMPLE OUTPUT FORMAT:
` + "```diff" + `
--- a/path/to/file.go
+++ b/path/to/file.go
@@ -10,3 +10,3 @@
-    broken line
+    fixed line
` + "```" + `
CONFIDENCE: 0.85`

// BuildPromptFromBundle constructs the user prompt from a structured context bundle.
func BuildPromptFromBundle(bundle *ContextBundle) string {
	var b strings.Builder

	b.WriteString("## Remediation Context\n\n")

	if bundle.FilePath != "" {
		fmt.Fprintf(&b, "**Target File:** `%s`\n", bundle.FilePath)
	}
	if bundle.Branch != "" {
		fmt.Fprintf(&b, "**Branch:** `%s`\n", bundle.Branch)
	}
	if bundle.CommitSHA != "" {
		fmt.Fprintf(&b, "**Commit:** `%s`\n", bundle.CommitSHA)
	}
	if bundle.TriggerSource != "" {
		fmt.Fprintf(&b, "**Trigger:** %s", bundle.TriggerSource)
		if bundle.TriggerRef != "" {
			fmt.Fprintf(&b, " (ref: %s)", bundle.TriggerRef)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n---\n\n")

	for i, frag := range bundle.Fragments {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "### %s\n", frag.Label)
		fmt.Fprintf(&b, "_Source: %s_", string(frag.Source))
		if frag.FilePath != "" {
			fmt.Fprintf(&b, " | _File: `%s`_", frag.FilePath)
		}
		if frag.TrimmedBytes > 0 {
			fmt.Fprintf(&b, " | _Trimmed: %d bytes_", frag.TrimmedBytes)
		}
		b.WriteString("\n\n```\n")
		b.WriteString(frag.Content)
		b.WriteString("\n```\n")
	}

	b.WriteString("\n---\n\n")
	b.WriteString("Analyze the context fragments above and generate a unified diff patch to fix the issue.\n")

	return b.String()
}

// GetSystemPrompt returns the system prompt for AI remediation.
func GetSystemPrompt() string {
	return systemPrompt
}
