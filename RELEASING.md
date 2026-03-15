# Releasing riteway

Publishing is always triggered by a human — never by an AI agent.

`./scripts/release.sh vX.Y.Z` is the standard release command. It runs preflight
checks, creates an annotated git tag, and pushes it. GitHub Actions detects the
tag and handles creating the GitHub release.

## Prerequisites

- Push access to `github.com/mycargus/riteway-golang`
- Go 1.21+ installed
- `gh` CLI installed and authenticated (`gh auth login`)

### Commit signing setup

All commits must be signed (see ADR-012). Configure Git to sign automatically:

**GPG signing:**

```sh
gpg --list-secret-keys --keyid-format=long   # find your key ID
git config --global commit.gpgsign true
git config --global user.signingkey <your-key-id>
```

Register your public key on GitHub: **Settings → SSH and GPG keys → New GPG key**.

**SSH signing (modern alternative):**

```sh
git config --global commit.gpgsign true
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/id_ed25519.pub
```

Register your public key on GitHub: **Settings → SSH and GPG keys → New signing key**.

---

### One-time setup: release environment

GitHub Actions uses a `release` environment to gate tag-triggered runs. Configure
it once per repository:

1. **Settings → Environments → New environment** → name it `release`
2. Under **Deployment branches and tags**, change to **Selected branches and tags**
3. **Add deployment branch or tag rule** → enter `v*` as a **tag** pattern
4. Save

## Steps

1. Update `CHANGELOG.md` — move items from `[Unreleased]` into a new versioned
   section:

   ```markdown
   ## [X.Y.Z] - YYYY-MM-DD

   ### Added
   - ...
   ```

2. Commit and push to `main`:

   ```sh
   git add CHANGELOG.md
   git commit -m "Release vX.Y.Z"
   git push origin main
   ```

3. Run the release script:

   ```sh
   make release VERSION=vX.Y.Z
   ```

   This validates the version format, extracts release notes from `CHANGELOG.md`,
   requires a clean working tree, checks the tag does not already exist, runs
   `make check` (fmt + vet + test with race detector), verifies `go mod tidy`
   produces no changes, shows a summary with the release notes, and prompts for
   confirmation before tagging and pushing.

4. GitHub Actions takes over:
   - Verifies `CHANGELOG.md` has an entry for the tag
   - Runs lint and full test suite
   - Creates the GitHub release with the CHANGELOG notes

   Monitor progress at: <https://github.com/mycargus/riteway-golang/actions>

---

## Before the first release

Ensure the repo exists on GitHub at the path matching the module declaration in
`go.mod`:

```
github.com/mycargus/riteway-golang
```

---

## Versioning

Follow [Semantic Versioning](https://semver.org/) with a `v` prefix:

- **PATCH** (`v0.1.x`) — backwards-compatible bug fixes
- **MINOR** (`v0.x.0`) — new backwards-compatible features
- **MAJOR** (`vx.0.0`) — breaking changes

Stay on `v0.x.x` while the API is still settling. Once the API is stable and you
are ready to commit to no breaking changes, tag `v1.0.0`.

**Important:** A major version bump past `v1` requires changing the module path in
`go.mod` and all import paths:

```
module github.com/mycargus/riteway-golang/v2
```

---

## After releasing

The script prints the pkg.go.dev URL. Visit it to trigger documentation indexing
(pkg.go.dev indexes on first request):

```
https://pkg.go.dev/github.com/mycargus/riteway-golang@vX.Y.Z
```

---

## Verifying the release

From any directory outside this repo, confirm the module is fetchable:

```sh
mkdir /tmp/verify-release && cd /tmp/verify-release
go mod init verify
go get github.com/mycargus/riteway-golang@vX.Y.Z
```

---

## Yanking a broken release

Go has no "unpublish". If a release is broken:

1. Tag a patch release immediately with the fix.
2. To discourage use of the broken version, add a `retract` directive to `go.mod`:

   ```
   retract vX.Y.Z // reason: describes the problem
   ```

3. Release the new version. `go get` will warn users on the retracted version to
   upgrade.

---

## What is intentionally blocked

AI agents (Claude Code) cannot trigger a release. The following are blocked by a
PreToolUse hook in `.claude/hooks/block-publish.sh` (see ADR-011):

- `make release` / `bash scripts/release.sh`
- `git push origin v*` (version tags)
- `gh release create`

Human terminal usage is unaffected — the hook only runs inside Claude Code.
