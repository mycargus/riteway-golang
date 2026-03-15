# Testing Patterns for riteway-golang

## Package and imports

All tests are in the external test package:

```go
package riteway_test

import (
    "strings"
    "testing"

    "github.com/google/go-cmp/cmp"
    "github.com/google/go-cmp/cmp/cmpopts"
    riteway "github.com/mycargus/riteway-golang"
)
```

## fakeT spy

`fake_t_test.go` provides a `fakeT` struct that implements `testing.TB`. Use it to test failure paths without affecting the real test runner.

```go
ft := &fakeT{}
riteway.Assert(ft, riteway.Case[int]{...})
// Now inspect ft.errors ([]string) and ft.failed (bool)
```

Key: `fakeT.Fatalf` and `fakeT.FailNow` call `runtime.Goexit()`, so tests that exercise fatal paths (including Require) must run in a separate goroutine:

```go
ft := &fakeT{}
done := make(chan struct{})
go func() {
    defer close(done)
    riteway.Require(ft, riteway.Case[int]{...})
}()
<-done
// Now inspect ft.errors and ft.failed
```

## Testing panic paths

Use `riteway.Try` to catch panics:

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

## Assertion style

Every assertion uses `riteway.Assert` or `riteway.Require` (dogfooding). Given and Should must be non-empty, non-whitespace, and descriptive:

```go
riteway.Assert(t, riteway.Case[bool]{
    Given:    "nil map vs empty map",
    Should:   "detect them as unequal",
    Actual:   ft.failed,
    Expected: true,
})
```

## Test naming

Test names follow `TestFunction_Scenario`:
- `TestAssert_NilSlicesAreEqual`
- `TestMatchRegexp_ZeroMatchPattern_Panics`
- `TestTry_PanicNonError_ZeroValueResult`
- `TestRequire_IsFatal`

## File placement

- Assert, Require, Try tests → `riteway_test.go`
- Match, MatchRegexp tests → `match_test.go`
