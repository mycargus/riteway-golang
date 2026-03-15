# ADR-002: Generic `Case[T]` struct for type-safe assertions

**Status:** Accepted

## Context

The original JavaScript riteway uses a plain object `{ given, should, actual, expected }`. Porting this to Go requires a decision about how to represent the assertion without sacrificing type safety.

Options considered:

1. **`interface{}`/`any` fields** — `Actual any` and `Expected any`. Simple, but allows comparing values of different types at runtime, producing confusing errors or silent type mismatches.
2. **Separate typed assertion functions** — `AssertInt`, `AssertString`, etc. Verbose, does not scale.
3. **Generic struct `Case[T any]`** — `Actual T` and `Expected T`. The compiler enforces that both sides have the same type.

## Decision

Use a generic struct `Case[T any]` where both `Actual` and `Expected` are of type `T`.

## Consequences

- Type mismatches between `Actual` and `Expected` are compile errors, not runtime panics.
- A single `Assert[T]` function handles all types without reflection or type assertions.
- Requires Go 1.18+ (already required for other reasons — see ADR-001).
- Users must sometimes provide the type parameter explicitly (e.g., `Case[int]{...}`) when it cannot be inferred, which is a minor ergonomic cost.
