package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueRunes(t *testing.T) {
	tt := []struct {
		title      string
		text       string
		ignoreCase bool
		expected   map[rune]int
	}{
		{
			title:      "simple case",
			text:       "abcA",
			ignoreCase: false,
			expected:   map[rune]int{'a': 1, 'b': 1, 'c': 1, 'A': 1},
		},
		{
			title:      "ignore case",
			text:       "abcA",
			ignoreCase: true,
			expected:   map[rune]int{'a': 2, 'b': 1, 'c': 1},
		},
		{
			title:      "repetitions",
			text:       "abc-abc",
			ignoreCase: true,
			expected:   map[rune]int{'a': 2, 'b': 2, 'c': 2, '-': 1},
		},
	}

	for _, test := range tt {
		res := UniqueRunes(test.text, test.ignoreCase)
		assert.Equal(t, test.expected, res, "Wrong unique runes in %s", test.title)
	}
}

func TestFindAdditions(t *testing.T) {
	tt := []struct {
		title    string
		orig     string
		altered  string
		expected []rune
	}{
		{
			title:    "no additions",
			orig:     "abc",
			altered:  "abc",
			expected: []rune{},
		},
		{
			title:    "no additions, different order",
			orig:     "abc",
			altered:  "bca",
			expected: []rune{},
		},
		{
			title:    "one addition",
			orig:     "abc",
			altered:  "abcd",
			expected: []rune("d"),
		},
		{
			title:    "two additions",
			orig:     "abc",
			altered:  "abcde",
			expected: []rune("de"),
		},
		{
			title:    "two additions somewhere",
			orig:     "abc",
			altered:  "adbec",
			expected: []rune("de"),
		},
		{
			title:    "repetitions",
			orig:     "abc",
			altered:  "abcccde",
			expected: []rune("ccde"),
		},
	}

	for _, test := range tt {
		runes := UniqueRunes(test.orig, true)
		res := FindAdditions(runes, test.altered, true)
		assert.Equal(t, test.expected, res, "Wrong reported additions in %s", test.title)
	}
}
