# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-03-14

### Added

- `Case[T]` generic struct for structured test assertions answering the five
  questions every unit test must answer.
- `Assert[T]` function using `go-cmp` for deep equality and human-readable
  diffs. Validates that `Given` and `Should` are non-empty and non-whitespace.
- `Try[T]` function for recovering panics as errors. Correctly propagates
  `runtime.Goexit` so `t.FailNow` and `t.Fatal` behave normally inside
  `Try`-wrapped functions.
- `Match` for literal substring search. Empty substring always returns `""`
  to avoid ambiguity with "not found".
- `MatchRegexp` for regular expression search. Panics on invalid patterns.
