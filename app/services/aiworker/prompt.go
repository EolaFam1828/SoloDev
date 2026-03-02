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

package aiworker

import (
	"fmt"
	"strings"

	"github.com/EolaFam1828/SoloDev/types"
)

const systemPrompt = `You are an expert software engineer specializing in automated code remediation.
Your task is to analyze error logs and source code, then generate a precise fix as a unified diff patch.

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

// BuildUserPrompt constructs the user prompt from a remediation record.
func BuildUserPrompt(rem *types.Remediation) string {
	var b strings.Builder

	b.WriteString("## Error Context\n\n")

	if rem.FilePath != "" {
		fmt.Fprintf(&b, "**File:** `%s`\n", rem.FilePath)
	}
	if rem.Branch != "" {
		fmt.Fprintf(&b, "**Branch:** `%s`\n", rem.Branch)
	}
	if rem.CommitSHA != "" {
		fmt.Fprintf(&b, "**Commit:** `%s`\n", rem.CommitSHA)
	}

	b.WriteString("\n### Error Log\n\n```\n")
	b.WriteString(rem.ErrorLog)
	b.WriteString("\n```\n\n")

	if rem.SourceCode != "" {
		b.WriteString("### Source Code\n\n```\n")
		b.WriteString(rem.SourceCode)
		b.WriteString("\n```\n\n")
	}

	b.WriteString("Please analyze the error and generate a unified diff patch to fix it.\n")

	return b.String()
}

// GetSystemPrompt returns the system prompt for AI remediation.
func GetSystemPrompt() string {
	return systemPrompt
}
