package riteway

import (
	"fmt"
	"regexp"
	"strings"
)

// Match reports whether text contains substring (case-sensitive). If found,
// it returns the substring. If not found, it returns an empty string.
//
// An empty substring argument always returns an empty string, even
// though the empty string is technically contained in every string.
// This avoids a confusing result where Match("anything", "") returns ""
// and is indistinguishable from "not found".
func Match(text, substring string) string {
	if substring == "" {
		return ""
	}
	if strings.Contains(text, substring) {
		return substring
	}
	return ""
}

// MatchRegexp reports whether text matches the regular expression
// pattern. If found, it returns the matched text. If not found, it
// returns an empty string. MatchRegexp panics if pattern is not a
// valid regular expression.
//
// An empty pattern argument always returns an empty string, matching
// Match's empty-argument guard.
//
// Patterns that can match the empty string (e.g., "x*", ".*") also
// panic, because their result would be indistinguishable from "not
// found". Use patterns that require at least one character (e.g., "x+").
//
// By default, . does not match newlines. Use the (?s) inline flag to
// enable dotall mode: MatchRegexp("foo\nbar", "(?s)foo.bar").
func MatchRegexp(text, pattern string) string {
	if pattern == "" {
		return ""
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("riteway.MatchRegexp: invalid pattern %q: %v", pattern, err))
	}
	if re.MatchString("") {
		panic(fmt.Sprintf("riteway.MatchRegexp: pattern %q can match an empty string, making the result indistinguishable from \"not found\"; use a pattern that requires at least one character", pattern))
	}
	return re.FindString(text)
}
