package util

import (
	"os"
	"path/filepath"
)

// ReadFile reads a textfile's entire content into a string and returns it.
func ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes the passed string content to a textfile.
func WriteFile(filePath string, str string) error {
	if str[len(str)-1] != '\n' {
		str += "\n"
	}
	return os.WriteFile(filePath, []byte(str), 0644)
}

// AppendFile writes the passed string content at the end of the textfile.
func AppendFile(filePath string, str string) error {
	if str[len(str)-1] != '\n' {
		str += "\n"
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(str)
	return err
}

// RemoveExtension returns the filename without extension
// E.g. todo.txt -> todo
func RemoveExtension(s string) string {
	return s[:len(s)-len(filepath.Ext(s))]
}
