# ADR-004: Use `go-cmp` for deep equality and diffs

**Status:** Accepted

## Context

Go's `reflect.DeepEqual` performs deep equality but provides no diff output and panics on structs with unexported fields in some configurations. A failed assertion needs to show *what* differed, not just *that* it differed.

Options considered:

1. **`reflect.DeepEqual` + manual formatting** — no diff, panics on unexported fields, requires significant extra code for useful output.
2. **`github.com/google/go-cmp`** — deep equality with human-readable diffs, extensible via `cmp.Option` (e.g., `cmpopts.IgnoreUnexported`, `cmp.AllowUnexported`, custom transformers), maintained by Google.
3. **`github.com/stretchr/testify`** — brings a large dependency with its own assertion style, which conflicts with riteway's philosophy of a minimal, focused API.

## Decision

Use `github.com/google/go-cmp` for both equality checks (`cmp.Equal`) and diff generation (`cmp.Diff`).

Pass `...cmp.Option` through `Assert[T]` to the caller, enabling full control over comparison behavior without any special-casing in riteway itself.

## Consequences

- Structs with unexported fields do not panic by default; callers choose how to handle them via `cmp.Option`.
- Diff output uses go-cmp's `(-want +got)` convention, which is the established Go ecosystem standard.
- Adds one external dependency (`go-cmp`). This is a widely-used, stable library with no transitive dependencies.
