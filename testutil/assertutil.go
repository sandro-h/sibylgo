package testutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// AssertContains checks if the actual string contains the expected string.
func AssertContains(t *testing.T, expected string, actual string) {
	assert.True(t, strings.Contains(actual, expected),
		fmt.Sprintf("Should contain '%s', but was '%s'", expected, actual))
}
