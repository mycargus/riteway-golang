# ADR-011: No Publishing Automation

**Context:** Go module tags are permanent — once a tag is pushed and indexed by the module proxy (`sum.golang.org`), the version cannot be retracted without leaving a visible retraction notice in `go.mod`. AI agents running release commands would remove the deliberate friction that makes accidental releases hard.

**Decision:** `scripts/release.sh` does not run inside Claude Code. The following are blocked by a PreToolUse hook in `.claude/hooks/block-publish.sh`:

- `bash scripts/release.sh`
- `make release`
- `git push origin v*` (version tags)
- `gh release create`

**Consequences:**
- No accidental releases from AI sessions or muscle memory.
- Explicit friction before every release is intentional.
- Human terminal usage is unaffected — the hook only fires inside Claude Code.
