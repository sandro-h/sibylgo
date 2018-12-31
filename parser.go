package main

import (
	"fmt"
	"os"
)

func ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := NewFileLineScanner(file)
	for scanner.Scan() {
		handleLine(scanner.Line())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func handleLine(line *Line) {
	fmt.Printf("%s\n", line.content)
}
