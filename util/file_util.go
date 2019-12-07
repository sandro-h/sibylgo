package util

import (
	"io/ioutil"
)

// ReadFile reads a textfile's entire content into a string and returns it.
func ReadFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes the passed string content to a textfile.
func WriteFile(filePath string, str string) error {
	return ioutil.WriteFile(filePath, []byte(str+"\n"), 0644)
}