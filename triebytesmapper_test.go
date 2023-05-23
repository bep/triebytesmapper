package triebytesmapper_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/bep/triebytesmapper"
	qt "github.com/frankban/quicktest"
)

func Example() {
	// Testdata:
	// Dickens, Charles. A Christmas Carol. 1843.
	// 1000 random words from the MacOS dictionary.
	christmascarol, err := os.ReadFile("testdata/christmascarol.txt")
	if err != nil {
		log.Fatal(err)
	}
	thousandwords, err := os.ReadFile("testdata/thousandwords.txt")
	if err != nil {
		log.Fatal(err)
	}
	keywords := strings.Split(string(thousandwords), "\n")
	tolower := func(r rune) rune {
		return unicode.ToLower(r)
	}
	opts := &triebytesmapper.Options{NormalizeRune: tolower}
	m := triebytesmapper.New(opts, keywords...)

	matches := m.Map(christmascarol)
	first := matches.Keyword(0, christmascarol)

	fmt.Printf("Found %d matches. First match is %q.\n", len(matches), first)
	//  Output: Found 8 matches. First match is "Dickens".

}

func TestMapper(t *testing.T) {
	c := qt.New(t)
	m := triebytesmapper.New(nil, "foo", "bar", "baz")
	c.Assert(m, qt.IsNotNil)
	match, more := m.MatchBytes([]byte("foo"))
	c.Assert(match, qt.Equals, "foo")
	c.Assert(more, qt.IsFalse)
	match, more = m.MatchBytes([]byte("fo"))
	c.Assert(match, qt.Equals, "")
	c.Assert(more, qt.IsTrue)

	src := []byte("abc foo defg bar baz")
	matches := m.Map(src)
	c.Assert(matches, qt.HasLen, 3)

	c.Assert(matches.Keyword(0, src), qt.DeepEquals, []byte("foo"))
	c.Assert(matches.Keyword(1, src), qt.DeepEquals, []byte("bar"))
	c.Assert(matches.Keyword(2, src), qt.DeepEquals, []byte("baz"))
	c.Assert(matches.Keyword(3, src), qt.IsNil)

}

func TestMapperMixedCase(t *testing.T) {
	c := qt.New(t)
	tolower := func(r rune) rune {
		return unicode.ToLower(r)
	}
	opts := &triebytesmapper.Options{NormalizeRune: tolower}
	m := triebytesmapper.New(opts, "foo", "BAR", "baZ")
	match, _ := m.MatchBytes([]byte("bAr"))
	c.Assert(match, qt.Equals, "BAR")
}

func BenchmarkMap(b *testing.B) {
	content := []byte(strings.Repeat("abc foo defg asdfasdfa asdfas asdfasdfasd a sdf bar baz ", 1000))
	keywords := []string{"foo", "bar", "baz"}

	b.Run("No options", func(b *testing.B) {
		m := triebytesmapper.New(nil, keywords...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Map(content)
		}
	})

	b.Run("To lowercase", func(b *testing.B) {
		tolower := func(r rune) rune {
			return unicode.ToLower(r)
		}
		opts := &triebytesmapper.Options{NormalizeRune: tolower}
		m := triebytesmapper.New(opts, keywords...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = m.Map(content)
		}
	})

}

func BenchmarkChristmasCarrol(b *testing.B) {
	// Testdata:
	// Dickens, Charles. A Christmas Carol. 1843.
	// 1000 random words from the MacOS dictionary.
	christmascarol, err := os.ReadFile("testdata/christmascarol.txt")
	if err != nil {
		b.Fatal(err)
	}
	thousandwords, err := os.ReadFile("testdata/thousandwords.txt")
	if err != nil {
		b.Fatal(err)
	}
	keywords := strings.Split(string(thousandwords), "\n")
	m := triebytesmapper.New(nil, keywords...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Map(christmascarol)
	}
}
