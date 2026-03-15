# Riteway Golang — Library Map

## Source Files

| File | Responsibility |
|------|---------------|
| `riteway.go` | `Assert[T]`, `Require[T]`, `Try[T]`; core assertion and panic recovery |
| `match.go` | `Match`, `MatchRegexp` — text search helpers |

## Test Files

| File | What it covers |
|------|---------------|
| `riteway_test.go` | Assert validation, value comparison, Require fatality, Try panic recovery, Goexit propagation |
| `match_test.go` | Match and MatchRegexp edge cases including Unicode, empty inputs, invalid regex, zero-match panic |
| `example_test.go` | Runnable usage examples (ExampleMatch, ExampleMatchRegexp, ExampleTry) and BenchmarkAssert |
| `fake_t_test.go` | `fakeT` spy implementing testing.TB — used to test failure paths without affecting parent test |

## Public API

```go
type Case[T any] struct {
    Given    string // describes the input or precondition
    Should   string // describes the expected behavior
    Actual   T      // the computed value
    Expected T      // the value we expect
}

func Assert[T any](t testing.TB, c Case[T], opts ...cmp.Option)
func Require[T any](t testing.TB, c Case[T], opts ...cmp.Option)
func Try[T any](fn func() T) (result T, err error)
func Match(text, substring string) string
func MatchRegexp(text, pattern string) string
```

## Design Decisions

See `docs/decisions/` — one ADR per file. Key entries:

| ADR | Decision |
|-----|---------|
| 001 | Go 1.21 minimum — needed for unambiguous `panic(nil)` via `*runtime.PanicNilError` |
| 002 | Generic `Case[T any]` — compiler enforces Actual/Expected same type |
| 003 | Accept `testing.TB` interface — works with `*testing.T`, `*testing.B`, `*testing.F` |
| 004 | Use go-cmp for deep equality — human-readable diffs, customizable via `cmp.Option` |
| 005 | Separate `Match` and `MatchRegexp` — idiomatic Go, no magic string prefixes |
| 006 | Empty substring returns `""` — distinguishes "not found" from ambiguous empty match |
| 007 | `Try[T]` propagates `runtime.Goexit` — uses `completed` flag to re-propagate termination |
| 008 | Validate `Given`/`Should` as non-empty, non-whitespace — rejects meaningless labels |
| 009 | `fakeT` spy for testing failure paths — embeds interface for forward compatibility |
| 010 | Omits JavaScript APIs — `CountKeys` (Go has `len()`), `Describe` factory (footgun) |

## Known Non-Obvious Behaviors

- `Assert` is non-fatal (`t.Errorf`); `Require` is fatal (`t.Fatalf`) — both share `doAssert` internally
- `Try` does not catch `runtime.Goexit` — uses a `completed` boolean to distinguish from `panic(nil)`, then re-calls `runtime.Goexit()` to re-propagate
- `panic(nil)` returns a `*runtime.PanicNilError` (non-nil) on Go 1.21+, so it is caught and returned as an error — not silently ignored
- `panic(error)` is returned as-is — `errors.Is` chains work correctly
- `panic(string)` is converted via `errors.New` — original string is preserved as `err.Error()`
- `panic(other)` is formatted as `"panic(%T): %v"` — type name is included
- `Assert` validates `Given` before `Should` — if both are empty, only the `Given` error fires (early return)
- `Assert`/`Require` return nothing; failure is reported via `t.Errorf`/`t.Fatalf`
- Validation error messages include the function name (`riteway.Assert:` vs `riteway.Require:`) and the bad value
- `go-cmp` panics when comparing structs with unexported fields unless `cmpopts.IgnoreUnexported` or `cmp.AllowUnexported` is passed via `opts`
- `Match("anything", "")` returns `""` — empty substring is treated as not found, even though `""` is technically in every string
- `Match` is case-sensitive
- `MatchRegexp` panics on invalid regex AND on patterns that can match the empty string (e.g., `x*`, `.*`) — the zero-match result would be indistinguishable from "not found"
- `MatchRegexp` returns only the first match (`FindString`) — no multi-match variant exists
- go-cmp diff format is `(-want +got)` where want = `Expected`, got = `Actual`
