package main

import (
	"fmt"
	"unicode"
)

func countStartWhitespaces(str string) int {
	for i, c := range str {
		if !unicode.IsSpace(c) {
			return i
		}
	}
	return len(str)
}

func lengthWithoutStartEndWhitespaces(str string) int {
	runes := []rune(str)
	st := -1
	en := -1
	for i := 0; i < len(runes); i++ {
		if st == -1 && !unicode.IsSpace(runes[i]) {
			st = i
		}
		if en == -1 && !unicode.IsSpace(runes[len(runes)-1-i]) {
			en = i
		}
	}
	if st == -1 {
		return 0
	}
	return len(runes) - st - en
}

func newParseError(line *Line, msg string, args ...interface{}) error {
	args = append([]interface{}{line.LineNumber() + 1}, args...)
	return fmt.Errorf("Line %d: "+msg, args...)
}
