package util

import (
	"io/ioutil"
	"os"
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
	if str[len(str)-1] != '\n' {
		str += "\n"
	}
	return ioutil.WriteFile(filePath, []byte(str), 0644)
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
