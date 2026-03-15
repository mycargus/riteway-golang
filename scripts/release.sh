#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/release.sh v0.1.0

VERSION="${1:-}"

# ── Validate argument ──────────────────────────────────────────────────────────

if [[ -z "$VERSION" ]]; then
  echo "usage: $0 <version>"
  echo "  example: $0 v0.1.0"
  exit 1
fi

if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "error: version must match vMAJOR.MINOR.PATCH (e.g. v0.1.0)"
  exit 1
fi

# ── Must be run from repo root ─────────────────────────────────────────────────

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

# ── Extract release notes from CHANGELOG.md ────────────────────────────────────

# Strip the leading 'v' to match CHANGELOG headings (e.g. [0.1.0])
BARE_VERSION="${VERSION#v}"

if [[ ! -f CHANGELOG.md ]]; then
  echo "error: CHANGELOG.md not found"
  exit 1
fi

RELEASE_NOTES="$(awk \
  "/^## \[${BARE_VERSION}\]/{found=1; next} found && /^## \[/{exit} found{print}" \
  CHANGELOG.md)"

if [[ -z "$RELEASE_NOTES" ]]; then
  echo "error: no entry for ${VERSION} found in CHANGELOG.md"
  echo "  add a '## [${BARE_VERSION}]' section before releasing"
  exit 1
fi

# ── Clean working tree ─────────────────────────────────────────────────────────

if [[ -n "$(git status --porcelain)" ]]; then
  echo "error: working tree is not clean — commit or stash changes first"
  git status --short
  exit 1
fi

# ── Tag must not already exist ─────────────────────────────────────────────────

if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "error: tag $VERSION already exists"
  exit 1
fi

# ── Pre-release checks ─────────────────────────────────────────────────────────

echo "==> running checks..."
make check

# ── go.sum must be tidy ────────────────────────────────────────────────────────

echo "==> verifying go.sum is tidy..."
go mod tidy
if [[ -n "$(git status --porcelain go.sum go.mod)" ]]; then
  echo "error: go mod tidy changed go.sum or go.mod — commit those changes first"
  git diff go.sum go.mod
  exit 1
fi

# ── Confirm ────────────────────────────────────────────────────────────────────

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
COMMIT="$(git rev-parse --short HEAD)"

echo ""
echo "  module : $(go list -m)"
echo "  version: $VERSION"
echo "  commit : $COMMIT ($CURRENT_BRANCH)"
echo ""
echo "release notes:"
echo "$RELEASE_NOTES" | sed 's/^/  /'
echo ""
read -r -p "tag, push, and create GitHub release? [y/N] " CONFIRM || CONFIRM=""
if [[ "$CONFIRM" != "y" && "$CONFIRM" != "Y" ]]; then
  echo "aborted"
  exit 0
fi

# ── Tag and push ───────────────────────────────────────────────────────────────

git tag "$VERSION"
git push origin "$VERSION"

echo ""
echo "tagged and pushed $VERSION — GitHub Actions will create the release"
echo "  https://github.com/mycargus/riteway-golang/releases/tag/$VERSION"
echo "  https://pkg.go.dev/$(go list -m)@$VERSION"
