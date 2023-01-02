package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mlange-42/xwrd/anagram"
)

const version = "0.1.0-dev"

func main() {
	letters := []rune(anagram.Letters)
	tree := anagram.NewTree(letters)

	path := "./data/german.txt"
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	words := strings.Split(string(fileContent), "\n")
	_ = words

	tree.AddWords(words)

	fmt.Println(tree.FindAnagrams("Regen"))

	maxCount := 0
	maxLeaves := [][]string{}
	for _, l := range tree.Leaves {
		if len(l) > maxCount {
			maxCount = len(l)
			maxLeaves = [][]string{l}
		} else if len(l) == maxCount {
			maxLeaves = append(maxLeaves, l)
		}
	}
	fmt.Println(maxLeaves)

	maxLeaves = [][]string{}
	for _, l := range tree.Leaves {
		if len(l) > 5 {
			maxCount = len(l)
			maxLeaves = append(maxLeaves, l)
		}
	}
	for _, l := range maxLeaves {
		fmt.Println(l)
	}
}
