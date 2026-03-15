---
name: explore
description: This skill should be used when adding new functions, after refactoring, or when the user says "explore edge cases", "find bugs", "test boundaries", "what happens if...", "what are the edge cases", or "check for regressions". It runs exploratory edge-case tests against the riteway-golang API, discovers actual behavior, adds missing test cases to the test suite, and reports surprises.
version: 0.2.0
disable-model-invocation: false
user-invocable: true
---

# Exploratory Testing

Discover actual behavior by running real edge cases against the riteway-golang API. Every experiment becomes a permanent test case — nothing is throwaway.

## Mindset

You are not verifying expected behavior (that's what the existing tests do). You are discovering behavior at the boundaries — the places where assumptions break, documentation is silent, or two features interact unexpectedly.

## Process

### Step 1: Read source and existing tests

Read all source files and all test files. You need both to identify gaps.

Source files:
- `riteway.go` — Assert, Require, Try, Case
- `match.go` — Match, MatchRegexp

Test files:
- `riteway_test.go`
- `match_test.go`
- `fake_t_test.go` — fakeT spy (see `references/testing-patterns.md`)
- `example_test.go` — runnable examples and benchmarks (less relevant for edge cases, but shows existing coverage)

### Step 2: Identify coverage gaps

For each public function, derive what's worth testing from the source code itself: what types does it accept, what branches does it have, what can go wrong at the boundaries? Then check the existing tests to see what's already covered.

Ask for each function:
- What happens with nil inputs? Zero values? Empty strings?
- What happens when two similar-looking things are compared (nil map vs empty map, nil slice vs empty slice)?
- What happens with each panic path? Is the error message informative?
- What happens with unusual type parameters (T = string, *int, interface{}, anonymous struct)?
- What happens when features interact (nested Try, multiple opts, Goexit inside Try)?
- What happens at the extremes of documented behavior (e.g., patterns that barely qualify as zero-match)?

Focus on gaps — don't re-test what the suite already covers.

### Step 3: Write experiments as test functions

Write new test functions directly into the existing test files (`riteway_test.go` or `match_test.go`). Follow the patterns in `references/testing-patterns.md`.

Rules:
- **Never write throwaway files.** Every experiment is a permanent test case.
- **Check for duplicates before writing.** If an existing test covers the same scenario, skip it.
- **Use descriptive test names** that explain the edge case: `TestAssert_NilMapVsEmptyMap`, not `TestEdgeCase1`.
- **Use fakeT** for failure-path tests (capturing error messages without failing the runner).
- **Use Try** for panic-path tests.
- **One assertion per behavior.** Don't bundle unrelated checks into one test function.

### Step 4: Run the full test suite

Run `go test ./... -v` after adding new tests. All tests must pass — if a new test reveals a bug, that's a finding, but the test should still document the actual behavior (not the behavior you wish existed).

If a test fails because the actual behavior doesn't match your expectation, update the test to match reality and note it as a finding.

### Step 5: Report findings

For each new test added, report:
- **What it tests** (the edge case)
- **What happens** (actual behavior)
- **Surprise level**: Expected / Mildly surprising / Very surprising
- **Action needed**: None (behavior is fine) / Should document / Should fix

## What Good Output Looks Like

```
Added 7 new test cases, found 2 surprises:

New tests:
- TestAssert_NilMapVsEmptyMap — nil and empty maps are not equal (expected)
- TestTry_PanicAnonymousStruct — produces verbose type string in message
    (mildly surprising — should document)
- ...

Findings:
1. [Mildly surprising] panic(anonymous struct) produces a verbose type string
   like "panic(struct { x int }): { 42}" — technically correct but worth
   documenting so users know what to expect.
2. ...
```

## Scope

By default, explore all public functions. If the user specifies a function or area (e.g., "explore Try edge cases"), focus only on that.

## References

- `references/testing-patterns.md` — How to write tests in this project (fakeT, Try for panics, naming conventions)
