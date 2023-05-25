package triebytesmapper_test

import (
	"fmt"
	"log"
	"math/rand"
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
	// The thousandwords file is a list of words separated by newlines.
	// We need to split it into a slice of strings.
	keywords := strings.Split(string(thousandwords), "\n")
	// Trim any Windows line endings.
	for i, k := range keywords {
		keywords[i] = strings.TrimSpace(k)
	}
	tolower := func(r rune) rune {
		return unicode.ToLower(r)
	}
	opts := &triebytesmapper.Options{NormalizeRune: tolower}
	m := triebytesmapper.New(opts, keywords...)

	matches := m.Map(christmascarol)
	first := matches.Keyword(0, christmascarol)

	fmt.Printf("Found %d matches. First match is %q.\n", len(matches), first)
	//  Output: Found 11 matches. First match is "Dickens".

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

	src := []byte("abc foo defg bar. baz")
	matches := m.Map(src)
	c.Assert(matches, qt.HasLen, 3)

	c.Assert(matches.Keyword(0, src), qt.DeepEquals, []byte("foo"))
	c.Assert(matches.Keyword(1, src), qt.DeepEquals, []byte("bar"))
	c.Assert(matches.Keyword(2, src), qt.DeepEquals, []byte("baz"))
	c.Assert(matches.Keyword(3, src), qt.IsNil)
}

func TestMapperMultiBytesUnicode(t *testing.T) {
	c := qt.New(t)
	m := triebytesmapper.New(nil, "👍", "👎")
	c.Assert(m, qt.IsNotNil)
	match, _ := m.MatchBytes([]byte("👍"))
	c.Assert(match, qt.Equals, "👍")
	match, _ = m.MatchBytes([]byte("👎"))
	c.Assert(match, qt.Equals, "👎")

	src := []byte("abc 👍 defg 👎.")
	matches := m.Map(src)
	c.Assert(matches, qt.HasLen, 2)
	c.Assert(matches.Keyword(0, src), qt.DeepEquals, []byte("👍"))

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
	// https://poets.org/poem/do-not-go-gentle-good-night
	// Translated to runes with https://valhyr.com/pages/rune-converter
	// These runes should live in the 5792-5887 range of the unicode table.
	const (
		yeatsEnglish = `Do not go gentle into that good night, Old age should burn and rave at close of day; Rage, rage against the dying of the light. Though wise men at their end know dark is right, Because their words had forked no lightning they Do not go gentle into that good night.  Good men, the last wave by, crying how bright Their frail deeds might have danced in a green bay, Rage, rage against the dying of the light.  Wild men who caught and sang the sun in flight, And learn, too late, they grieved it on its way, Do not go gentle into that good night.  Grave men, near death, who see with blinding sight Blind eyes could blaze like meteors and be gay, Rage, rage against the dying of the light.  And you, my father, there on the sad height, Curse, bless, me now with your fierce tears, I pray. Do not go gentle into that good night. Rage, rage against the dying of the light.`
		yeatsRunes   = `ᛞᛟ ᚾᛟᛏ ᚷᛟ ᚷᛖᚾᛏᛚᛖ ᛁᚾᛏᛟ ᚦᚨᛏ ᚷᛟᛟᛞ ᚾᛁᚷᚺᛏ, ᛟᛚᛞ ᚨᚷᛖ ᛊᚺᛟᚢᛚᛞ ᛒᚢᚱᚾ ᚨᚾᛞ ᚱᚨᚢᛖ ᚨᛏ ᚲᛚᛟᛊᛖ ᛟᚠ ᛞᚨᛁ; ᚱᚨᚷᛖ, ᚱᚨᚷᛖ ᚨᚷᚨᛁᚾᛊᛏ ᚦᛖ ᛞᛁᛁᛜ ᛟᚠ ᚦᛖ ᛚᛁᚷᚺᛏ. ᛏᚺᛟᚢᚷᚺ ᚹᛁᛊᛖ ᛗᛖᚾ ᚨᛏ ᚦᛖᛁᚱ ᛖᚾᛞ ᚲᚾᛟᚹ ᛞᚨᚱᚲ ᛁᛊ ᚱᛁᚷᚺᛏ, ᛒᛖᚲᚨᚢᛊᛖ ᚦᛖᛁᚱ ᚹᛟᚱᛞᛊ ᚺᚨᛞ ᚠᛟᚱᚲᛖᛞ ᚾᛟ ᛚᛁᚷᚺᛏᚾᛁᛜ ᚦᛖᛁ ᛞᛟ ᚾᛟᛏ ᚷᛟ ᚷᛖᚾᛏᛚᛖ ᛁᚾᛏᛟ ᚦᚨᛏ ᚷᛟᛟᛞ ᚾᛁᚷᚺᛏ. ᚷᛟᛟᛞ ᛗᛖᚾ, ᚦᛖ ᛚᚨᛊᛏ ᚹᚨᚢᛖ ᛒᛁ, ᚲᚱᛁᛁᛜ ᚺᛟᚹ ᛒᚱᛁᚷᚺᛏ ᛏᚺᛖᛁᚱ ᚠᚱᚨᛁᛚ ᛞᛖᛖᛞᛊ ᛗᛁᚷᚺᛏ ᚺᚨᚢᛖ ᛞᚨᚾᚲᛖᛞ ᛁᚾ ᚨ ᚷᚱᛖᛖᚾ ᛒᚨᛁ, ᚱᚨᚷᛖ, ᚱᚨᚷᛖ ᚨᚷᚨᛁᚾᛊᛏ ᚦᛖ ᛞᛁᛁᛜ ᛟᚠ ᚦᛖ ᛚᛁᚷᚺᛏ. ᚹᛁᛚᛞ ᛗᛖᚾ ᚹᚺᛟ ᚲᚨᚢᚷᚺᛏ ᚨᚾᛞ ᛊᚨᛜ ᚦᛖ ᛊᚢᚾ ᛁᚾ ᚠᛚᛁᚷᚺᛏ, ᚨᚾᛞ ᛚᛖᚱᚾ, ᛏᛟᛟ ᛚᚨᛏᛖ, ᚦᛖᛁ ᚷᚱᛁᛖᚢᛖᛞ ᛁᛏ ᛟᚾ ᛁᛏᛊ ᚹᚨᛁ, ᛞᛟ ᚾᛟᛏ ᚷᛟ ᚷᛖᚾᛏᛚᛖ ᛁᚾᛏᛟ ᚦᚨᛏ ᚷᛟᛟᛞ ᚾᛁᚷᚺᛏ. ᚷᚱᚨᚢᛖ ᛗᛖᚾ, ᚾᛖᚱ ᛞᛖᚦ, ᚹᚺᛟ ᛊᛖᛖ ᚹᛁᚦ ᛒᛚᛁᚾᛞᛁᛜ ᛊᛁᚷᚺᛏ ᛒᛚᛁᚾᛞ ᛖᛁᛖᛊ ᚲᛟᚢᛚᛞ ᛒᛚᚨᛉᛖ ᛚᛁᚲᛖ ᛗᛖᛏᛖᛟᚱᛊ ᚨᚾᛞ ᛒᛖ ᚷᚨᛁ, ᚱᚨᚷᛖ, ᚱᚨᚷᛖ ᚨᚷᚨᛁᚾᛊᛏ ᚦᛖ ᛞᛁᛁᛜ ᛟᚠ ᚦᛖ ᛚᛁᚷᚺᛏ. ᚨᚾᛞ ᛁᛟᚢ, ᛗᛁ ᚠᚨᚦᛖᚱ, ᚦᛖᚱᛖ ᛟᚾ ᚦᛖ ᛊᚨᛞ ᚺᛖᛁᚷᚺᛏ, ᚲᚢᚱᛊᛖ, ᛒᛚᛖᛊᛊ, ᛗᛖ ᚾᛟᚹ ᚹᛁᚦ ᛁᛟᚢᚱ ᚠᛁᛖᚱᚲᛖ ᛏᛖᚱᛊ, ᛁ ᛈᚱᚨᛁ. ᛞᛟ ᚾᛟᛏ ᚷᛟ ᚷᛖᚾᛏᛚᛖ ᛁᚾᛏᛟ ᚦᚨᛏ ᚷᛟᛟᛞ ᚾᛁᚷᚺᛏ. ᚱᚨᚷᛖ, ᚱᚨᚷᛖ ᚨᚷᚨᛁᚾᛊᛏ ᚦᛖ ᛞᛁᛁᛜ ᛟᚠ ᚦᛖ ᛚᛁᚷᚺᛏ.`
	)

	createYeatsContentAndKeywords := func(sampleText string, lenContent, numKeywords int) ([]byte, []string) {
		content := sampleText
		for len(content) < lenContent {
			content = content + " " + sampleText
		}
		content = content[:lenContent]
		sampleTextWords := strings.Fields(string(sampleText))

		// To avoid getting a perfect match we insert some random words into the keywords.
		numDummyWOrds := numKeywords / 10
		for i := 0; i < numDummyWOrds; i++ {
			sampleTextWords = append(sampleTextWords, fmt.Sprintf("dummy%d", i))
		}

		// Shuffle the words.
		for i := range sampleTextWords {
			j := rand.Intn(i + 1)
			sampleTextWords[i], sampleTextWords[j] = sampleTextWords[j], sampleTextWords[i]
		}
		keywords := make([]string, numKeywords)
		for i := 0; i < numKeywords; i++ {
			keywords[i] = sampleTextWords[i%len(sampleTextWords)]
		}
		return []byte(string(content)), keywords
	}

	var basicContent, basicKeywords = createYeatsContentAndKeywords(yeatsEnglish, 1000, 30)

	b.Run("New", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m := triebytesmapper.New(nil, basicKeywords...)
			if m == nil {
				b.Fatal("m is nil")
			}
		}
	})

	b.Run("New tolower", func(b *testing.B) {
		tolower := func(r rune) rune {
			return unicode.ToLower(r)
		}
		opts := &triebytesmapper.Options{NormalizeRune: tolower}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m := triebytesmapper.New(opts, basicKeywords...)
			if m == nil {
				b.Fatal("m is nil")
			}
		}
	})

	b.Run("Map", func(b *testing.B) {
		m := triebytesmapper.New(nil, basicKeywords...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = m.Map(basicContent)
		}
	})

	b.Run("Map tolower", func(b *testing.B) {
		tolower := func(r rune) rune {
			return unicode.ToLower(r)
		}
		opts := &triebytesmapper.Options{NormalizeRune: tolower}
		m := triebytesmapper.New(opts, basicKeywords...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = m.Map(basicContent)
		}
	})

	for _, methodToTest := range []string{"Map", "New"} {
		for _, lang := range []string{"English", "Runes"} {
			sampleText := yeatsEnglish
			if lang == "Runes" {
				sampleText = yeatsRunes
			}
			for _, contentLength := range []int{1000, 10000} {
				for _, numKeywords := range []int{10, 100} {
					content, keywords := createYeatsContentAndKeywords(sampleText, contentLength, numKeywords)
					b.Run(fmt.Sprintf("%s %s Yeats %d words %d keywords", methodToTest, lang, contentLength, numKeywords), func(b *testing.B) {
						if methodToTest == "Map" {
							m := triebytesmapper.New(nil, keywords...)
							b.ResetTimer()
							for i := 0; i < b.N; i++ {
								m.Map(content)
							}
						} else {
							for i := 0; i < b.N; i++ {
								m := triebytesmapper.New(nil, keywords...)
								if m == nil {
									b.Fatal("m is nil")
								}
							}
						}

					})
				}
			}
		}
	}

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
