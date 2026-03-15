# ADR-007: `Try[T]` propagates `runtime.Goexit`

**Status:** Accepted

## Context

`Try[T]` uses a `defer`/`recover()` block to catch panics. However, `runtime.Goexit` (called internally by `t.FailNow()`, `t.Fatal()`, etc.) also causes deferred functions to run and `recover()` to return `nil`.

Without special handling, a `Try`-wrapped function that calls `t.FailNow()` would:
1. Trigger `runtime.Goexit`
2. Run the deferred recovery block in `Try`
3. See `recover()` return `nil` and `completed == false`
4. Return normally with a zero value and `nil` error

This silently swallows the `Goexit` signal. The test that called `t.FailNow()` would appear to succeed (returning `nil` error) instead of terminating the goroutine.

## Decision

`Try[T]` uses a `completed` boolean flag set to `true` only after `fn()` returns normally. In the deferred recovery block:

- If `recover()` returns a non-nil value: a panic occurred — convert it to an error.
- If `recover()` returns `nil` and `completed` is `false`: `Goexit` was called — call `runtime.Goexit()` again to re-propagate it.
- If `recover()` returns `nil` and `completed` is `true`: `fn` completed normally — do nothing.

## Consequences

- `t.FailNow()` and `t.Fatal()` inside a `Try`-wrapped function behave correctly: the goroutine terminates rather than returning a nil error.
- Code after the `Try(...)` call is not reached when `Goexit` is propagated — this is the invariant verified in `TestTry_RuntimeGoexit`.
- The `panic(nil)` case (which also causes `recover()` to return a `*runtime.PanicNilError` on Go 1.21+) is handled correctly: `*runtime.PanicNilError` is a non-nil value and takes the panic branch, not the Goexit branch.
