package anagram

// Letters used for the anagram tree. Sorted for memory-efficient trees
const Letters = "äöü?qxyjßvpzkfwbomgcluhdtarnise"

// Histogram returns a histogram of rune counts for a word
func Histogram(word string, letters map[rune]int, subtract bool, result []int) {
	unknown := letters['?']
	unit := 1
	if subtract {
		unit = -1
	}
	for _, char := range word {
		if idx, ok := letters[char]; ok {
			result[idx] += unit
		} else {
			result[unknown] += unit
		}
	}
}
