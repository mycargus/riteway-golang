package riteway_test

import (
	"fmt"
	"testing"

	riteway "github.com/mycargus/riteway-golang"
)

func ExampleMatch() {
	fmt.Println(riteway.Match("hello world", "world"))
	fmt.Println(riteway.Match("hello world", "xyz"))
	fmt.Println(riteway.Match("hello", ""))
	// Output:
	// world
	//
	//
}

func ExampleMatchRegexp() {
	fmt.Println(riteway.MatchRegexp("hello world", `w\w+`))
	fmt.Println(riteway.MatchRegexp("hello world", `\d+`))
	// Output:
	// world
	//
}

func BenchmarkAssert(b *testing.B) {
	c := riteway.Case[int]{
		Given:    "two equal integers",
		Should:   "pass without error",
		Actual:   42,
		Expected: 42,
	}
	for i := 0; i < b.N; i++ {
		riteway.Assert(b, c)
	}
}

func ExampleTry() {
	_, err := riteway.Try(func() int {
		panic("something went wrong")
	})
	fmt.Println(err)
	// Output:
	// something went wrong
}
