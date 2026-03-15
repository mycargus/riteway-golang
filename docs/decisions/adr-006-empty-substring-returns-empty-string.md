# ADR-006: `Match` with an empty substring returns `""`

**Status:** Accepted

## Context

`strings.Contains("anything", "")` returns `true` in Go — the empty string is technically a substring of every string. If `Match` followed this convention, `Match("anything", "")` would return `""`.

The return value `""` already means "not found" in `Match`'s API. Returning `""` when the empty string is "found" is therefore indistinguishable from "not found", producing a confusing and useless result.

## Decision

`Match` returns `""` when the `substring` argument is empty, regardless of the `text` value. This is documented explicitly in the function's doc comment.

## Consequences

- The return value of `Match` is unambiguous: non-empty means found, `""` means not found.
- Callers searching for the empty string get a clear, predictable result rather than a misleadingly successful one.
- This deviates from `strings.Contains` behavior, which is intentional and documented.

A similar note applies to `MatchRegexp`: patterns that can match the empty string (e.g., `x*`) may return `""` even on a successful match. The doc comment recommends anchored patterns (e.g., `x+`) to avoid this ambiguity.
