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

// countIndent counts the whitespace indent at the start of the string up to maxIndent.
// Spaces count as 1, Tabs count as tabSize indentation.
// Returns the indentation value and the actual number of whitespace characters.
func countIndent(str string, tabSize int, maxIndent int) (int, int) {
	indent := 0
	cnt := 0
	for _, c := range str {
		if c == '\t' {
			indent += tabSize
			cnt++
		} else if c == ' ' {
			indent++
			cnt++
		} else {
			break
		}

		if indent >= maxIndent {
			break
		}
	}
	return indent, cnt
}

func newParseError(line *Line, msg string, args ...interface{}) error {
	args = append([]interface{}{line.LineNumber() + 1}, args...)
	return fmt.Errorf("Line %d: "+msg, args...)
}

func parsePriority(str string) (int, string) {
	prio := 0
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] != ParseConfig.GetPriorityMark() {
			break
		}
		prio++
	}
	return prio, strings.TrimSpace(str[0 : len(str)-prio])
}

// LastRuneIndex returns the rune index of the last instance of substr in s, or -1 if substr is not present in s.
func LastRuneIndex(s string, substr string) int {
	i := strings.LastIndex(s, substr)
	if i < 0 {
		return i
	}
	return utf8.RuneCountInString(s[:i])
}
