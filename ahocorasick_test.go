// ahocorasick_test.go: test suite for ahocorasick
//
// Copyright (c) 2013 CloudFlare, Inc.

package ahocorasick

import (
	"regexp"
	"strings"
	"testing"
)

func assert(t *testing.T, b bool) {
	if !b {
		t.Fail()
	}
}

func TestNoPatterns(t *testing.T) {
	m := NewStringMatcher([]string{})
	hits := m.Match([]byte("foo bar baz"))
	assert(t, len(hits) == 0)
}

func TestNoData(t *testing.T) {
	m := NewStringMatcher([]string{"foo", "baz", "bar"})
	hits := m.Match([]byte(""))
	assert(t, len(hits) == 0)
}

func TestSuffixes(t *testing.T) {
	m := NewStringMatcher([]string{"Superman", "uperman", "perman", "erman"})
	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == Hit{0, 18})
	assert(t, hits[1] == Hit{1, 19})
	assert(t, hits[2] == Hit{2, 20})
	assert(t, hits[3] == Hit{3, 21})
}

func TestPrefixes(t *testing.T) {
	m := NewStringMatcher([]string{"Superman", "Superma", "Superm", "Super"})
	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == Hit{3, 18})
	assert(t, hits[1] == Hit{2, 18})
	assert(t, hits[2] == Hit{1, 18})
	assert(t, hits[3] == Hit{0, 18})
}

func TestInterior(t *testing.T) {
	m := NewStringMatcher([]string{"Steel", "tee", "e"})
	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 6)
	assert(t, hits[0] == Hit{2, 2})
	assert(t, hits[1] == Hit{2, 13})
	assert(t, hits[2] == Hit{1, 12})
	assert(t, hits[3] == Hit{2, 14})
	assert(t, hits[4] == Hit{0, 11})
	assert(t, hits[5] == Hit{2, 21})
}

func TestMatchAtStart(t *testing.T) {
	m := NewStringMatcher([]string{"The", "Th", "he"})
	hits := m.Match([]byte("The Man Of Steel: Superman"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == Hit{1, 0})
	assert(t, hits[1] == Hit{0, 0})
	assert(t, hits[2] == Hit{2, 1})
}

func TestMatchAtEnd(t *testing.T) {
	m := NewStringMatcher([]string{"teel", "eel", "el"})
	hits := m.Match([]byte("The Man Of Steel"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == Hit{0, 12})
	assert(t, hits[1] == Hit{1, 13})
	assert(t, hits[2] == Hit{2, 14})
}

func TestOverlappingPatterns(t *testing.T) {
	m := NewStringMatcher([]string{"Man ", "n Of", "Of S"})
	hits := m.Match([]byte("The Man Of Steel"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == Hit{0, 4})
	assert(t, hits[1] == Hit{1, 6})
	assert(t, hits[2] == Hit{2, 8})
}

func TestMultipleMatches(t *testing.T) {
	m := NewStringMatcher([]string{"The", "Man", "an"})
	hits := m.Match([]byte("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 10)
	assert(t, hits[0] == Hit{1, 2})
	assert(t, hits[1] == Hit{2, 3})
	assert(t, hits[2] == Hit{2, 10})
	assert(t, hits[3] == Hit{2, 16})
	assert(t, hits[4] == Hit{2, 23})
	assert(t, hits[5] == Hit{1, 36})
	assert(t, hits[6] == Hit{2, 37})
	assert(t, hits[7] == Hit{2, 42})
	assert(t, hits[8] == Hit{0, 48})
	assert(t, hits[9] == Hit{2, 53})
}

func TestSingleCharacterMatches(t *testing.T) {
	m := NewStringMatcher([]string{"a", "M", "z"})
	hits := m.Match([]byte("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 13)
	assert(t, hits[0] == Hit{1, 2})
	assert(t, hits[1] == Hit{0, 3})
	assert(t, hits[2] == Hit{0, 10})
	assert(t, hits[3] == Hit{0, 16})
	assert(t, hits[4] == Hit{0, 18})
	assert(t, hits[5] == Hit{0, 23})
	assert(t, hits[6] == Hit{0, 25})
	assert(t, hits[7] == Hit{0, 27})
	assert(t, hits[8] == Hit{1, 36})
	assert(t, hits[9] == Hit{0, 37})
	assert(t, hits[10] == Hit{0, 42})
	assert(t, hits[11] == Hit{0, 53})
	assert(t, hits[12] == Hit{0, 55})
}

func TestNothingMatches(t *testing.T) {
	m := NewStringMatcher([]string{"baz", "bar", "foo"})
	hits := m.Match([]byte("A Man A Plan A Canal: Panama, which Man Planned The Canal"))
	assert(t, len(hits) == 0)
}

func TestWikipedia(t *testing.T) {
	m := NewStringMatcher([]string{"a", "ab", "bc", "bca", "c", "caa"})
	hits := m.Match([]byte("abccab"))
	assert(t, len(hits) == 7)
	assert(t, hits[0] == Hit{0, 0})
	assert(t, hits[1] == Hit{1, 0})
	assert(t, hits[2] == Hit{2, 1})
	assert(t, hits[3] == Hit{4, 2})
	assert(t, hits[4] == Hit{4, 3})
	assert(t, hits[5] == Hit{0, 4})
	assert(t, hits[6] == Hit{1, 4})

	hits = m.Match([]byte("bccab"))
	assert(t, len(hits) == 5)
	assert(t, hits[0] == Hit{2, 0})
	assert(t, hits[1] == Hit{4, 1})
	assert(t, hits[2] == Hit{4, 2})
	assert(t, hits[3] == Hit{0, 3})
	assert(t, hits[4] == Hit{1, 3})

	hits = m.Match([]byte("bccb"))
	assert(t, len(hits) == 3)
	assert(t, hits[0] == Hit{2, 0})
	assert(t, hits[1] == Hit{4, 1})
	assert(t, hits[2] == Hit{4, 2})
}

func TestMatch(t *testing.T) {
	m := NewStringMatcher([]string{"Mozilla", "Mac", "Macintosh", "Safari", "Sausage"})
	hits := m.Match([]byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 5)
	assert(t, hits[0] == Hit{0, 0})
	assert(t, hits[1] == Hit{1, 13})
	assert(t, hits[2] == Hit{2, 13})
	assert(t, hits[3] == Hit{1, 30})
	assert(t, hits[4] == Hit{3, 107})

	hits = m.Match([]byte("Mozilla/5.0 (Mac; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 4)
	assert(t, hits[0] == Hit{0, 0})
	assert(t, hits[1] == Hit{1, 13})
	assert(t, hits[2] == Hit{1, 24})
	assert(t, hits[3] == Hit{3, 101})

	hits = m.Match([]byte("Mozilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36"))
	assert(t, len(hits) == 2)
	assert(t, hits[0] == Hit{0, 0})
	assert(t, hits[1] == Hit{3, 106})

	hits = m.Match([]byte("Mozilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Sofari/537.36"))
	assert(t, len(hits) == 1)
	assert(t, hits[0] == Hit{0, 0})

	hits = m.Match([]byte("Mazilla/5.0 (Moc; Intel Computer OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Sofari/537.36"))
	assert(t, len(hits) == 0)
}

var bytes = []byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36")
var sbytes = string(bytes)
var dictionary = []string{"Mozilla", "Mac", "Macintosh", "Safari", "Sausage"}
var precomputed = NewStringMatcher(dictionary)

func BenchmarkMatchWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed.Match(bytes)
	}
}

func BenchmarkContainsWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary {
			if strings.Contains(sbytes, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re = regexp.MustCompile("(" + strings.Join(dictionary, "|") + ")")

func BenchmarkRegexpWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re.FindAllIndex(bytes, -1)
	}
}

var dictionary2 = []string{"Googlebot", "bingbot", "msnbot", "Yandex", "Baiduspider"}
var precomputed2 = NewStringMatcher(dictionary2)

func BenchmarkMatchFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed2.Match(bytes)
	}
}

func BenchmarkContainsFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary2 {
			if strings.Contains(sbytes, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re2 = regexp.MustCompile("(" + strings.Join(dictionary2, "|") + ")")

func BenchmarkRegexpFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re2.FindAllIndex(bytes, -1)
	}
}

var bytes2 = []byte("Firefox is a web browser, and is Mozilla's flagship software product. It is available in both desktop and mobile versions. Firefox uses the Gecko layout engine to render web pages, which implements current and anticipated web standards. As of April 2013, Firefox has approximately 20% of worldwide usage share of web browsers, making it the third most-used web browser. Firefox began as an experimental branch of the Mozilla codebase by Dave Hyatt, Joe Hewitt and Blake Ross. They believed the commercial requirements of Netscape's sponsorship and developer-driven feature creep compromised the utility of the Mozilla browser. To combat what they saw as the Mozilla Suite's software bloat, they created a stand-alone browser, with which they intended to replace the Mozilla Suite. Firefox was originally named Phoenix but the name was changed so as to avoid trademark conflicts with Phoenix Technologies. The initially-announced replacement, Firebird, provoked objections from the Firebird project community. The current name, Firefox, was chosen on February 9, 2004.")
var sbytes2 = string(bytes2)

var dictionary3 = []string{"Mozilla", "Mac", "Macintosh", "Safari", "Phoenix"}
var precomputed3 = NewStringMatcher(dictionary3)

func BenchmarkLongMatchWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed3.Match(bytes2)
	}
}

func BenchmarkLongContainsWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary3 {
			if strings.Contains(sbytes2, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re3 = regexp.MustCompile("(" + strings.Join(dictionary3, "|") + ")")

func BenchmarkLongRegexpWorks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re3.FindAllIndex(bytes2, -1)
	}
}

var dictionary4 = []string{"12343453", "34353", "234234523", "324234", "33333"}
var precomputed4 = NewStringMatcher(dictionary4)

func BenchmarkLongMatchFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed4.Match(bytes2)
	}
}

func BenchmarkLongContainsFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary4 {
			if strings.Contains(sbytes2, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re4 = regexp.MustCompile("(" + strings.Join(dictionary4, "|") + ")")

func BenchmarkLongRegexpFails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re4.FindAllIndex(bytes2, -1)
	}
}

var dictionary5 = []string{"12343453", "34353", "234234523", "324234", "33333", "experimental", "branch", "of", "the", "Mozilla", "codebase", "by", "Dave", "Hyatt", "Joe", "Hewitt", "and", "Blake", "Ross", "mother", "frequently", "performed", "in", "concerts", "around", "the", "village", "uses", "the", "Gecko", "layout", "engine"}
var precomputed5 = NewStringMatcher(dictionary5)

func BenchmarkMatchMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed5.Match(bytes)
	}
}

func BenchmarkContainsMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary4 {
			if strings.Contains(sbytes, s) {
				hits = append(hits, i)
			}
		}
	}
}

var re5 = regexp.MustCompile("(" + strings.Join(dictionary5, "|") + ")")

func BenchmarkRegexpMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re5.FindAllIndex(bytes, -1)
	}
}

func BenchmarkLongMatchMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precomputed5.Match(bytes2)
	}
}

func BenchmarkLongContainsMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hits := make([]int, 0)
		for i, s := range dictionary4 {
			if strings.Contains(sbytes2, s) {
				hits = append(hits, i)
			}
		}
	}
}

func BenchmarkLongRegexpMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re5.FindAllIndex(bytes2, -1)
	}
}
