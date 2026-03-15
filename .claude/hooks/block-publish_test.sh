#!/usr/bin/env bash

# Tests for block-publish.sh hook
# Run: bash .claude/hooks/block-publish_test.sh

HOOK=".claude/hooks/block-publish.sh"
pass=0
fail=0

assert() {
  local description="$1" expected="$2" actual="$3"
  if [ "$actual" = "$expected" ]; then
    echo "PASS: $description"
    pass=$((pass + 1))
  else
    echo "FAIL: $description (expected exit $expected, got exit $actual)"
    fail=$((fail + 1))
  fi
}

# Block cases
echo '{"tool_input":{"command":"bash scripts/release.sh v0.1.0"}}' | bash "$HOOK" 2>/dev/null
assert "blocks 'bash scripts/release.sh'" "2" "$?"

echo '{"tool_input":{"command":"git push origin v0.1.0"}}' | bash "$HOOK" 2>/dev/null
assert "blocks 'git push origin v<tag>'" "2" "$?"

echo '{"tool_input":{"command":"git push origin v1.2.3"}}' | bash "$HOOK" 2>/dev/null
assert "blocks 'git push origin v<semver>'" "2" "$?"

echo '{"tool_input":{"command":"gh release create v0.1.0 --title v0.1.0"}}' | bash "$HOOK" 2>/dev/null
assert "blocks 'gh release create'" "2" "$?"

# Allow cases
echo '{"tool_input":{"command":"go test ./..."}}' | bash "$HOOK" 2>/dev/null
assert "allows 'go test ./...'" "0" "$?"

echo '{"tool_input":{"command":"git commit -m blocks release"}}' | bash "$HOOK" 2>/dev/null
assert "allows commit message mentioning release" "0" "$?"

echo '{"tool_input":{"command":"git status"}}' | bash "$HOOK" 2>/dev/null
assert "allows 'git status'" "0" "$?"

echo '{"tool_input":{"command":"bash scripts/test-release.sh"}}' | bash "$HOOK" 2>/dev/null
assert "allows 'bash scripts/test-release.sh'" "0" "$?"

echo '{"tool_input":{"command":"git push origin main"}}' | bash "$HOOK" 2>/dev/null
assert "allows 'git push origin main'" "0" "$?"

echo '{"tool_input":{"command":"make check"}}' | bash "$HOOK" 2>/dev/null
assert "allows 'make check'" "0" "$?"

echo ""
echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ] || exit 1
