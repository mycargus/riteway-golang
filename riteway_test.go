package riteway_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		Should:   "record a validation error for Given showing the bad value",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Given must not be empty (got \"\")",
	})
}

func TestAssert_WhitespaceOnlyGiven(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "  ", Should: "something", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "whitespace-only Given field",
		Should:   "record a validation error for Given showing the bad value",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Given must not be empty (got \"  \")",
	})
}

func TestAssert_EmptyShould(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "valid given", Should: "", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "empty Should field with valid Given",
		Should:   "record a validation error for Should showing the bad value",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Should must not be empty (got \"\")",
	})
}

func TestAssert_WhitespaceOnlyShould(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "valid given", Should: "  ", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "whitespace-only Should field",
		Should:   "record a validation error for Should showing the bad value",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Should must not be empty (got \"  \")",
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
		Should:   "record the Given validation error first showing the bad value",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Given must not be empty (got \"\")",
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
		Should:   "return an error with message 'panic(int): 42' preserving the type name",
		Actual:   err.Error(),
		Expected: "panic(int): 42",
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

func TestAssert_UnexportedFields_WithoutOpts_RecoversPanic(t *testing.T) {
	ft := &fakeT{}
	a := configWithSecret{Port: 8080, secret: "a"}
	b := configWithSecret{Port: 8080, secret: "b"}
	riteway.Assert(ft, riteway.Case[configWithSecret]{
		Given:    "struct with unexported field and no cmp.Option",
		Should:   "record an error instead of panicking",
		Actual:   a,
		Expected: b,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "Assert called with unexported field struct and no opts",
		Should:   "mark fakeT as failed without crashing the test binary",
		Actual:   ft.failed,
		Expected: true,
	})
	riteway.Assert(t, riteway.Case[int]{
		Given:    "Assert called with unexported field struct and no opts",
		Should:   "record exactly one error (Equal and Diff share one recover)",
		Actual:   len(ft.errors),
		Expected: 1,
	})
	if len(ft.errors) > 0 {
		riteway.Assert(t, riteway.Case[bool]{
			Given:    "Assert called with unexported field struct and no opts",
			Should:   "mention unexported fields in the error message",
			Actual:   strings.Contains(ft.errors[0], "unexported"),
			Expected: true,
		})
	}
}

func TestAssert_PanicFromCmpOption_NoFalseUnexportedGuidance(t *testing.T) {
	ft := &fakeT{}
	panicOpt := cmp.Comparer(func(x, y int) bool { panic("comparer panicked") })
	riteway.Assert(ft, riteway.Case[int]{
		Given:    "cmp.Option that panics with an unrelated message",
		Should:   "record an error without suggesting unexported-fields guidance",
		Actual:   1,
		Expected: 2,
	}, panicOpt)
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "cmp.Option that panics with unrelated message",
		Should:   "mark fakeT as failed",
		Actual:   ft.failed,
		Expected: true,
	})
	if len(ft.errors) > 0 {
		riteway.Assert(t, riteway.Case[bool]{
			Given:    "cmp.Option panic unrelated to unexported fields",
			Should:   "not mention unexported fields in the error message",
			Actual:   strings.Contains(ft.errors[0], "unexported"),
			Expected: false,
		})
	}
}

func TestTry_PanicNonError_ZeroValueResult(t *testing.T) {
	result, _ := riteway.Try(func() int { panic("oops") })
	riteway.Assert(t, riteway.Case[int]{
		Given:    "a function that panics",
		Should:   "return the zero value of T (0 for int, not a partial result)",
		Actual:   result,
		Expected: 0,
	})
}

func TestRequire_HappyPath(t *testing.T) {
	riteway.Require(t, riteway.Case[int]{
		Given:    "two equal integers",
		Should:   "pass without error",
		Actual:   42,
		Expected: 42,
	})
}

func TestRequire_IsFatal(t *testing.T) {
	ft := &fakeT{}
	postRequireReached := false
	done := make(chan struct{})
	go func() {
		defer close(done)
		riteway.Require(ft, riteway.Case[int]{
			Given:    "mismatched values",
			Should:   "stop the test immediately",
			Actual:   1,
			Expected: 2,
		})
		postRequireReached = true // must NOT be reached
	}()
	<-done
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a failed Require",
		Should:   "not execute code after Require",
		Actual:   postRequireReached,
		Expected: false,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a failed Require",
		Should:   "mark the test as failed",
		Actual:   ft.failed,
		Expected: true,
	})
}

func TestRequire_EmptyGiven_IsFatal(t *testing.T) {
	ft := &fakeT{}
	postRequireReached := false
	done := make(chan struct{})
	go func() {
		defer close(done)
		riteway.Require(ft, riteway.Case[int]{Given: "", Should: "something", Actual: 1, Expected: 1})
		postRequireReached = true // must NOT be reached
	}()
	<-done
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "Require with empty Given",
		Should:   "not execute code after Require",
		Actual:   postRequireReached,
		Expected: false,
	})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "Require with empty Given",
		Should:   "record a validation error naming riteway.Require",
		Actual:   ft.errors[0],
		Expected: "riteway.Require: Given must not be empty (got \"\")",
	})
}

// --- edge case tests ---

type address struct {
	Street string
	City   string
}

type person struct {
	Name    string
	Address address
}

type multiSecret struct {
	X int
	a int
	b string
}

func TestAssert_NilSlicesAreEqual(t *testing.T) {
	var s1 []int
	var s2 []int
	riteway.Assert(t, riteway.Case[[]int]{
		Given:    "two nil []int slices",
		Should:   "pass Assert without recording a failure",
		Actual:   s1,
		Expected: s2,
	})
}

func TestAssert_NilVsEmptySlice_Fails(t *testing.T) {
	ft := &fakeT{}
	var nilSlice []int
	emptySlice := []int{}
	riteway.Assert(ft, riteway.Case[[]int]{
		Given:    "nil []int vs empty []int{}",
		Should:   "record a failure because nil and empty slices differ",
		Actual:   nilSlice,
		Expected: emptySlice,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "nil []int vs empty []int{}",
		Should:   "mark fakeT as failed",
		Actual:   ft.failed,
		Expected: true,
	})
}

func TestAssert_EqualMaps(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	riteway.Assert(t, riteway.Case[map[string]int]{
		Given:    "two map[string]int with the same content",
		Should:   "pass Assert without recording a failure",
		Actual:   m1,
		Expected: m2,
	})
}

func TestAssert_UnequalMaps_RecordsFailure(t *testing.T) {
	ft := &fakeT{}
	m1 := map[string]int{"a": 1}
	m2 := map[string]int{"a": 2}
	riteway.Assert(ft, riteway.Case[map[string]int]{
		Given:    "two map[string]int with different values",
		Should:   "record a failure",
		Actual:   m1,
		Expected: m2,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "two map[string]int with different values",
		Should:   "mark fakeT as failed",
		Actual:   ft.failed,
		Expected: true,
	})
}

func TestAssert_NestedStructEqual(t *testing.T) {
	p1 := person{Name: "Alice", Address: address{Street: "1 Main St", City: "Springfield"}}
	p2 := person{Name: "Alice", Address: address{Street: "1 Main St", City: "Springfield"}}
	riteway.Assert(t, riteway.Case[person]{
		Given:    "two person structs with identical nested address fields",
		Should:   "pass Assert without recording a failure",
		Actual:   p1,
		Expected: p2,
	})
}

func TestAssert_MultipleOptions(t *testing.T) {
	a := multiSecret{X: 7, a: 1, b: "secret"}
	b := multiSecret{X: 7, a: 99, b: "different"}
	riteway.Assert(t, riteway.Case[multiSecret]{
		Given:    "two multiSecret structs with matching exported field and differing unexported fields",
		Should:   "pass when unexported fields are ignored via cmpopts.IgnoreUnexported",
		Actual:   a,
		Expected: b,
	}, cmpopts.IgnoreUnexported(multiSecret{}))
}

func TestAssert_TabAndNewlineGiven(t *testing.T) {
	ft := &fakeT{}
	riteway.Assert(ft, riteway.Case[int]{Given: "\t\n", Should: "something", Actual: 1, Expected: 1})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "Given field containing only tabs and newlines",
		Should:   "record a validation error showing the bad value",
		Actual:   ft.errors[0],
		Expected: "riteway.Assert: Given must not be empty (got \"\\t\\n\")",
	})
}

// Ensure cmp import is used (compile-time guard).
var _ = cmp.Equal

// --- exploratory edge case tests ---

func TestAssert_NilMapVsEmptyMap_Fails(t *testing.T) {
	ft := &fakeT{}
	var nilMap map[string]int
	emptyMap := map[string]int{}
	riteway.Assert(ft, riteway.Case[map[string]int]{
		Given:    "nil map vs empty map[string]int{}",
		Should:   "record a failure because nil and empty maps differ",
		Actual:   nilMap,
		Expected: emptyMap,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "nil map vs empty map[string]int{}",
		Should:   "mark fakeT as failed",
		Actual:   ft.failed,
		Expected: true,
	})
}

func TestAssert_AllowUnexported_EqualFields_Passes(t *testing.T) {
	a := configWithSecret{Port: 8080, secret: "same"}
	b := configWithSecret{Port: 8080, secret: "same"}
	riteway.Assert(t, riteway.Case[configWithSecret]{
		Given:    "two structs with identical exported and unexported fields",
		Should:   "pass when cmp.AllowUnexported is used",
		Actual:   a,
		Expected: b,
	}, cmp.AllowUnexported(configWithSecret{}))
}

func TestAssert_AllowUnexported_UnequalFields_Fails(t *testing.T) {
	ft := &fakeT{}
	a := configWithSecret{Port: 8080, secret: "alpha"}
	b := configWithSecret{Port: 8080, secret: "beta"}
	riteway.Assert(ft, riteway.Case[configWithSecret]{
		Given:    "two structs with differing unexported fields",
		Should:   "record a failure when cmp.AllowUnexported is used",
		Actual:   a,
		Expected: b,
	}, cmp.AllowUnexported(configWithSecret{}))
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "two structs with differing unexported fields",
		Should:   "mark fakeT as failed",
		Actual:   ft.failed,
		Expected: true,
	})
}

func TestTry_ZeroValueString(t *testing.T) {
	result, _ := riteway.Try(func() string { panic("oops") })
	riteway.Assert(t, riteway.Case[string]{
		Given:    "a Try[string] function that panics",
		Should:   "return empty string as zero value",
		Actual:   result,
		Expected: "",
	})
}

func TestTry_ZeroValuePointer(t *testing.T) {
	result, _ := riteway.Try(func() *int { panic("oops") })
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a Try[*int] function that panics",
		Should:   "return nil as zero value",
		Actual:   result == nil,
		Expected: true,
	})
}

func TestTry_PanicAnonymousStruct(t *testing.T) {
	_, err := riteway.Try(func() int {
		panic(struct{ Code int }{Code: 42})
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "a function that panics with an anonymous struct value",
		Should:   "include 'struct' in the error message type",
		Actual:   strings.Contains(err.Error(), "struct"),
		Expected: true,
	})
}

func TestTry_NilFn(t *testing.T) {
	_, err := riteway.Try[int](nil)
	riteway.Assert(t, riteway.Case[string]{
		Given:    "a nil function passed to Try",
		Should:   "return a descriptive error naming riteway.Try",
		Actual:   err.Error(),
		Expected: "riteway.Try: fn must not be nil",
	})
}

func TestTry_NestedTry(t *testing.T) {
	outer, outerErr := riteway.Try(func() string {
		_, innerErr := riteway.Try(func() int { panic("inner panic") })
		if innerErr == nil {
			return "wrong"
		}
		return "outer completed"
	})
	riteway.Assert(t, riteway.Case[string]{
		Given:    "nested Try where the inner function panics",
		Should:   "outer Try completes normally",
		Actual:   outer,
		Expected: "outer completed",
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "nested Try where the inner function panics",
		Should:   "outer error is nil",
		Actual:   outerErr == nil,
		Expected: true,
	})
}
