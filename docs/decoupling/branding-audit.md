# SoloDev Decoupling Audit

Status as of March 2, 2026, after the canonical env and packaging sweep.

## Completed In Phase 2

- Primary Go module path moved from `github.com/harness/gitness` to `github.com/EolaFam1828/SoloDev`.
- Primary command/build target moved to `cmd/solodev`.
- `make build`, Docker build, and Wire generation now target `cmd/solodev`.
- Compatibility shims remain in place for the old `gitness` binary/script surface.

## Remaining Inventory

- `GITNESS_*` tokens: `49`
- `Harness` tokens: `4640`
- `cmd/gitness` refs: `7`
- Chart/package legacy refs: `15`

Important: the `Harness` count is mostly copyright/provenance text. Do not bulk-remove those occurrences.

## Must Rename Soon

These are still part of the active product or operator surface and should be removed before calling the platform fully decoupled.

- Legacy env compatibility still exists intentionally in [Dockerfile](../../Dockerfile), [Dockerfile.uiv2](../../Dockerfile.uiv2), and [registry/app/remote/clients/registry/client.go](../../registry/app/remote/clients/registry/client.go).
- Compatibility source tree still exists in [cmd/gitness](../../cmd/gitness).
- Some internal/legal/compliance docs still reference Harness by design in [docs/compliance/apache-2.0-derivative-checklist.md](../../docs/compliance/apache-2.0-derivative-checklist.md).

## Compatibility Layer Kept Intentionally

These are acceptable to keep temporarily because they preserve existing local workflows or deployments.

- Old command tree in [cmd/gitness](../../cmd/gitness) remains as a compatibility source path.
- Old wire wrapper in [scripts/wire/gitness.sh](../../scripts/wire/gitness.sh) now forwards to the new generator.
- Old binary symlink is preserved in [Makefile](../../Makefile) and [Dockerfile](../../Dockerfile).
- `SOLODEV_*` to `GITNESS_*` aliasing remains in [cli/operations/server/config.go](../../cli/operations/server/config.go).

## Keep For Legal Provenance

These should not be scrubbed as part of branding cleanup.

- Apache-2.0 `LICENSE`
- `NOTICE`
- Existing upstream copyright headers
- `Modified by ...` notices added to materially changed files

## Recommended Next Pass

1. Reduce the remaining compatibility-only `GITNESS_*` fallbacks in container builds and registry init paths once old deploys no longer need them.
2. Collapse [cmd/gitness](../../cmd/gitness) into a thinner shim once the new command path has baked for a while.
3. Re-run codegen where needed so generated artifacts stay aligned with the renamed module/package surface.
