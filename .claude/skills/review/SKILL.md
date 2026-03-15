---
name: review
description: This skill should be used when asked to review riteway-golang from an engineer's perspective, think like a user of this library, look for UX issues, check for pitfalls, review the API, identify failure modes, or improve the developer experience. It performs a systematic walk-through of every public function, checking error messages, go-cmp integration, documentation accuracy, and real-world usage scenarios.
version: 0.2.0
disable-model-invocation: false
user-invocable: true
---

# Riteway Golang — Engineer Perspective Review

Perform a systematic review of riteway-golang from the perspective of an engineer who is using it — not building it. The goal is to find friction, surprises, and places where the library could make it easier to do the right thing and harder to do the wrong thing.

## Prerequisite

If test coverage looks thin for the areas you're reviewing, suggest running the explore skill (`/explore`) first. It adds edge-case tests to the suite, giving you concrete evidence for findings instead of mental simulation.

## Mindset

Think like an engineer encountering this library for the first time, or using it daily on a real project. Ask at every step:

- What would I assume this does?
- What would actually happen?
- Would I be confused, surprised, or blocked?
- Is this easy to use correctly and hard to use incorrectly?

## Review Process

### Step 1: Read the source files

Read all source files before forming opinions. Assumptions about implementation details are the most common source of missed issues.

Files to read (see `references/library-map.md` for the full map):
- `riteway.go`
- `match.go`

### Step 2: Read the test suite for evidence

Read all test files. The test suite — especially tests added by `/explore` — contains concrete evidence of actual behavior. Use test names, assertions, and captured error messages as the basis for your findings.

Do not mentally simulate behavior that a test already demonstrates. Cite the test.

### Step 3: Walk through the API surface

For each public function, consider it from the user's perspective. Reference specific test cases as evidence.

**Assert** — What does the failure message look like for simple vs complex types? What happens with whitespace-only Given/Should? Are go-cmp panics surfaced clearly? Is it obvious that Assert is non-fatal?

**Require** — Is it clear when to use Require vs Assert? Does the error message identify itself as "riteway.Require"? Is the fatal behavior correct when validation fails?

**Try** — Can I tell from the error that the failure was a panic vs a normal error return? Does it work with all type parameters? What does the zero value result look like when panic occurs?

**Match** — Is the "not found" signal (empty string) ambiguous in practice? Is case-sensitivity obvious?

**MatchRegexp** — Is the panic-on-zero-match behavior clear? Is the panic message for invalid patterns actionable? What does the first-match-only behavior mean for users?

### Step 4: Check go-cmp integration UX

Using test evidence, assess:
- How clear is the diff format (`-want +got`) for common types?
- Is `cmpopts.IgnoreUnexported` easy to discover and apply?
- Does the go-cmp panic on unexported fields produce a helpful error message?
- Can opts be passed in a composable way, or is there friction?

### Step 5: Check documentation vs reality

- Does the README example code match actual behavior?
- Are import paths accurate?
- Does the README explain what happens on failure?
- Is every panic condition documented?
- Is Go 1.21 requirement prominent?

### Step 6: Use the review checklist

Work through `references/review-checklist.md` systematically. Each section covers a category of issues. Mark items that surface real problems. Cross-reference with test evidence.

### Step 7: Score and report

Rate each issue found by:
- **Severity**: Critical / Significant / Moderate / Minor / Polish
- **Category**: API ergonomics / Error messages / go-cmp integration / Documentation / Generic type UX
- **Fix direction**: What would make this better?

Organize the report with the highest-severity issues first. Include specific code examples that reproduce each issue. Cite test names where applicable.

## What Good Output Looks Like

A good review surfaces concrete, specific issues — not vague observations. Each issue should include:

1. A specific scenario that triggers it
2. What the user expects
3. What actually happens (cite test name or actual output)
4. Why it matters (severity)
5. A suggested improvement

**Example of a good finding:**
> **[Significant] `Assert` is non-fatal but README doesn't say so**
> `TestAssert_MismatchedValues_ProducesDiff` shows that after a failed Assert,
> the test continues. A user expecting `t.Fatalf` behavior will be confused when
> subsequent assertions produce noise. Fix: add "Non-fatal: calls t.Errorf" to
> the README, and point users to `Require` for fatal behavior.

**Example of a weak finding to avoid:**
> "The API could be more intuitive." — Too vague, no actionable direction.

## Additional Resources

- **`references/library-map.md`** — All source files, test files, public API, key design decisions, and known non-obvious behaviors
- **`references/review-checklist.md`** — Systematic checklist organized by category
