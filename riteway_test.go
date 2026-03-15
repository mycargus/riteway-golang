package riteway_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	riteway "github.com/mycargus/riteway-golang"
)

// configWithSecret is a struct with an unexported field used to test cmp.Option support.
type configWithSecret struct {
	Port   int
	secret string
}

func TestAssert_HappyPath(t *testing.T) {
	riteway.Assert(t, riteway.Case[int]{
		Given:    "two equal integers",
		Should:   "pass without error",
		Actual:   42,
		Expected: 42,
	})
}

func TestAssert_ZeroValues(t *testing.T) {
	riteway.Assert(t, riteway.Case[int]{
		Given:    "zero value of int",
		Should:   "match another zero value",
		Actual:   0,
		Expected: 0,
	})
}

func TestAssert_EmptyGiven(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "", Should: "something", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "empty Given field",
		Should:   "record a validation error for Given",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Given must not be empty",
	})
}

func TestAssert_WhitespaceOnlyGiven(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "  ", Should: "something", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "whitespace-only Given field",
		Should:   "record a validation error for Given",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Given must not be empty",
	})
}

func TestAssert_EmptyShould(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "valid given", Should: "", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "empty Should field with valid Given",
		Should:   "record a validation error for Should",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Should must not be empty",
	})
}

func TestAssert_WhitespaceOnlyShould(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "valid given", Should: "  ", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "whitespace-only Should field",
		Should:   "record a validation error for Should",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Should must not be empty",
	})
}

func TestAssert_ZeroValueCase_FirstErrorIsGiven(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{})
	riteway.Assert(t, riteway.Case[int]{
		Given:    "zero-value Case[int]{}",
		Should:   "record exactly one error (Given, not Should, due to early return)",
		Actual:   len(ft.errors),
		Expected: 1,
	})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "zero-value Case[int]{}",
		Should:   "record the Given validation error first",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Given must not be empty",
	})
}

func TestAssert_MismatchedValues_ProducesDiff(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "mismatched values", Should: "produce a diff", Actual: 1, Expected: 2})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "mismatched Actual and Expected",
		Should:   "include (-want +got) annotation in error message",
		Actual:   strings.Contains(ft.errors[0], "(-want +got)"),
		Expected: true,
	})
}

func TestAssert_UnexportedFields(t *testing.T) {
	a := configWithSecret{Port: 8080, secret: "a"}
	b := configWithSecret{Port: 8080, secret: "b"}
	riteway.Assert(t, riteway.Case[configWithSecret]{
		Given:    "two structs with matching exported fields and differing unexported fields",
		Should:   "pass when unexported fields are ignored via cmpopts.IgnoreUnexported",
		Actual:   a,
		Expected: b,
	}, cmpopts.IgnoreUnexported(configWithSecret{}))
}

func TestTry_NoPanic(t *testing.T) {
	result, err := riteway.Try(func() int { return 42 })
	riteway.Assert(t, riteway.Case[int]{
		Given:    "a function that does not panic",
		Should:   "return the result",
		Actual:   result,
		Expected: 42,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a function that does not panic",
		Should:   "return nil error",
		Actual:   err == nil,
		Expected: true,
	})
}

func TestTry_PanicString(t *testing.T) {
	_, err := riteway.Try(func() int { panic("boom") })
	riteway.Assert(t, riteway.Case[string]{
		Given:    "a function that panics with a string",
		Should:   "return the panic message as an error",
		Actual:   err.Error(),
		Expected: "boom",
	})
}

func TestTry_PanicError(t *testing.T) {
	sentinel := errors.New("sentinel error")
	_, err := riteway.Try(func() int { panic(sentinel) })
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a function that panics with an error value",
		Should:   "return that exact error via errors.Is",
		Actual:   errors.Is(err, sentinel),
		Expected: true,
	})
}

func TestTry_PanicInt(t *testing.T) {
	_, err := riteway.Try(func() int { panic(42) })
	riteway.Assert(t, riteway.Case[string]{
		Given:    "a function that panics with an int",
		Should:   "return an error with message 'panic: 42'",
		Actual:   err.Error(),
		Expected: "panic: 42",
	})
}

func TestTry_PanicNil(t *testing.T) {
	_, err := riteway.Try(func() int { panic(nil) })
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a function that calls panic(nil) on Go 1.21+ (returns *runtime.PanicNilError)",
		Should:   "return a non-nil error",
		Actual:   err != nil,
		Expected: true,
	})
}

func TestTry_RuntimeGoexit(t *testing.T) {
	ft := &fakeT{}
	postTryReached := false
	done := make(chan struct{})
	go func() {
		defer close(done)
		riteway.Try(func() int { //nolint:errcheck
			ft.FailNow() // calls runtime.Goexit on this goroutine
			return 0
		})
		postTryReached = true // must NOT be reached
	}()
	<-done
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a function calling FailNow inside Try",
		Should:   "not execute code after Try (Goexit propagated to goroutine exit)",
		Actual:   postTryReached,
		Expected: false,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a function calling FailNow inside Try",
		Should:   "mark the caller as failed",
		Actual:   ft.failed,
		Expected: true,
	})
}
