# Contributing to riteway-golang

## Prerequisites

- Go 1.21 or later
- Git with commit signing configured (see [Commit signing](#commit-signing))

## Setup

```sh
git clone https://github.com/your-username/riteway-golang
cd riteway-golang
```

No additional setup is needed — the module has one dependency (`go-cmp`) which
Go fetches automatically.

## Development

```sh
make check        # fmt + vet + test + test-scripts (mirrors CI)
make test         # go test -count=1 -race ./...
make fmt          # check formatting (fails if unformatted files exist)
make vet          # go vet ./...
make test-scripts # bash scripts/test-release.sh
```

To run a single test:

```sh
go test -run TestName ./...
```

## Commit signing

All commits must be signed. See [RELEASING.md](RELEASING.md) for GPG and SSH
signing setup instructions.

## Pull request process

1. Fork the repo and create a branch from `main`.
2. Make your changes and ensure `make check` passes locally.
3. Open a pull request with a clear title and description.

Before a PR can merge, the following are required:

- All CI jobs pass (`Lint`, `Test (Go 1.21–1.24)`, `Script tests`, `Hook tests`)
- One approving review
- All review comments resolved
- Branch up-to-date with `main` (use the "Update branch" button on the PR if needed)

PRs are merged via squash — your branch's commits are collapsed into one commit
on `main` using the PR title and description.

## Architecture

See [docs/decisions/](docs/decisions/) for architectural decision records explaining
key design choices, and [CLAUDE.md](CLAUDE.md) for a full overview of the codebase.
