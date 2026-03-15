# ADR-009: `fakeT` spy for testing failure paths

**Status:** Accepted

## Context

Some `Assert[T]` behaviors can only be verified by observing what gets recorded on the `testing.TB` value — specifically, the error messages produced for validation failures (empty `Given`, empty `Should`, etc.) and mismatched value diffs.

A real `*testing.T` cannot be used for this: there is no way to read back the messages it recorded. Passing a real `t` to `Assert` with invalid inputs would also mark the enclosing test as failed, conflating "the test under test failed correctly" with "this test has a bug".

## Decision

Introduce a `fakeT` spy in `fake_t_test.go` (test-only, unexported) that implements `testing.TB` by embedding the interface and overriding only the methods under observation:

- `Helper()` — no-op (avoids nil pointer dereference from the embedded interface)
- `Error(args ...any)` — appends formatted message to `errors []string`, sets `failed = true`
- `Errorf(format string, args ...any)` — same with format string
- `FailNow()` — sets `failed = true`, calls `runtime.Goexit()` to simulate `t.FailNow()` behavior

Embedding `testing.TB` (the interface, not a concrete type) provides forward compatibility: if Go adds new methods to `testing.TB` in future versions, `fakeT` will still satisfy the interface via the embedded field without any changes.

All failure-path tests are fully dogfooded: `fakeT` is passed to `Assert`, and then a real `Assert(t, ...)` call verifies what `fakeT` recorded.

## Consequences

- Failure messages are verified exactly (e.g., `"riteway.Assert: Given must not be empty"`), not just approximately.
- Validation tests do not pollute the parent test's failure state.
- `fakeT` is scoped to `_test.go` files and has no impact on the public API.
- The `TestTry_RuntimeGoexit` test uses `fakeT.FailNow()` in a goroutine to verify `Goexit` propagation without causing the outer test to fail.
