# ADR-003: Accept `testing.TB` instead of `*testing.T`

**Status:** Accepted

## Context

`Assert[T]` needs a way to report test failures. The obvious parameter type is `*testing.T`, but Go's testing package also provides `*testing.B` (benchmarks) and `*testing.F` (fuzz tests), both of which implement the `testing.TB` interface.

## Decision

`Assert[T]` accepts `testing.TB` rather than `*testing.T`.

## Consequences

- `Assert` can be called from benchmarks and fuzz tests without any adapter.
- The `fakeT` spy used in failure-path tests implements `testing.TB`, making it forward-compatible with new methods added to the interface in future Go versions (because `fakeT` embeds `testing.TB` and only overrides the specific methods it needs to observe).
- No behavioral difference for the common `*testing.T` case.
