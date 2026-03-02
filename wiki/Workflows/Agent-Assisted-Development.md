# Agent-Assisted Development

A demonstration of how AI agents interact with SoloDev during development using the MCP protocol.

## Scenario

A developer has configured Claude Desktop to connect to SoloDev via MCP. They are working on a feature and want the agent to handle operational tasks while they focus on writing code.

## Setup

Claude Desktop configuration (`~/.config/claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "solodev": {
      "command": "gitness",
      "args": ["mcp", "stdio"],
      "env": {
        "SOLODEV_TOKEN": "<your-personal-access-token>"
      }
    }
  }
}
```

## Workflow 1: Onboard a New Repository

Developer: "Set up my new project in SoloDev."

Agent uses `onboard_repo` compound tool:
1. Triggers a security scan on the repository
2. Generates a CI/CD pipeline from detected stack
3. Creates quality gate rules for coverage, complexity, and style
4. Sets up a health check for the main API endpoint

Result: The repository has CI/CD, security scanning, quality enforcement, and health monitoring — all configured without manual YAML editing.

## Workflow 2: Fix a Failing Pipeline

Developer: "My pipeline is failing. Can you fix it?"

Agent workflow:
1. Reads `solodev://errors/active` resource — sees pipeline failure error
2. Calls `fix_this(error_id="pipeline-42-failure")`
3. Compound tool triggers remediation, waits for AI Worker to generate patch
4. Calls `remediation_get()` to retrieve the completed patch
5. Presents the diff to the developer for review

Result: The agent diagnosed the problem and generated a fix without the developer reading build logs.

## Workflow 3: Pre-PR Quality Check

Developer: "Is my code ready for PR?"

Agent uses `pr_ready` compound tool:
1. Runs a security scan — 0 critical, 0 high findings
2. Evaluates quality gates — all rules pass
3. Checks tech debt — 2 medium items (not blocking)
4. Returns verdict: "PR ready. 2 medium tech debt items noted."

Result: Comprehensive quality check without the developer running separate tools.

## Workflow 4: Incident Investigation

Developer: "Users are reporting slow responses."

Agent uses `incident_triage` compound tool:
1. Checks `solodev://health/status` — API endpoint showing degraded response times
2. Checks `solodev://errors/active` — database timeout errors appearing
3. Checks `solodev://security/open-findings` — no related security issues
4. Correlates signals: health degradation started when database timeout errors appeared
5. Reports: "Health degradation correlates with database timeout errors. Recommending remediation of the connection handling code."

Agent follows up with `fix_this` to trigger remediation.

## Workflow 5: Sprint Planning for Tech Debt

Developer: "Help me plan which tech debt to tackle this sprint."

Agent uses `solodev_debt_sprint` prompt:
1. Reads `solodev://tech-debt/hotspots` — sees all debt items
2. Prioritizes by severity, age, and file impact
3. Produces a ranked list with effort estimates
4. Suggests: "Focus on the 3 critical items in the database layer. Estimated total effort: 6 hours."

## Key Observation

In each workflow, the developer describes what they need in natural language. The agent translates that into MCP tool calls, orchestrates the platform operations, and returns results. The developer supervises rather than executes.
