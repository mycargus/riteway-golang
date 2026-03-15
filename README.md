# riteway-golang

A Go port of the [paralleldrive/riteway](https://github.com/paralleldrive/riteway) testing philosophy. Every test must answer five questions:

1. **What is the unit under test?**
2. **What should it do?**
3. **What was the actual output?**
4. **What was the expected output?**
5. **How do you reproduce the failure?**

Enforcing these questions produces failure messages that are immediately actionable — no guessing, no archaeology.

## Install

```sh
go get github.com/mycargus/riteway-golang
```

## Import

The module path is `github.com/mycargus/riteway-golang` but the Go package name is `riteway`:

```go
import riteway "github.com/mycargus/riteway-golang"
// then use as: riteway.Assert(...)
```

## API

### `Case[T]`

```go
type Case[T any] struct {
    Given    string // describes the input or precondition
    Should   string // describes the expected behavior
    Actual   T      // the computed value
    Expected T      // the value we expect
}
```

### `Assert[T]`

```go
func Assert[T any](t testing.TB, c Case[T], opts ...cmp.Option)
```

Compares `Actual` and `Expected` using [go-cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp) deep equality. On mismatch, reports a failure with a structured message and human-readable diff.

- Validates that `Given` and `Should` are non-empty and non-whitespace.
- Accepts optional `cmp.Option` values for custom comparison (e.g., `cmpopts.IgnoreUnexported` to skip unexported fields, `cmp.AllowUnexported` to compare them).
- Works with `*testing.T`, `*testing.B`, and `*testing.F`.

### `Try[T]`

```go
func Try[T any](fn func() T) (result T, err error)
```

Calls `fn` and recovers from any panic, returning it as an error. Useful for asserting panic behavior in tests. Does **not** catch `runtime.Goexit` (i.e., `t.FailNow`/`t.Fatal` inside `Try` still terminate the subtest normally).

### `Match`

```go
func Match(text, substring string) string
```

Returns `substring` if found in `text`, otherwise `""`. An empty `substring` always returns `""` to avoid the ambiguous case where `Match("anything", "")` returns `""` and is indistinguishable from "not found".

### `MatchRegexp`

```go
func MatchRegexp(text, pattern string) string
```

Returns the first match of `pattern` in `text`, or `""` if not found. Panics if `pattern` is not a valid regular expression. Use `Try` to test for that panic.

## Usage

### Basic assertion

```go
func TestAdd(t *testing.T) {
    riteway.Assert(t, riteway.Case[int]{
        Given:    "no arguments",
        Should:   "return 0",
        Actual:   Add(),
        Expected: 0,
    })
}
```

### Table-driven tests

```go
func TestSquare(t *testing.T) {
    cases := []riteway.Case[int]{
        {Given: "zero",     Should: "return 0", Actual: Square(0), Expected: 0},
        {Given: "positive", Should: "return 4", Actual: Square(2), Expected: 4},
        {Given: "negative", Should: "return 9", Actual: Square(-3), Expected: 9},
    }
    for _, c := range cases {
        t.Run("Given "+c.Given, func(t *testing.T) {
            riteway.Assert(t, c)
        })
    }
}
```

### Structs with unexported fields

```go
import "github.com/google/go-cmp/cmp/cmpopts"

riteway.Assert(t, riteway.Case[Config]{
    Given:    "default settings",
    Should:   "use port 8080",
    Actual:   NewConfig(),
    Expected: Config{Port: 8080},
}, cmpopts.IgnoreUnexported(Config{}))
```

### Panic testing with Try

```go
_, err := riteway.Try(func() int { panic("boom") })
riteway.Assert(t, riteway.Case[string]{
    Given:    "a panicking function",
    Should:   "return the panic message as an error",
    Actual:   err.Error(),
    Expected: "boom",
})
```

### Text matching

```go
riteway.Assert(t, riteway.Case[string]{
    Given:    "rendered HTML with a title",
    Should:   "contain the page title",
    Actual:   riteway.Match(html, "Welcome"),
    Expected: "Welcome",
})
```

### MatchRegexp panic testing

```go
_, err := riteway.Try(func() string {
    return riteway.MatchRegexp("text", "[invalid")
})
riteway.Assert(t, riteway.Case[bool]{
    Given:    "an invalid regexp pattern",
    Should:   "panic",
    Actual:   err != nil,
    Expected: true,
})
```

## Failure output

When a test fails, riteway produces:

```
--- FAIL: TestSquare/Given_negative (0.00s)
    riteway_test.go:42: Given negative: should return 9 (-want +got):
          int(
        -     9,
        +     10,
          )
```

## Requirements

- Go 1.21+

## Attribution

This library is a Go port of [paralleldrive/riteway](https://github.com/paralleldrive/riteway), originally created by [Eric Elliott](https://github.com/ericelliott). The five-question testing philosophy, API design, and naming conventions are derived from that work.

## License

MIT
