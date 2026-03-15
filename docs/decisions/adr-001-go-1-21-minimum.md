# ADR-001: Go 1.21 as minimum version

**Status:** Accepted

## Context

The library uses generics (`Case[T]`, `Assert[T]`, `Try[T]`) which require Go 1.18+. It also needs to handle `panic(nil)` unambiguously inside `Try[T]`.

Prior to Go 1.21, calling `recover()` after `panic(nil)` returned `nil`, which was indistinguishable from no panic occurring at all. Go 1.21 introduced `*runtime.PanicNilError` as the recovery value for `panic(nil)`, eliminating the ambiguity.

## Decision

Require Go 1.21 as the minimum version.

## Consequences

- `Try[T]` can reliably distinguish `panic(nil)` from a normal return.
- Users on Go versions older than 1.21 cannot use this library.
