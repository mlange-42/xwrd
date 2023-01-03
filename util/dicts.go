package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Dict represents a known word list
type Dict struct {
	Name     string
	Language string
	URL      string
}

// NewDict creates a new dictionary
func NewDict(langName string) Dict {
	parts := strings.Split(langName, "/")
	lang := "en"
	name := parts[0]
	if len(parts) > 1 {
		lang = parts[0]
		name = parts[1]
	}
	return Dict{
		Language: lang,
		Name:     name,
	}
}

// FullName returns the full name with language prefix
func (d *Dict) FullName() string {
	return fmt.Sprintf("%s/%s", d.Language, d.Name)
}

// Dictionaries list dictionaries available for download
var Dictionaries = map[string][]Dict{
	"de": {
		{
			Name:     "enz",
			Language: "de",
			URL:      "https://raw.githubusercontent.com/enz/german-wordlist/master/words",
		},
	},
	"en": {
		{
			Name:     "yawl",
			Language: "en",
			URL:      "https://raw.githubusercontent.com/elasticdog/yawl/master/yawl-0.3.2.03/word.list",
		},
	},
}

// LoadDictionary reads a file into a slice of words
func LoadDictionary(dict Dict) ([]string, error) {
	if !HasDictionary(dict) {
		return nil, fmt.Errorf("no dictionary '%s/%s'. Download with: xwrd dict install %[1]s/%[2]s", dict.Language, dict.Name)
	}
	fileContent, err := ioutil.ReadFile(DictPath(dict))
	if err != nil {
		return nil, err
	}
	words := strings.Split(string(fileContent), "\n")
	return words, nil
}

// HasDictionary checks if a dict exists
func HasDictionary(dict Dict) bool {
	path := DictPath(dict)
	return FileExists(path)
}

// AllDictionaries lists all installed dictionaries
func AllDictionaries() (map[string]Dict, error) {
	basePath := DictDir()

	languages, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	results := map[string]Dict{}

	for _, lang := range languages {
		if !lang.IsDir() {
			continue
		}
		langPath := filepath.Join(basePath, lang.Name())

		dicts, err := os.ReadDir(langPath)
		if err != nil {
			return nil, err
		}

		for _, dict := range dicts {
			if dict.IsDir() {
				continue
			}
			d := NewDict(lang.Name() + "/" + dict.Name())
			results[d.FullName()] = d
		}
	}
	return results, nil
}

// DownloadDictionary downloads a dictionary
func DownloadDictionary(dict Dict) error {
	path := DictPath(dict)

	resp, err := http.Get(dict.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
