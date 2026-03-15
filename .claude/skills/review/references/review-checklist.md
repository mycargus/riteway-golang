# Riteway Golang — Review Checklist

Use this checklist systematically when reviewing riteway-golang from an engineer's perspective.

**Important:** Do not mentally simulate behavior that is already covered by a test. Cite the test. If a scenario has no test coverage, suggest running `/explore` to fill the gap.

---

## 1. API Ergonomics

### Assert
- [ ] Can a user accidentally pass the wrong type for `Actual` vs `Expected`? What does the compiler error look like?
- [ ] What happens with whitespace-only `Given` or `Should`? Is the error message actionable?
- [ ] Does the failure message include both `Given`/`Should` context AND the diff?
- [ ] What does the diff look like for a simple string mismatch? For a nested struct? For a nil pointer?
- [ ] Is it clear from the README that Assert is non-fatal (`t.Errorf`)?
- [ ] Is the `opts ...cmp.Option` parameter discoverable? Would a new user know it exists?
- [ ] What happens when go-cmp panics due to unexported fields — is the error message clear?

### Require
- [ ] Is it clear when to use Require vs Assert?
- [ ] Does the error message say "riteway.Require" (not "riteway.Assert")?
- [ ] Does validation failure (empty Given) also halt the test?
- [ ] Is the Require example in the README compelling — does it show the value over Assert?

### Try
- [ ] Does `Try[T]` work with all type parameters — string, int, struct, pointer, interface?
- [ ] What is the zero value result when a panic occurs? Is this documented?
- [ ] What does the error message look like for `panic("message")` vs `panic(someError)` vs `panic(42)`?
- [ ] Can I tell from the error that the failure was a panic vs a normal error return?
- [ ] Does `Try` work correctly when the wrapped function calls `t.FailNow()` or `t.Fatal()`?
- [ ] What happens if `fn` is nil?
- [ ] Can I nest `Try` calls?

### Match
- [ ] Does `Match("anything", "")` return `""` or `"anything"`? Is the behavior surprising?
- [ ] Does `Match("", "pattern")` return `""` without panicking?
- [ ] Is the "not found" signal (empty string) ambiguous when the expected match is itself an empty string?
- [ ] Is case sensitivity clear from the README and godoc?
- [ ] Would a user know to assert `!= ""` to check "was found"?

### MatchRegexp
- [ ] Does `MatchRegexp("hello", "")` return `""`?
- [ ] Does `MatchRegexp("hello", "x*")` panic as expected?
- [ ] Is the panic-on-zero-match behavior documented in the README?
- [ ] Is the panic-on-invalid-regex surfaced with a helpful message?
- [ ] Does `.` match newlines by default?
- [ ] Is "first match only" behavior documented?
- [ ] Would a user know about `(?s)` and `(?i)` flags?

---

## 2. go-cmp Integration

- [ ] Is `cmpopts.IgnoreUnexported` easy to discover from the README or error messages?
- [ ] Does the panic message for unexported fields mention the field name and struct type?
- [ ] Is the `(-want +got)` diff format consistent with what Go developers expect?
- [ ] Are `cmp.Option` examples in the README?
- [ ] Can users pass multiple options composably?
- [ ] For large nested structs, is the diff readable or overwhelming?

---

## 3. Error Messages

For each error a user might encounter, ask:
- Is it obvious what went wrong?
- Is it obvious what to do next?
- Does it point to the right location (user code, not library internals)?

Key error paths to check:
- [ ] Empty `Given` field — includes bad value in message
- [ ] Whitespace-only `Should` field — includes bad value in message
- [ ] Assertion failure — includes Given/Should context and diff
- [ ] Require failure — same format as Assert but halts the test
- [ ] go-cmp panic on unexported field — includes field name and struct type
- [ ] `MatchRegexp` with invalid pattern — includes pattern and parse error
- [ ] `MatchRegexp` with zero-match pattern — explains why it panicked
- [ ] `Try` wrapping `panic(nil)` on Go 1.21+
- [ ] `Try` wrapping a non-error panic type — includes `%T` type name

---

## 4. Documentation vs Reality

- [ ] Does the README example code match actual behavior?
- [ ] Are the import paths in the README accurate?
- [ ] Is the `Case[T]` struct field order in the README consistent with the source?
- [ ] Does the README explain that Assert is non-fatal and Require is fatal?
- [ ] Is the `Try`/panic behavior clearly documented (what types are caught, what the error looks like)?
- [ ] Does the README API section for each function match its godoc? Check for behaviors documented in godoc that are absent from the README (zero-value results, nil guards, default behaviors).
- [ ] Is the `Match` empty-substring behavior documented?
- [ ] Is the `Match` case-sensitivity noted?
- [ ] Is the `MatchRegexp` panic-on-invalid AND panic-on-zero-match documented?
- [ ] Is the go-cmp unexported-field behavior documented with the fix?
- [ ] Is Go 1.21 requirement documented prominently?
- [ ] Does the README show how to use `cmpopts.IgnoreUnexported`?
- [ ] Does the README link to `regexp/syntax` for inline flags?

---

## 5. Make Right Things Easy, Wrong Things Hard

For each pattern below, assess difficulty (Easy / Medium / Hard / Impossible):

| Pattern | Should be | Assess |
|---------|-----------|--------|
| Write a passing assertion | Easy | |
| Write a failing assertion with full context | Easy | |
| Use `Require` for a guard assertion | Easy | |
| Use `Try` to test a panic | Easy | |
| Pass mismatched types to `Assert` | Impossible (compiler) | |
| Forget `Given` or `Should` | Hard (caught at runtime) | |
| Use whitespace-only `Given` | Hard (caught at runtime) | |
| Compare structs with unexported fields (without opts) | Medium (go-cmp panics) | |
| Use `MatchRegexp` with invalid pattern | Hard (panics immediately) | |
| Use `MatchRegexp` with zero-match pattern (e.g., `x*`) | Hard (panics immediately) | |
| Use `Try` and miss that result is zero value on panic | Medium | |
| Use `Assert` thinking it stops the test (it doesn't) | Easy (README says non-fatal) | |

---

## 6. Real-World Scenarios

Walk through these scenarios using test evidence. Cite specific tests that demonstrate the behavior.

### New engineer onboarding
1. `go get github.com/mycargus/riteway-golang`
2. Write first test using `riteway.Assert`
3. Make the test fail intentionally — is the output clear?
4. What does the diff look like?

### Testing an error/panic case
1. Write a function that panics
2. Use `riteway.Try` to capture it
3. Assert on the error message
4. What if the function returns a value and panics — what is `result`?

### Using Require as a guard
1. Require a precondition (e.g., no error from setup)
2. Assert on the result — does this only run if Require passed?
3. Is the pattern obvious from the README example?

### Testing string content
1. Use `riteway.Match` to search rendered output
2. Assert on what was found (`!= ""`)
3. Assert something was NOT found (`== ""`)
4. What does a failure message look like for a Match assertion?

### Testing structs with unexported fields
1. Write a test that compares structs with unexported fields
2. Observe the go-cmp panic
3. Apply `cmpopts.IgnoreUnexported` — how discoverable is this fix?

### Table-driven tests
1. Use `Assert` inside a loop over test cases
2. Make one case fail — does the output identify which case?
3. Does `t.Run` help? Is this documented?
