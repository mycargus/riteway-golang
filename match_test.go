package riteway_test

import (
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
}

func TestMatchRegexp_PatternMatchesEmpty(t *testing.T) {
	riteway.Assert(t, riteway.Case[string]{
		Given:    "pattern 'x*' that can match empty string",
		Should:   "return empty string (ambiguous match documented behavior)",
		Actual:   riteway.MatchRegexp("hello", `x*`),
		Expected: "",
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
