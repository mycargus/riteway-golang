---
name: explore
description: This skill should be used when adding new functions, after refactoring, or when the user says "explore edge cases", "find bugs", "test boundaries", "what happens if...", "what are the edge cases", or "check for regressions". It runs exploratory edge-case tests against the riteway-golang API, discovers actual behavior, adds missing test cases to the test suite, and reports surprises.
version: 0.1.0
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

For each public function, check which edge case categories are already covered in the test files. See `references/edge-case-categories.md` for the full category list.

Compare existing test names and assertions against the categories. Note which categories have no coverage.

Focus your exploration on the gaps — don't re-test things the suite already covers.

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
- TestTry_NilFn — nil function panics, Try recovers with nil-deref message
    (mildly surprising — undocumented)
- ...

Findings:
1. [Mildly surprising] Try(nil) produces "runtime error: invalid memory
   address or nil pointer dereference" — technically correct but undocumented.
   Consider whether to document or guard against nil fn.
2. ...
```

## Scope

By default, explore all public functions. If the user specifies a function or area (e.g., "explore Try edge cases"), focus only on that.

## References

- `references/testing-patterns.md` — How to write tests in this project (fakeT, Try for panics, naming conventions)
- `references/edge-case-categories.md` — Edge case categories by function type
