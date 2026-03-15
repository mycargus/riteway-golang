# ADR-008: Validate `Given` and `Should` as non-empty, non-whitespace

**Status:** Accepted

## Context

The riteway philosophy requires every test to answer five questions. `Given` and `Should` are the fields that communicate the test's intent in human-readable failure messages. If either field is empty or contains only whitespace, the failure message becomes:

```
Given : should  (-want +got):
```

This is meaningless and defeats the purpose of the library.

## Decision

`Assert[T]` validates that both `Given` and `Should` are non-empty after trimming whitespace (`strings.TrimSpace`). If either fails validation:

1. A specific error message is reported via `t.Error(...)` naming the offending field.
2. `Assert` returns immediately (early return after the first failure).

The early return means that if both fields are empty (e.g., a zero-value `Case[T]{}`), only the `Given` error is recorded. `Should` is not checked until `Given` passes. This makes failure messages deterministic and allows each validation to be independently tested.

## Consequences

- Callers are forced to provide meaningful field values at compile time (strings) but are caught at test runtime if they supply empty or whitespace-only strings.
- The `fakeT` spy in tests can verify each validation error message exactly.
- Whitespace-only strings (`"   "`) are rejected for the same reason as empty strings — they produce visually blank failure messages.
