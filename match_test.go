package riteway_test

import (
	"strings"
	"testing"

	riteway "github.com/mycargus/riteway-golang"
)

func TestMatch_FoundInMiddle(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "substring in the middle of text",
		Should:   "return the substring",
		Actual:   riteway.Match("hello world", "lo wo"),
		Expected: "lo wo",
	})
}

func TestMatch_FoundAtStart(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "substring at the start of text",
		Should:   "return the substring",
		Actual:   riteway.Match("hello world", "hello"),
		Expected: "hello",
	})
}

func TestMatch_FoundAtEnd(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "substring at the end of text",
		Should:   "return the substring",
		Actual:   riteway.Match("hello world", "world"),
		Expected: "world",
	})
}

func TestMatch_SubstringEqualsFullText(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "substring equal to the full text",
		Should:   "return the substring",
		Actual:   riteway.Match("hello", "hello"),
		Expected: "hello",
	})
}

func TestMatch_NotFound(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "substring not present in text",
		Should:   "return empty string",
		Actual:   riteway.Match("hello world", "xyz"),
		Expected: "",
	})
}

func TestMatch_EmptySubstring_NonEmptyText(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "empty substring with non-empty text",
		Should:   "return empty string to avoid ambiguity",
		Actual:   riteway.Match("hello", ""),
		Expected: "",
	})
}

func TestMatch_EmptyText_NonEmptySubstring(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "empty text with non-empty substring",
		Should:   "return empty string",
		Actual:   riteway.Match("", "hello"),
		Expected: "",
	})
}

func TestMatch_BothEmpty(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "both text and substring are empty",
		Should:   "return empty string",
		Actual:   riteway.Match("", ""),
		Expected: "",
	})
}

func TestMatch_CaseMismatch(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "substring with different case than text",
		Should:   "return empty string (case-sensitive)",
		Actual:   riteway.Match("Hello World", "hello"),
		Expected: "",
	})
}

func TestMatch_Unicode(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "unicode multi-byte substring",
		Should:   "return the substring",
		Actual:   riteway.Match("こんにちは世界", "世界"),
		Expected: "世界",
	})
}

func TestMatchRegexp_MatchesSubstring(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "pattern matching a substring",
		Should:   "return the matched text",
		Actual:   riteway.MatchRegexp("hello world", `w\w+`),
		Expected: "world",
	})
}

func TestMatchRegexp_MatchesFullText(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "pattern matching the full text",
		Should:   "return the full text",
		Actual:   riteway.MatchRegexp("hello", `h\w+`),
		Expected: "hello",
	})
}

func TestMatchRegexp_NotFound(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "pattern not matching text",
		Should:   "return empty string",
		Actual:   riteway.MatchRegexp("hello world", `\d+`),
		Expected: "",
	})
}

func TestMatchRegexp_InvalidPattern_Panics(t *testing.T) {
	_, err := riteway.Try(func() string {
		return riteway.MatchRegexp("text", "[invalid")
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "an invalid regexp pattern",
		Should:   "panic (caught via Try)",
		Actual:   err != nil,
		Expected: true,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "an invalid regexp pattern",
		Should:   "include 'riteway.MatchRegexp' in the panic message for context",
		Actual:   strings.Contains(err.Error(), "riteway.MatchRegexp"),
		Expected: true,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "an invalid regexp pattern",
		Should:   "include the bad pattern in the panic message for debuggability",
		Actual:   strings.Contains(err.Error(), "[invalid"),
		Expected: true,
	})
}

func TestMatchRegexp_EmptyPattern(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "empty pattern",
		Should:   "return empty string, matching Match's empty-substring guard",
		Actual:   riteway.MatchRegexp("hello", ""),
		Expected: "",
	})
}

func TestMatchRegexp_DotDoesNotMatchNewline(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "pattern with . and text containing a newline",
		Should:   "return empty string because . does not match \\n by default",
		Actual:   riteway.MatchRegexp("foo\nbar", "foo.bar"),
		Expected: "",
	})
}

func TestMatchRegexp_DotAllFlag(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "pattern with (?s) dotall flag and text containing a newline",
		Should:   "return the matched text spanning the newline",
		Actual:   riteway.MatchRegexp("foo\nbar", "(?s)foo.bar"),
		Expected: "foo\nbar",
	})
}

func TestMatchRegexp_ZeroMatchPattern_Panics(t *testing.T) {
	_, err := riteway.Try(func() string {
		return riteway.MatchRegexp("hello", `x*`)
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "pattern 'x*' that can match an empty string",
		Should:   "panic to prevent ambiguous result",
		Actual:   err != nil,
		Expected: true,
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "pattern 'x*' that can match an empty string",
		Should:   "include the pattern in the panic message",
		Actual:   strings.Contains(err.Error(), "x*"),
		Expected: true,
	})
}

func TestMatchRegexp_GreedyZeroMatchPattern_Panics(t *testing.T) {
	_, err := riteway.Try(func() string {
		return riteway.MatchRegexp("hello world", `.*`)
	})
	riteway.Assert(t, riteway.Case[bool]{
		Given:    "pattern '.*' that can match an empty string",
		Should:   "panic to prevent ambiguous result",
		Actual:   err != nil,
		Expected: true,
	})
}

func TestMatchRegexp_AnchoredPattern(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "anchored pattern 'x+' against text with x",
		Should:   "return the matched text",
		Actual:   riteway.MatchRegexp("fooxxx bar", `x+`),
		Expected: "xxx",
	})
}

func TestMatchRegexp_CaseSensitive(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "uppercase pattern against lowercase text",
		Should:   "return empty string (case-sensitive)",
		Actual:   riteway.MatchRegexp("hello", `HELLO`),
		Expected: "",
	})
}

// --- edge case tests ---

func TestMatch_RepeatedSubstring(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "text with a repeated substring",
		Should:   "return the first occurrence of the substring",
		Actual:   riteway.Match("ababab", "ab"),
		Expected: "ab",
	})
}

func TestMatch_SubstringWithNewline(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "substring that spans a newline character",
		Should:   "return the matched substring including the newline",
		Actual:   riteway.Match("line1\nline2", "1\nline"),
		Expected: "1\nline",
	})
}

func TestMatchRegexp_CaseInsensitiveFlag(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "pattern with inline case-insensitive flag (?i)",
		Should:   "return the matched text regardless of case",
		Actual:   riteway.MatchRegexp("Hello World", `(?i)hello`),
		Expected: "Hello",
	})
}

func TestMatchRegexp_ReturnsFirstMatch(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "text with multiple words matching \\w+at",
		Should:   "return the first match",
		Actual:   riteway.MatchRegexp("cat bat rat", `\w+at`),
		Expected: "cat",
	})
}

func TestMatchRegexp_DigitPattern(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "text containing a digit sequence surrounded by letters",
		Should:   "return the digit sequence",
		Actual:   riteway.MatchRegexp("abc 123 def", `\d+`),
		Expected: "123",
	})
}
