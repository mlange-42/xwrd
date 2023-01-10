package anagram

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestHistogram(t *testing.T) {
	letters := []rune(Letters)
	lettersMap := make(map[rune]int, 2*len(letters))
	revLettersMap := map[int]rune{}
	for i, letter := range letters {
		lettersMap[letter] = i
		lettersMap[unicode.ToUpper(letter)] = i
		revLettersMap[i] = letter
	}

	tt := []struct {
		title    string
		text     string
		expected map[rune]int
	}{
		{
			title:    "one per letter",
			text:     "abc",
			expected: map[rune]int{'a': 1, 'b': 1, 'c': 1},
		},
		{
			title:    "multiple per letter",
			text:     "aaabcc",
			expected: map[rune]int{'a': 3, 'b': 1, 'c': 2},
		},
		{
			title:    "uppercase letters",
			text:     "ABC",
			expected: map[rune]int{'a': 1, 'b': 1, 'c': 1},
		},
		{
			title:    "special characters",
			text:     "abc+-",
			expected: map[rune]int{'a': 1, 'b': 1, 'c': 1, '?': 2},
		},
	}

	for _, test := range tt {
		res := make([]int, len(lettersMap), len(lettersMap))
		Histogram(test.text, lettersMap, false, res)

		resMap := map[rune]int{}
		for idx, rn := range revLettersMap {
			if res[idx] > 0 {
				resMap[rn] = res[idx]
			}
		}

		Histogram(test.text, lettersMap, false, res)
		assert.Equal(t, test.expected, resMap, "Wrong unique runes in %s", test.title)
	}
}
