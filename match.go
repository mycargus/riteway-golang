package riteway

import (
	"regexp"
	"strings"
)

// Match reports whether text contains substring. If found, it returns
// the substring. If not found, it returns an empty string.
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
// Patterns that can match the empty string (e.g., "x*") may return ""
// even when the pattern technically matches. Use anchored patterns
// (e.g., "x+") to avoid this ambiguity.
func MatchRegexp(text, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.FindString(text)
}
