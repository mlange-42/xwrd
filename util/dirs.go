package util

import (
	"os"
	"path/filepath"
)

const (
	rootDirName = ".xwrd"
	dictDirName = "dict"
	configName  = "config.yml"
	defaultDict = "german-700k.txt"
)

// RootDir returns the root storage directory
func RootDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, rootDirName)
}

// DictDir returns the dictionary storage directory
func DictDir() string {
	return filepath.Join(RootDir(), dictDirName)
}

// LanguageDir returns the dictionary storage directory for a language
func LanguageDir(lang string) string {
	return filepath.Join(RootDir(), dictDirName, lang)
}

// DictPath returns the path to a dictionary file
func DictPath(dict Dict) string {
	return filepath.Join(RootDir(), dictDirName, dict.Language, dict.Name+".lst")
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	return filepath.Join(RootDir(), configName)
}

// EnsureDirs creates storage directories if not present
func EnsureDirs() {
	err := CreateDir(RootDir())
	if err != nil {
		panic(err)
	}
	err = CreateDir(DictDir())
	if err != nil {
		panic(err)
	}
	for lang := range Dictionaries {
		path := filepath.Join(DictDir(), lang)
		err = CreateDir(path)
		if err != nil {
			panic(err)
		}
	}
}

// CreateDir creates directories recursively
func CreateDir(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
