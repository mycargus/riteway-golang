package riteway_test

import (
	"fmt"
	"runtime"
	"testing"
)

type fakeT struct {
	testing.TB
	errors []string
	failed bool
}

func (f *fakeT) Helper() {}

func (f *fakeT) Error(args ...any) {
	f.errors = append(f.errors, fmt.Sprint(args...))
	f.failed = true
}

func (f *fakeT) Errorf(format string, args ...any) {
	f.errors = append(f.errors, fmt.Sprintf(format, args...))
	f.failed = true
}

func (f *fakeT) FailNow() {
	f.failed = true
	runtime.Goexit()
}

func (f *fakeT) Fatal(args ...any) {
	f.errors = append(f.errors, fmt.Sprint(args...))
	f.failed = true
	runtime.Goexit()
}

func (f *fakeT) Fatalf(format string, args ...any) {
	f.errors = append(f.errors, fmt.Sprintf(format, args...))
	f.failed = true
	runtime.Goexit()
}

func (f *fakeT) Log(args ...any)                 {}
func (f *fakeT) Logf(format string, args ...any) {}
func (f *fakeT) Name() string                    { return "fakeT" }
