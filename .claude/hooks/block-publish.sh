#!/usr/bin/env bash

# PreToolUse hook: block automated publishing (ADR 011)
# Reads JSON from stdin, pattern-matches the "command" field directly.
# Exit 2 = block, Exit 0 = allow.

input=$(cat)

case "$input" in
  *'"command":"bash scripts/release.sh'*|\
  *'"command":"make release'*|\
  *'"command":"git push origin v'*|\
  *'"command":"gh release create'*)
    echo "HOOK_BLOCKED: Publishing must be done manually. See ADR 011." >&2
    exit 2
    ;;
  *'"command":"git push origin main'*|\
  *'"command":"git push -u origin main'*)
    echo "HOOK_BLOCKED: Direct pushes to main are not allowed. Open a pull request instead." >&2
    exit 2
    ;;
esac

exit 0
