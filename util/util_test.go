package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualsIgnoreTrailingNewlines(t *testing.T) {
	type test struct {
		s1     string
		s2     string
		equals bool
	}

	tests := []test{
		{"abc\ndef", "abc\ndef", true},
		{"abc\ndef", "abc\ndef\n", true},
		{"abc\ndef", "abc\ndef\r\n", true},
		{"abc\ndef\n", "abc\ndef", true},
		{"abc\ndef\n", "abc\ndef\n\n\r\n", true},
		{"", "", true},
		{"", "\n", true},
		{"abc\ndef", "abc\nxyz", false},
		{"abc\n\ndef", "abc\ndef", false},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.equals, EqualsIgnoreTrailingNewlines(tc.s1, tc.s2), "'%s' ?= '%s'", tc.s1, tc.s2)
	}
}
