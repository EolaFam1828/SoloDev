# Source Control

## Purpose

Provides Git hosting, code review, and repository management. This is the foundational layer that stores the code SoloDev operates on. Inherited from the Gitness upstream with no SoloDev-specific modifications to the core SCM logic.

## Inputs

- Git push operations (SSH on port 3022, HTTP on port 3000)
- Pull request creation and review actions
- Webhook configurations
- Branch and tag management operations

## Processing

- Stores Git objects and references using the native Git protocol
- Manages pull request state, comments, and merge operations
- Fires webhooks on repository events (push, PR create, PR merge)
- Enforces branch protection rules (when configured)

## Outputs

- Git repository content accessible via HTTP and SSH
- Pull request state changes and merge results
- Webhook payloads to registered endpoints
- Repository metadata for other modules (used by Auto-Pipeline, AI Worker)

## Integration with SoloDev

The SCM layer is the target for remediation outputs. When the AI Worker generates a patch diff:
1. The patch targets files in a hosted repository
2. Fix branches are created in the SCM
3. Pull requests are opened with the AI analysis as description (planned)
4. Auto-merge applies validated patches to the target branch (planned)

The Auto-Pipeline module reads repository file paths from SCM to detect the tech stack.

## Key Paths

| Purpose | Path |
|---------|------|
| Repository controller | `app/api/controller/repo/` |
| Git backend | `git/` |
| SSH server | `ssh/` |

## Status

**Implemented** — Inherited from Gitness. Full Git hosting, pull requests, code review, and webhooks.

## Future Work

- Auto-PR creation from AI-generated patches
- Auto-merge for high-confidence patches
- Tighter integration between SCM events and the remediation loop
