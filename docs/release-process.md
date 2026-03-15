# Release Process

## Prerequisites

- Push access to `github.com/mycargus/riteway-golang`
- Go 1.21+ installed
- Git configured with your identity
- `gh` CLI installed and authenticated (`gh auth login`)

## Before the first release

Ensure the repo exists on GitHub at the path matching the module declaration in `go.mod`:

```
github.com/mycargus/riteway-golang
```

## Releasing

**Step 1:** Update `CHANGELOG.md`

Move items from `[Unreleased]` into a new versioned section:

```markdown
## [0.1.0] - 2026-03-14

### Added
- ...
```

Commit the changelog update before running the release script.

**Step 2:** Run the release script

```sh
./scripts/release.sh v0.1.0
```

The script:
1. Validates the version format (`vMAJOR.MINOR.PATCH`)
2. Extracts release notes from the matching `CHANGELOG.md` section — fails if none found
3. Requires a clean working tree
4. Checks the tag does not already exist
5. Runs `make check` (fmt + vet + test with race detector)
6. Runs `go mod tidy` and fails if `go.sum` or `go.mod` changed
7. Shows a summary with the release notes and prompts for confirmation
8. Tags and pushes, then creates a GitHub release via `gh release create`

## Running checks without releasing

```sh
make check   # fmt + vet + test
make test    # tests only
make fmt     # format check only
make vet     # vet only
```

## Versioning

This module follows [Semantic Versioning](https://semver.org) with a `v` prefix.

| Change | Version bump | Example |
|---|---|---|
| Bug fix, doc improvement | patch | `v0.1.0` → `v0.1.1` |
| New function or field (backwards-compatible) | minor | `v0.1.0` → `v0.2.0` |
| Breaking API change (removed/renamed function, changed signature) | major | `v0.x.x` → `v1.0.0` |

**Important:** A major version bump past `v1` requires changing the module path in `go.mod` and all import paths:

```
module github.com/mycargus/riteway-golang/v2
```

Stay on `v0.x.x` while the API is still settling. Once the API is stable and you are ready to commit to no breaking changes, tag `v1.0.0`.

## After releasing

The script creates the GitHub release automatically and prints the pkg.go.dev URL.
Visit it to trigger documentation indexing (pkg.go.dev indexes on first request):

```
https://pkg.go.dev/github.com/mycargus/riteway-golang@v0.1.0
```

## Verifying the release

From any directory outside this repo, confirm the module is fetchable:

```sh
mkdir /tmp/verify-release && cd /tmp/verify-release
go mod init verify
go get github.com/mycargus/riteway-golang@v0.1.0
```

## Yanking a broken release

Go has no "unpublish". If a release is broken:

1. Tag a patch release immediately with the fix.
2. To discourage use of the broken version, add a `retract` directive to `go.mod`:

   ```
   retract v0.1.0 // reason: describes the problem
   ```

3. Release the new version. `go get` will warn users on the retracted version to upgrade.
