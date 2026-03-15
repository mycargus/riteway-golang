# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`riteway-golang` is a Go port of the [paralleldrive/riteway](https://github.com/paralleldrive/riteway) JavaScript testing library. It enforces a structured testing philosophy where every assertion answers five questions: unit under test, expected behavior, actual output, expected output, and how to reproduce failure.

**Module:** `github.com/mycargus/riteway-golang`
**Go version:** 1.21+ (required for unambiguous `panic(nil)` via `*runtime.PanicNilError`)
**Direct dependency:** `github.com/google/go-cmp v0.7.0`

## Commands

```bash
make test          # go test -count=1 -race ./...  (no cache, race detector)
make fmt           # gofmt -l . (fails if unformatted files exist)
make vet           # go vet ./...
make check         # fmt + vet + test + test-scripts (full CI suite)
make test-scripts  # bash scripts/test-release.sh
```

To run a single test:
```bash
go test -run TestName ./...
```

## Architecture

The library is intentionally minimal: two source files, one package.

- **`riteway.go`** — `Case[T]`, `Assert[T]`, `Require[T]`, `Try[T]`, shared `doAssert` helper
- **`match.go`** — `Match` (literal substring) and `MatchRegexp` (regex)

### Key Design Decisions (see `docs/decisions/` for full ADRs)

- **`Case[T any]`** is generic: `Actual` and `Expected` must share the same type at compile time (ADR-002)
- **`Assert` vs `Require`**: both delegate to `doAssert(name, errorf)` — only the error function differs (`t.Errorf` vs `t.Fatalf`)
- **`Try[T]`** catches panics but re-propagates `runtime.Goexit` using a `completed` boolean flag (ADR-007); also rejects `fn == nil`
- **`Match("anything", "")`** returns `""` — empty substring is treated as not-found to avoid ambiguity (ADR-006)
- **`MatchRegexp`** panics on zero-match patterns like `x*` or `.*` — intentional ambiguity guard
- **`go-cmp` diff format**: `(-want +got)` where `want = Expected`, `got = Actual`
- `Given` and `Should` fields are validated as non-empty, non-whitespace before comparison runs (ADR-008)

### Test Infrastructure

- **`fake_t_test.go`** — `fakeT` spy type implementing `testing.TB` for testing failure paths without contaminating the parent test
- **`example_test.go`** — runnable examples (`ExampleMatch`, `ExampleMatchRegexp`, `ExampleTry`) and `BenchmarkAssert`

## Publishing Guard

`make release` and related publish commands are blocked by a PreToolUse hook (`.claude/hooks/block-publish.sh`). Releases must be triggered from a human terminal. Do not attempt to automate releases.

## Skills

Two project-specific skills are available:

- `/explore` — runs exploratory edge-case tests against the public API, adds missing test cases, and reports surprises
- `/review` — systematic walk-through of every public function checking error messages, go-cmp integration, docs accuracy, and real-world usage scenarios
