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
	root               *Node
}

// Keyword returns the keyword at the given index in the matches, nil if the index is out of range.
func (m Matches) Keyword(i int, src []byte) []byte {
	if i < 0 || i >= len(m) {
		return nil
	}
	return src[m[i].Lo : m[i].Hi+1]
}

// Matches is a slice of low (inclusively) and high (exclusively) slice indices.
type Matches []LoHi

// Map a byte slice against the keywords in the trie.
func (m *Mapper) Map(s []byte) Matches {
	matches := make(Matches, 0, 10)

	r, lo, w := rune(0), 0, 0
	for i := 0; i < len(s); i += w {
		r, w = utf8.DecodeRune(s[i:])
		if unicode.IsSpace(r) {
			word := s[lo:i]
			match, more := m.MatchBytes(word)
			if match != "" {
				matches = append(matches, LoHi{lo, i - 1})
			}
			if !more {
				// No more matches possible for this byte sequence.
				lo = i + w
			}
			continue
		}
	}

	// Check last word.
	word := s[lo:]
	match, _ := m.MatchBytes(word)
	if match != "" {
		matches = append(matches, LoHi{lo, len(s) - 1})
	}

	return matches
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
		if m.opts.RuneNormalizer != nil {
			r = m.opts.RuneNormalizer(r)
		}
		if _, ok := node.Children[r]; !ok {
			return "", false
		}
		node = node.Children[r]
	}
	return node.Keyword, node.Keyword == ""
}

type LoHi struct {
	Lo int
	Hi int
}

// Options for the mapper.
type Options struct {
	// Will be applied to both the input and the keywords before matching.
	// Defaults to nil.
	// Typically used to  do lower casing and accent folding.
	RuneNormalizer func(rune) rune
}

func New(opts *Options, keywords ...string) *Mapper {
	if opts == nil {
		opts = &Options{}
	}

	keywordsNormalized := keywords
	if opts.RuneNormalizer != nil {
		keywordsNormalized = make([]string, len(keywords))
		for i, k := range keywords {
			keywordsNormalized[i] = strings.Map(opts.RuneNormalizer, k)
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

type Node struct {
	Children map[rune]*Node
	Keyword  string
}

func (m *Mapper) build() *Node {
	root := &Node{Children: make(map[rune]*Node)}
	for i, keyword := range m.keywordsNormalized {
		node := root
		for _, r := range keyword {
			if _, ok := node.Children[r]; !ok {
				node.Children[r] = &Node{Children: make(map[rune]*Node)}
			}
			node = node.Children[r]
		}
		node.Keyword = m.keywordsOrig[i]
	}
	return root
}
