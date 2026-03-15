# Edge Case Categories

Use this as a checklist when identifying coverage gaps. For each public function, check which categories are already covered by existing tests.

## Assert / Require

### Value types
- [ ] int (zero, positive, negative)
- [ ] string (empty, whitespace, unicode, multiline)
- [ ] bool
- [ ] nil pointer vs non-nil pointer
- [ ] nil slice vs empty slice
- [ ] nil map vs empty map
- [ ] equal maps (order-independent)
- [ ] unequal maps
- [ ] nested structs (deep equality)
- [ ] structs with unexported fields (with and without opts)
- [ ] interface values

### Validation
- [ ] Empty Given
- [ ] Whitespace-only Given
- [ ] Empty Should
- [ ] Whitespace-only Should
- [ ] Both Given and Should empty (early return — only Given error)
- [ ] Tab/newline-only Given

### go-cmp integration
- [ ] cmpopts.IgnoreUnexported
- [ ] cmp.AllowUnexported
- [ ] Multiple cmp.Options
- [ ] Custom cmp.Comparer
- [ ] go-cmp panic on unexported fields (no opts) — recovered as error
- [ ] go-cmp panic from custom option — no false unexported-field mention

### Diff format
- [ ] Simple type mismatch (int, string)
- [ ] Nested struct mismatch — readability
- [ ] Large struct mismatch — is it overwhelming?
- [ ] Nil pointer in diff

### Require-specific
- [ ] Happy path (passes)
- [ ] Fatal on mismatch (stops test)
- [ ] Fatal on empty Given (stops test)
- [ ] Error message says "riteway.Require" not "riteway.Assert"

## Try

### Panic types
- [ ] panic(string) — message preserved as error
- [ ] panic(error) — errors.Is identity preserved
- [ ] panic(int) — formatted as "panic(int): 42"
- [ ] panic(nil) on Go 1.21+ — *runtime.PanicNilError
- [ ] panic(struct{...}) — anonymous struct in message

### Return values
- [ ] No panic — result returned, err nil
- [ ] Panic — result is zero value of T
- [ ] Try[string] zero value
- [ ] Try[*int] zero value (nil)
- [ ] Try[struct{}] zero value

### Edge cases
- [ ] nil fn — recovers nil pointer dereference
- [ ] Nested Try — inner panic doesn't affect outer
- [ ] fn calls runtime.Goexit (t.FailNow) — re-propagated

## Match

### Inputs
- [ ] Substring found at start, middle, end
- [ ] Substring equals full text
- [ ] Substring not found
- [ ] Empty substring (returns "")
- [ ] Empty text
- [ ] Both empty
- [ ] Case mismatch (case-sensitive)
- [ ] Unicode / multibyte characters
- [ ] Substring with newline
- [ ] Repeated substring

## MatchRegexp

### Patterns
- [ ] Pattern matching substring
- [ ] Pattern matching full text
- [ ] Pattern not matching
- [ ] Empty pattern (returns "")
- [ ] Invalid pattern (panics with clear message)
- [ ] Zero-match pattern like x* (panics)
- [ ] Greedy zero-match pattern like .* (panics)
- [ ] Anchored pattern (x+)

### Flags
- [ ] . does not match newline by default
- [ ] (?s) dotall flag
- [ ] (?i) case-insensitive flag

### Behavior
- [ ] Returns first match only
- [ ] Case-sensitive by default
- [ ] Digit pattern (\d+)
