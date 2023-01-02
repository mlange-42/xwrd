package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mlange-42/xwrd/anagram"
)

const version = "0.1.0-dev"

func main() {
	tree := anagram.NewTree([]rune(anagram.Letters))

	path := "./data/german-700k.txt"
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	words := strings.Split(string(fileContent), "\n")

	tree.AddWords(words)

	ana := tree.MultiAnagrams("Martin")

	for _, res := range ana {
		fmt.Println(res)
	}
}
