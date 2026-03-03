# Contributing to SoloDev

Thank you for your interest in contributing to SoloDev. GitHub is the source of truth for the project history, and this repo is configured to keep `main` linear and low-friction.

* If you are a new contributor see: [Steps to Contribute](#steps-to-contribute)

* If you have a minor fix or improvement, feel free to create a pull request. Please provide necessary details in the pull request description and use a meaningful title.

* If you plan to do something more involved, first discuss your ideas by [raising an issue](https://github.com/EolaFam1828/SoloDev/issues). This will avoid unnecessary work and reduce duplicated effort.

* Relevant coding style guidelines are 

    - For backend: the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) and the formatting and style section of Peter Bourgon's [Go: Best Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style)
    - For frontend: [Google TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html) and [Best practices for Typescript coding](https://medium.com/@eshagarg1996/best-practices-for-typescript-coding-8b1ea98d02f8). 

* Be explicit about behavior changes, tests, and any API or UX impact in your pull request description.

## Steps to Contribute

Should you wish to work on an issue, please claim it first by commenting on the GitHub issue that you want to work on. This is to prevent duplicated efforts from contributors on the same issue.

Please check the repository issues for work that is ready to pick up. If you are unsure whether something is already in motion, open or comment on an issue first.

### Local Development

Please review the root [README.md](README.md) to build and test the project locally.

### Local Hooks

This fork relies on local verification instead of GitHub Actions. Run `make init` once per clone to install the tracked repo hooks and local Git defaults, and use `make verify` to run the same validation on demand.

The pre-push hook runs local verification before changes leave your machine. If you intentionally need to bypass it, use `git push --no-verify`.

`make init` also configures this clone for a linear-history workflow:

- `pull.ff=only`
- `pull.rebase=false`
- `merge.ff=only`
- `fetch.prune=true`
- `push.default=simple`
- `push.autoSetupRemote=true`
- `git sync` alias for fetch + fast-forward pull

### Verification

`make verify` runs:

- Go lint checks aligned with the repo's local gate
- web typecheck, lint, prettier, tests, and production build

Run it before pushing changes.

## Git Workflow

Use a short-lived branch for non-trivial work and keep `main` clean:

```bash
make init
git sync
git switch -c codex/my-change
```

When the work is ready:

- run `make verify`
- push the branch
- use a squash merge PR if you want review history on GitHub
- let GitHub delete the branch after merge

Do not merge `main` into your branch. If your branch falls behind, either fast-forward from a fresh base or recreate the branch from the current `main`.

## Pull Request Checklist

* Branch from `main` after running `git sync`. Keep history linear. GitHub only allows squash PR merges and `main` rejects merge-commit history.

* Commits should be as small as possible, while ensuring that each commit is correct independently (i.e., each commit should compile and pass tests).

* If your patch is not getting reviewed or you need a specific person to review it, you can @-reply a reviewer asking for a review in the pull request or a comment.

* Add tests relevant to the fixed bug or new feature.

## Dependency management

Harness uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies on external packages.

To add or update a new dependency, use the `go get` command:

```bash
# Pick the latest tagged release.
go get example.com/some/module/pkg@latest

# Pick a specific version.
go get example.com/some/module/pkg@vX.Y.Z
```

Tidy up the `go.mod` and `go.sum` files:

```bash
# The GO111MODULE variable can be omitted when the code isn't located in GOPATH.
GO111MODULE=on go mod tidy
```

You have to commit the changes to `go.mod` and `go.sum` before submitting the pull request.
