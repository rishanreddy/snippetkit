package internal

import (
	"os"
	"path/filepath"
)

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// EnsureDirExists creates a directory if it doesn't exist
func EnsureDirExists(path string) error {
	return os.MkdirAll(filepath.Dir(path), os.ModePerm)
}

// WriteToFile writes content to a file
func WriteToFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
