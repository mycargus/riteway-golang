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
// context and a human-readable diff.
//
// Optional cmp.Options can configure comparison behavior (e.g.,
// cmpopts.IgnoreUnexported for structs with unexported fields).
//
// Assert validates that Given and Should are non-empty, non-whitespace
// strings. Whitespace-only values are rejected because they produce
// meaningless failure messages.
func Assert[T any](t testing.TB, c Case[T], opts ...cmp.Option) {
	t.Helper()

	if strings.TrimSpace(c.Given) == "" {
		t.Error("riteway.Assert: Given must not be empty")
		return
	}
	if strings.TrimSpace(c.Should) == "" {
		t.Error("riteway.Assert: Should must not be empty")
		return
	}

	if !cmp.Equal(c.Actual, c.Expected, opts...) {
		t.Errorf("Given %s: should %s (-want +got):\n%s",
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
//   - all other types become errors via fmt.Errorf("panic: %v", r)
//
// Try does not swallow runtime.Goexit. If fn calls t.FailNow or t.Fatal,
// Try detects the Goexit signal and re-propagates it, so the calling
// goroutine terminates rather than receiving a nil error.
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
			if r != nil {
				switch v := r.(type) {
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = fmt.Errorf("panic: %v", v)
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
