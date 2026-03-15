#!/usr/bin/env bash

# PreToolUse hook: block automated publishing (ADR 012)
# Reads JSON from stdin, pattern-matches the "command" field directly.
# Exit 2 = block, Exit 0 = allow.

input=$(cat)

case "$input" in
  *'"command":"bash scripts/release.sh'*|\
  *'"command":"git push origin v'*|\
  *'"command":"gh release create'*)
    echo "HOOK_BLOCKED: Publishing must be done manually. See ADR 012." >&2
    exit 2
    ;;
esac

exit 0
