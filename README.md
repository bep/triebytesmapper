[![Tests on Linux, MacOS and Windows](https://github.com/bep/triebytesmapper/workflows/Test/badge.svg)](https://github.com/bep/triebytesmapper/actions?query=workflow:Test)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/triebytesmapper)](https://goreportcard.com/report/github.com/bep/triebytesmapper)
[![GoDoc](https://godoc.org/github.com/bep/triebytesmapper?status.svg)](https://godoc.org/github.com/bep/triebytesmapper)

Mapper builds a trie from a set of keywords and can be used to map a byte slice (typically file content) against these keywords.

There are some settings:

```go
// Options for the mapper.
type Options struct {
	// NormalizeRune will be applied to both the input and the keywords before matching.
	// Defaults to nil.
	// Typically used to  do lower casing and accent folding.
	NormalizeRune func(rune) rune

	// IsWordBoundary is used to determine word boundaries.
	// A default implementation is used if nil.
	IsWordBoundary func(rune) bool
}
```

Note that this library currently does not work for CJK languages.