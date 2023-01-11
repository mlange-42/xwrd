package anagram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	tree := NewTree([]rune(Letters))

	words := []string{
		"abc", "bca", "cab",
		"abcdef", "fedcba",
	}

	progress := make(chan int, 8)
	go tree.AddWords(words, progress)

	for pr := range progress {
		_ = pr
	}

	ana := tree.Anagrams("abc")
	assert.Equal(t, Leaf{"abc", "bca", "cab"}, ana, "Wrong anagrams")

	anaUk := tree.AnagramsWithUnknown("abc", 0, 0)
	assert.Equal(t, []Leaf{{"abc", "bca", "cab"}}, anaUk, "Wrong anagrams")

	anaUk = tree.AnagramsWithUnknown("abc", 0, 1)
	assert.Equal(t, []Leaf{{"abc", "bca", "cab"}}, anaUk, "Wrong anagrams")

	anaUk = tree.AnagramsWithUnknown("abc", 0, 3)
	assert.Equal(t, []Leaf{{"abc", "bca", "cab"}, {"abcdef", "fedcba"}}, anaUk, "Wrong anagrams")

	anaPt := tree.PartialAnagrams("abcdef", 0)
	assert.Equal(t, []Leaf{{"abc", "bca", "cab"}, {"abcdef", "fedcba"}}, anaPt, "Wrong anagrams")

	anaPt = tree.PartialAnagrams("abcdef", 4)
	assert.Equal(t, []Leaf{{"abcdef", "fedcba"}}, anaPt, "Wrong anagrams")

	anaPt = tree.PartialAnagramsWithUnknown("abcd", 0, 0, 0)
	assert.Equal(t, []Leaf{{"abc", "bca", "cab"}}, anaPt, "Wrong anagrams")

	anaPt = tree.PartialAnagramsWithUnknown("abcd", 0, 0, 3)
	assert.Equal(t, []Leaf{{"abc", "bca", "cab"}, {"abcdef", "fedcba"}}, anaPt, "Wrong anagrams")

	anaMult := tree.MultiAnagrams("abcabc", 0, 0, false)
	assert.Equal(t, [][]Leaf{{{"abc", "bca", "cab"}, {"abc", "bca", "cab"}}}, anaMult, "Wrong anagrams")
}
