# ADR-005: Separate `Match` and `MatchRegexp` functions

**Status:** Accepted

## Context

The original JavaScript riteway uses a single `match` function with a magic `"/"` prefix convention to distinguish literal substring matching from regular expression matching:

```js
match(text, "/pattern/")  // regexp
match(text, "substring")  // literal
```

This works in JavaScript but is not idiomatic Go. Go convention strongly favors distinct functions for distinct behaviors rather than encoding behavior in string values.

## Decision

Provide two separate functions: `Match(text, substring string) string` for literal substring search and `MatchRegexp(text, pattern string) string` for regular expression search.

## Consequences

- Behavior is explicit at the call site — no magic prefix to remember or mistype.
- Each function has a clear, focused signature and doc comment.
- `MatchRegexp` can document its panic behavior independently without complicating `Match`.
- Aligns with Go standard library conventions (e.g., `strings.Contains` vs `regexp.MatchString`).
