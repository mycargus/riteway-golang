# Architectural Decision Records

Decisions made during riteway-golang development, with context and rationale.

| # | Decision | File |
|---|----------|------|
| 1 | Go 1.21 as minimum version | [001](001-go-1-21-minimum.md) |
| 2 | Generic `Case[T]` struct for type-safe assertions | [002](002-generic-case-struct.md) |
| 3 | Accept `testing.TB` instead of `*testing.T` | [003](003-testing-tb-interface.md) |
| 4 | Use `go-cmp` for deep equality and diffs | [004](004-go-cmp-for-equality-and-diffs.md) |
| 5 | Separate `Match` and `MatchRegexp` functions | [005](005-separate-match-and-matchregexp.md) |
| 6 | `Match` with an empty substring returns `""` | [006](006-empty-substring-returns-empty-string.md) |
| 7 | `Try[T]` propagates `runtime.Goexit` | [007](007-try-propagates-goexit.md) |
| 8 | Validate `Given` and `Should` as non-empty, non-whitespace | [008](008-validate-given-and-should.md) |
| 9 | `fakeT` spy for testing failure paths | [009](009-faket-spy-for-failure-paths.md) |
| 10 | APIs from JavaScript riteway not ported to Go | [010](010-omitted-javascript-apis.md) |
| 11 | No publishing automation | [011](011-no-publishing-automation.md) |
| 12 | Signed commits required | [012](012-signed-commits.md) |
| 13 | Branch protection for main | [013](013-branch-protection.md) |
