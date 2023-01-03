package util

import (
	"io/ioutil"
	"strings"
)

// ReadWordList reads a file into a slice of words
func ReadWordList(dict Dict) ([]string, error) {
	fileContent, err := ioutil.ReadFile(DictPath(dict))
	if err != nil {
		return nil, err
	}
	words := strings.Split(string(fileContent), "\n")
	return words, nil
}
