package triebytesmapper

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Mapper builds a trie from a set of keywords and can be used to map a byte slice against these keywords.
type Mapper struct {
	opts               *Options
	keywordsOrig       []string
	keywordsNormalized []string
	root               *node
}

// New creates a new mapper with the given options and keywords.
func New(opts *Options, keywords ...string) *Mapper {
	if opts == nil {
		opts = &Options{}
	}

	keywordsNormalized := keywords
	if opts.NormalizeRune != nil {
		keywordsNormalized = make([]string, len(keywords))
		for i, k := range keywords {
			keywordsNormalized[i] = strings.Map(opts.NormalizeRune, k)
		}
	}

	m := &Mapper{
		keywordsOrig:       keywords,
		keywordsNormalized: keywordsNormalized,
		opts:               opts,
	}

	m.root = m.build()
	return m
}

// Keyword returns the keyword at the given index in the matches, nil if the index is out of range.
func (m Matches) Keyword(i int, src []byte) []byte {
	if i < 0 || i >= len(m) {
		return nil
	}
	return src[m[i].Lo:m[i].Hi]
}

// Matches is a slice of low (inclusively) and high (exclusively) slice indices.
type Matches []LoHi

// Map a byte slice against the keywords in the trie.
func (m *Mapper) Map(s []byte) Matches {
	matches := make(Matches, 0, 10)

	wasWordBoundary := false
	r, lo, w := rune(0), 0, 0
	for i := 0; i < len(s); i += w {
		r, w = utf8.DecodeRune(s[i:])
		if isWordBoundary(r) {
			if wasWordBoundary {
				// Skip consecutive word boundaries.
				lo = i + w
				continue
			}
			wasWordBoundary = true
			word := s[lo:i]
			match, more := m.MatchBytes(word)
			if match != "" {
				matches = append(matches, LoHi{lo, i})
			}
			if !more {
				// No more matches possible for this byte sequence.
				lo = i + w
			}
		} else {
			wasWordBoundary = false
		}
	}

	// Check last word.
	word := s[lo:]
	match, _ := m.MatchBytes(word)
	if match != "" {
		matches = append(matches, LoHi{lo, len(s)})
	}

	return matches
}

func isWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || unicode.IsPunct(r)
}

// MatchBytes matches a byte slice against the keywords in the trie.
// A non empty string is returned if the byte slice matches a keyword.
// The boolean return value indicates if there could be a match with
// more bytes (b is a prefix of a keyword in the trie).
func (m *Mapper) MatchBytes(b []byte) (string, bool) {

	node := m.root
	r, w := rune(0), 0
	for i := 0; i < len(b); i += w {
		r, w = utf8.DecodeRune(b[i:])
		if m.opts.NormalizeRune != nil {
			r = m.opts.NormalizeRune(r)
		}
		if len(node.Children) <= 0 {
			return "", false
		}
		if len(node.Children) <= int(r) {
			return "", false
		}
		node = node.Children[r]
		if node == nil {
			return "", false
		}
	}
	return node.Keyword, node.Keyword == ""
}

// LoHi is a low (inclusively) and high (exclusively) slice indices.
type LoHi struct {
	Lo int
	Hi int
}

// Options for the mapper.
type Options struct {
	// NormalizeRune will be applied to both the input and the keywords before matching.
	// Defaults to nil.
	// Typically used to  do lower casing and accent folding.
	NormalizeRune func(rune) rune
}

func (m *Mapper) build() *node {
	root := &node{}
	for i, keyword := range m.keywordsNormalized {
		n := root
		for _, r := range keyword {
			ri := int(r)
			if len(n.Children) <= ri {
				n.Children = append(n.Children, make([]*node, ri-len(n.Children)+1)...)
			}
			if n.Children[r] == nil {
				n.Children[r] = &node{}
			}
			n = n.Children[r]
		}
		n.Keyword = m.keywordsOrig[i]
	}

	return root
}

type node struct {
	Children []*node
	Keyword  string
}
