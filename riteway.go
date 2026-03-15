// Package riteway provides structured test assertions that enforce
// the five questions every unit test must answer: what is the unit
// under test, what should it do, what was the actual output, what
// was the expected output, and how do you reproduce the failure.
//
// Assert is the core function. It uses google/go-cmp for deep
// equality comparison with human-readable diffs.
package riteway

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Case represents a single test assertion answering the five questions
// every test must answer: what is the unit under test, what should it do,
// what was the actual output, what was the expected output, and how do
// you reproduce the failure.
type Case[T any] struct {
	Given    string // describes the input or precondition
	Should   string // describes the expected behavior
	Actual   T      // the computed value
	Expected T      // the value we expect
}

// Assert compares Actual and Expected in the given Case using deep
// equality via go-cmp. On mismatch, it reports a failure with structured
// context and a human-readable diff. Assert is non-fatal: it calls
// t.Errorf (not t.Fatalf), so the test continues after a failed assertion.
//
// To stop the test immediately on failure, use Require instead.
//
// Optional cmp.Options can configure comparison behavior (e.g.,
// cmpopts.IgnoreUnexported for structs with unexported fields).
//
// Assert validates that Given and Should are non-empty, non-whitespace
// strings. Whitespace-only values are rejected because they produce
// meaningless failure messages.
//
// If go-cmp panics (e.g. due to unexported fields with no opts),
// Assert recovers the panic and reports it as a test failure.
func Assert[T any](t testing.TB, c Case[T], opts ...cmp.Option) {
	t.Helper()
	doAssert(t, c, "riteway.Assert", t.Errorf, opts...)
}

// Require compares Actual and Expected in the given Case using deep
// equality via go-cmp. On mismatch, it reports a failure with structured
// context and a human-readable diff, then stops the test immediately.
//
// Require is the fatal counterpart to Assert: it calls t.Fatalf (not
// t.Errorf), so the test halts on the first failed assertion.
//
// Optional cmp.Options and validation behavior are identical to Assert.
func Require[T any](t testing.TB, c Case[T], opts ...cmp.Option) {
	t.Helper()
	doAssert(t, c, "riteway.Require", t.Fatalf, opts...)
}

// doAssert is the shared implementation for Assert and Require.
// name is used in error messages; errorf is t.Errorf or t.Fatalf.
func doAssert[T any](t testing.TB, c Case[T], name string, errorf func(string, ...any), opts ...cmp.Option) {
	t.Helper()

	if strings.TrimSpace(c.Given) == "" {
		errorf("%s: Given must not be empty (got %q)", name, c.Given)
		return
	}
	if strings.TrimSpace(c.Should) == "" {
		errorf("%s: Should must not be empty (got %q)", name, c.Should)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errorf("%s: comparison panicked: %v", name, r)
		}
	}()

	if !cmp.Equal(c.Actual, c.Expected, opts...) {
		errorf("Given %s: should %s (-want +got):\n%s",
			c.Given, c.Should,
			cmp.Diff(c.Expected, c.Actual, opts...))
	}
}

// Try calls fn and recovers from any panic, returning it as an error.
// If fn succeeds, Try returns its result and a nil error.
//
// Panic values are converted to errors as follows:
//   - error values are returned as-is
//   - string values become errors via errors.New
//   - all other types become errors via fmt.Errorf("panic(%T): %v", r, r)
//
// Try does not swallow runtime.Goexit. If fn calls t.FailNow or t.Fatal,
// Try detects the Goexit signal and re-propagates it, so the calling
// goroutine terminates rather than receiving a nil error.
//
// On panic, result is the zero value of T.
//
// Note: panic(nil) behavior differs across Go versions. On Go 1.21+,
// recover() returns a *runtime.PanicNilError; on earlier versions it
// returns nil (indistinguishable from no panic). This module requires
// Go 1.21 to avoid that ambiguity.
func Try[T any](fn func() T) (result T, err error) {
	completed := false
	defer func() {
		r := recover()
		if !completed {
			result = *new(T) // explicitly zero the result on any non-completion path
			if r != nil {
				switch v := r.(type) {
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = fmt.Errorf("panic(%T): %v", v, v)
				}
			} else {
				// recover() returned nil and fn didn't complete: Goexit was called.
				runtime.Goexit()
			}
		}
	}()
	result = fn()
	completed = true
	return
}
