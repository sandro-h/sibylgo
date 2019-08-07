package format

import (
	"github.com/sandro-h/sibylgo/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFoldForVSCode(t *testing.T) {
	todos, _ := parse.String(`[] foo
[] bar
	[] hoho
	hello
	ye
	[] bla
		[] gi
[] ruu
	[] nog
[] pow
	hihi
[] wee
	`)

	fold := FoldForVSCode(todos)
	assert.Equal(t, `1-6
7-8
9-10
`, fold)
}
