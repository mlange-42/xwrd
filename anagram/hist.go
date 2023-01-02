package anagram

// Letters used for the anagram tree
const Letters = "abcdefghijklmnopqrstuvwxyzäöüß?"

// Histogram returns a histogram of rune counts for a word
func Histogram(word string, letters map[rune]int, result []int) {
	unknown := letters['?']
	for _, char := range word {
		if idx, ok := letters[char]; ok {
			result[idx]++
		} else {
			result[unknown]++
		}
	}
}
