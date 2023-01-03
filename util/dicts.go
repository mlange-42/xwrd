package util

import (
	"io"
	"net/http"
	"os"
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

// Dictionaries list dictionaries available for download
var Dictionaries = map[string][]Dict{
	"en": {
		{
			Name:     "yawl",
			Language: "en",
			URL:      "https://raw.githubusercontent.com/elasticdog/yawl/master/yawl-0.3.2.03/word.list",
		},
	},
	"de": {
		{
			Name:     "enz",
			Language: "de",
			URL:      "https://raw.githubusercontent.com/enz/german-wordlist/master/words",
		},
	},
}

// DownloadDict downloads a dictionary
func DownloadDict(dict Dict) error {
	path := DictPath(dict)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(dict.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
