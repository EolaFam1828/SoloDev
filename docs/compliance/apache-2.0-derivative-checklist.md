# Apache-2.0 Derivative Checklist

Working checklist for SoloDev as a derivative of upstream Harness Open Source. This is engineering guidance, not legal advice.

## Core Rules

- Keep the root `LICENSE` file in source distributions.
- Keep an accurate `NOTICE` file for what SoloDev ships.
- Preserve upstream copyright and license notices in retained source files.
- Mark materially modified files when practical and continue preserving `Modified by ...` notices already added in this repo.
- Do not use Harness names, logos, or marks in product branding, packaging, or marketing in a way that suggests endorsement.

## Source Distribution Checklist

- Ship `LICENSE`.
- Ship `NOTICE`.
- Keep Apache headers in retained upstream-derived files.
- Keep attribution history when copying or refactoring upstream code into new packages.
- Document major SoloDev-specific modifications in release notes or a changelog.

## Binary And Container Distribution Checklist

- Include `LICENSE` and `NOTICE` in the image, release archive, or adjacent release assets.
- Verify image/package names no longer use Harness trademarks unless explicitly preserved only for compatibility.
- Verify entrypoints, env var examples, and docs prefer SoloDev names.
- Keep compatibility aliases documented as transitional, not canonical.

## Repo Hygiene Checklist

- Separate branding cleanup from legal provenance cleanup.
- Do not mass-delete `Harness` strings without checking whether they are legal notices.
- Track remaining compatibility surfaces in [docs/decoupling/branding-audit.md](../decoupling/branding-audit.md).
- Review generated artifacts after regeneration because old module/type names can leak back in from templates or codegen.

## Notice Maintenance Checklist

- Confirm `NOTICE` still attributes upstream Harness work.
- Add SoloDev attribution only for net-new or materially modified portions where appropriate.
- Remove stale product-marketing text from `NOTICE`; keep it factual.
- Recheck `NOTICE` whenever Helm charts, containers, or bundled third-party assets change.

## Trademark Checklist

- Remove Harness logos from UI, docs, and packaging.
- Replace Harness product names in setup instructions and screenshots.
- Keep plain factual attribution such as "derived from Harness Open Source" where useful.
- Avoid language implying official Harness partnership, approval, or distribution.

## Release Gate Before Public Shipping

1. Review `LICENSE` and `NOTICE`.
2. Review [docs/decoupling/branding-audit.md](../decoupling/branding-audit.md) and clear all "Must Rename Soon" items intended for that release.
3. Verify the UI, README, Helm chart, Docker image metadata, and CLI help output are SoloDev-branded.
4. Run a repo grep for `Harness`, `GITNESS_`, and `gitness` and manually classify anything user-visible.
