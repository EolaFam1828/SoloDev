# Contributing to SoloDev

Thank you for your interest in contributing to SoloDev. SoloDev uses GitHub to manage reviews of pull requests.

* If you are a new contributor see: [Steps to Contribute](#steps-to-contribute)

* If you have a minor fix or improvement, feel free to create a pull request. Please provide necessary details in the pull request description and use a meaningful title.

* If you plan to do something more involved, first discuss your ideas by [raising an issue](https://github.com/EolaFam1828/SoloDev/issues). This will avoid unnecessary work and usually prevents duplicated effort.

* Relevant coding style guidelines are 

    - For backend: the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) and the formatting and style section of Peter Bourgon's [Go: Best Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style)
    - For frontend: [Google TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html) and [Best practices for Typescript coding](https://medium.com/@eshagarg1996/best-practices-for-typescript-coding-8b1ea98d02f8). 

* If this repository requires a contributor agreement, GitHub will prompt you during the pull request flow.

## Steps to Contribute

Should you wish to work on an issue, please claim it first by commenting on the GitHub issue that you want to work on. This is to prevent duplicated efforts from contributors on the same issue.

Please check the [`good-first-issue`](https://github.com/EolaFam1828/SoloDev/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) label to find issues that are good for getting started. If you have questions, comment on the issue or open a new one.

### Local Development

Please review the local [README](./README.md) and [Makefile](./Makefile) to build and test the project locally.

### Pre-commit Hook

We have a pre-commit hook to ensure code quality before committing changes. This hook checks for required binaries (grep, sed, and xargs) and runs checks specifically for Go files (*.go). If any issues are found during the checks, the commit process will be halted until the issues are resolved.

### Lint Check

Our CI linter pipeline conducts automated checks for code quality, with separate lint checks for Go and TypeScript in [ci-lint.yml](./.github/workflows/ci-lint.yml). These checks help ensure adherence to coding standards and identify potential issues early.

## Pull Request Checklist

* Branch from the main branch and, if needed, rebase to the current main branch before submitting your pull request. If it doesn't merge cleanly with main you may be asked to rebase your changes.

* Commits should be as small as possible, while ensuring that each commit is correct independently (i.e., each commit should compile and pass tests).

* If your patch is not getting reviewed or you need a specific person to review it, you can @-reply a reviewer asking for a review in the pull request or a comment.

* Add tests relevant to the fixed bug or new feature.

## Dependency management

SoloDev uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies on external packages.

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

