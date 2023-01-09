package util

import (
	"unicode"
)

// UniqueRunes returns a map of unique runes in a string
func UniqueRunes(text string, ignoreCase bool) map[rune]int {
	result := map[rune]int{}
	for _, char := range text {
		if ignoreCase {
			char = unicode.ToLower(char)
		}
		if _, ok := result[char]; ok {
			result[char]++
		} else {
			result[char] = 1
		}
	}
	return result
}

// FindAdditions finds runes that appear in the second argument, but not in the first.
// The rune-to-count map `orig` is consumed/internally altered by the function!
func FindAdditions(orig map[rune]int, altered string, ignoreCase bool) []rune {
	result := []rune{}
	for _, char := range altered {
		if ignoreCase {
			char = unicode.ToLower(char)
		}
		if v, ok := orig[char]; ok {
			if v > 0 {
				orig[char]--
			} else {
				result = append(result, char)
			}
		} else {
			result = append(result, char)
		}
	}
	return result
}
