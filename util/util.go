package util

import "regexp"

// Keys returns the map's keys as a list.
func Keys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// EqualsIgnoreTrailingNewlines checks if s1 equals s2, but without
// taking trailing newlines into account.
func EqualsIgnoreTrailingNewlines(s1 string, s2 string) bool {
	trailingPattern := regexp.MustCompile("[\r\n]+$")
	return trailingPattern.ReplaceAllString(s1, "") == trailingPattern.ReplaceAllString(s2, "")
}
