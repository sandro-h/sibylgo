package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func countStartWhitespaces(str string) int {
	for i, c := range str {
		if !unicode.IsSpace(c) {
			return i
		}
	}
	return utf8.RuneCountInString(str)
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

func parsePriority(str string) (int, string) {
	prio := 0
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] != priorityMark {
			break
		}
		prio++
	}
	return prio, strings.TrimSpace(str[0 : len(str)-prio])
}

func LastRuneIndex(s string, substr string) int {
	i := strings.LastIndex(s, substr)
	return utf8.RuneCountInString(s[:i])
}
