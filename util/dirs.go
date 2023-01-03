package util

import (
	"os"
	"path/filepath"
)

const (
	rootDirName = ".xwrd"
	dictDirName = "dict"
	defaultDict = "default.txt"
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

// DictPath returns the path to the default dictionary
func DictPath() string {
	return filepath.Join(DictDir(), defaultDict)
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
