# ADR-010: APIs from JavaScript riteway not ported to Go

**Status:** Accepted

## Context

The JavaScript riteway library includes several helpers beyond the core `assert` function. Each was evaluated for inclusion in this Go port.

## Decision

The following are intentionally excluded:

### `CountKeys(obj)`

Returns `Object.keys(obj).length` in JavaScript — a workaround for the absence of a built-in `len` on objects. In Go, `len(m)` works natively on maps, slices, strings, arrays, and channels. No equivalent helper is needed.

### `Describe(label, fn)` factory / closure pattern

In JavaScript riteway, `Describe` creates a factory function that closes over a `given` label, reducing repetition in related assertions. In Go, the equivalent would close over a `testing.TB` value.

Capturing `testing.TB` in a closure is a footgun: if the closure is called from a goroutine or subtest after the original test has finished, the captured `t` is invalid. Go's test framework explicitly warns against this pattern. The idiomatic Go alternative is a plain function or a table-driven loop, both of which are already well-supported by `Case[T]`.

### React rendering utilities, CLI runner, TAP output, framework adapters

These are JavaScript/Node.js-specific and have no applicable Go counterpart.

## Consequences

- The Go API surface is minimal: `Assert`, `Case`, `Try`, `Match`, `MatchRegexp`.
- Users familiar with JavaScript riteway will find the Go port smaller but fully sufficient for the core philosophy.
- No footgun around captured `testing.TB` values.
