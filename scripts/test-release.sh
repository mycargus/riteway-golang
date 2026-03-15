#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
RELEASE_SCRIPT="$SCRIPT_DIR/release.sh"
ORIG_DIR="$PWD"

PASS=0
FAIL=0

# ── Harness ────────────────────────────────────────────────────────────────────

pass() { echo "    PASS: $1"; PASS=$(( PASS + 1 )); }
fail() { echo "    FAIL: $1"; [[ -n "${2:-}" ]] && echo "          $2"; FAIL=$(( FAIL + 1 )); }

assert_exit() {
  local expected=$1 actual=$2 desc=$3
  if [[ "$actual" -eq "$expected" ]]; then
    pass "$desc"
  else
    fail "$desc" "expected exit $expected, got $actual"
  fi
}

assert_contains() {
  local needle=$1 haystack=$2 desc=$3
  if printf '%s' "$haystack" | grep -qF "$needle"; then
    pass "$desc"
  else
    fail "$desc" "output did not contain: $needle"
  fi
}

# Run the release script in the current directory.
# stdin is /dev/null so the confirm prompt never blocks.
run() {
  set +e
  output=$("$RELEASE_SCRIPT" "$@" 2>&1 </dev/null)
  exit_code=$?
  set -e
}

# ── Repo helpers ───────────────────────────────────────────────────────────────

make_repo() {
  local dir
  dir=$(mktemp -d)
  git -C "$dir" init -q
  git -C "$dir" config user.email "test@example.com"
  git -C "$dir" config user.name "Test"

  # Minimal go.mod — no deps, so go mod tidy is a no-op
  printf 'module example.com/test\n\ngo 1.21\n' > "$dir/go.mod"

  # make check is a no-op to avoid running real Go tests
  printf 'check:\n\t@true\n' > "$dir/Makefile"

  echo "$dir"
}

add_changelog() {
  local dir=$1 version=${2:-0.1.0}
  cat > "$dir/CHANGELOG.md" <<EOF
# Changelog

## [Unreleased]

## [${version}] - 2026-03-14

### Added
- Initial release
EOF
}

commit_all() {
  git -C "$1" add .
  git -C "$1" commit -q -m "initial"
}

# ── Stub gh to prevent accidental real GitHub calls ────────────────────────────

STUBS=$(mktemp -d)
printf '#!/usr/bin/env bash\necho "gh stub: $*"\nexit 0\n' > "$STUBS/gh"
chmod +x "$STUBS/gh"
export PATH="$STUBS:$PATH"

# ── 1. Version format validation ───────────────────────────────────────────────

echo ""
echo "1. version format validation"

dir=$(make_repo); add_changelog "$dir"; commit_all "$dir"; cd "$dir"

run;              assert_exit 1 $exit_code "no argument: exits 1"
                  assert_contains "usage:" "$output" "no argument: prints usage"

run "0.1.0";      assert_exit 1 $exit_code "missing v prefix: exits 1"
                  assert_contains "vMAJOR.MINOR.PATCH" "$output" "missing v prefix: prints format hint"

run "v1.0";       assert_exit 1 $exit_code "missing patch segment: exits 1"

run "v1.0.0-rc1"; assert_exit 1 $exit_code "pre-release suffix: exits 1"

run "v1.0.0a";    assert_exit 1 $exit_code "non-numeric patch: exits 1"

run "v0.1.0";     assert_exit 0 $exit_code "valid version: exits 0 (aborted at confirm)"
                  assert_contains "aborted" "$output" "valid version: reaches confirm prompt"

cd "$ORIG_DIR"

# ── 2. CHANGELOG extraction ────────────────────────────────────────────────────

echo ""
echo "2. CHANGELOG extraction"

dir=$(make_repo); commit_all "$dir"; cd "$dir"

run "v0.1.0";     assert_exit 1 $exit_code "no CHANGELOG: exits 1"
                  assert_contains "CHANGELOG.md not found" "$output" "no CHANGELOG: prints message"

# CHANGELOG exists but version section is for a different release
add_changelog "$dir" "0.2.0"
git -C "$dir" add . && git -C "$dir" commit -q -m "add changelog"

run "v0.1.0";     assert_exit 1 $exit_code "wrong version in CHANGELOG: exits 1"
                  assert_contains "no entry for v0.1.0" "$output" "wrong version in CHANGELOG: correct message"

# Correct version — notes should appear in the confirm summary
add_changelog "$dir" "0.1.0"
git -C "$dir" add . && git -C "$dir" commit -q -m "add 0.1.0 entry"

run "v0.1.0";     assert_exit 0 $exit_code "matching CHANGELOG entry: reaches confirm"
                  assert_contains "Initial release" "$output" "matching CHANGELOG entry: shows notes in confirm"

cd "$ORIG_DIR"

# ── 3. Clean working tree ──────────────────────────────────────────────────────

echo ""
echo "3. clean working tree"

dir=$(make_repo); add_changelog "$dir"; commit_all "$dir"; cd "$dir"

echo "dirty" > untracked.txt
run "v0.1.0";     assert_exit 1 $exit_code "untracked file: exits 1"
                  assert_contains "working tree is not clean" "$output" "untracked file: prints message"
rm untracked.txt

echo "dirty" >> go.mod
run "v0.1.0";     assert_exit 1 $exit_code "modified tracked file: exits 1"
                  assert_contains "working tree is not clean" "$output" "modified tracked file: prints message"
git checkout -- go.mod

cd "$ORIG_DIR"

# ── 4. Tag already exists ──────────────────────────────────────────────────────

echo ""
echo "4. tag already exists"

dir=$(make_repo); add_changelog "$dir"; commit_all "$dir"; cd "$dir"

git tag v0.1.0
run "v0.1.0";     assert_exit 1 $exit_code "existing tag: exits 1"
                  assert_contains "tag v0.1.0 already exists" "$output" "existing tag: prints message"

run "v0.2.0";     assert_exit 1 $exit_code "different version, no CHANGELOG entry: exits 1"
                  assert_contains "no entry for v0.2.0" "$output" "different version: CHANGELOG check fires"

cd "$ORIG_DIR"

# ── 5. go.sum tidy ────────────────────────────────────────────────────────────

echo ""
echo "5. go.sum tidy"

dir=$(make_repo); add_changelog "$dir"

# Commit a go.sum with a fake entry — go mod tidy will remove it
printf 'fake/module v1.0.0 h1:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa=\n' > "$dir/go.sum"
commit_all "$dir"
cd "$dir"

run "v0.1.0";     assert_exit 1 $exit_code "stale go.sum: exits 1"
                  assert_contains "go mod tidy changed" "$output" "stale go.sum: prints message"

cd "$ORIG_DIR"

# ── Results ───────────────────────────────────────────────────────────────────

echo ""
echo "results: $PASS passed, $FAIL failed"
echo ""

[[ $FAIL -eq 0 ]]
