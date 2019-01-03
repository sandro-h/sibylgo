package parse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadLines(t *testing.T) {
	sc := NewLineStringScanner(`line1
line2
line3
line4`)

	assert.True(t, sc.Scan())
	assertLine(t, 0, 0, "line1", sc.Line())
	assert.True(t, sc.Scan())
	assertLine(t, 1, 6, "line2", sc.Line())
	assert.True(t, sc.Scan())
	assertLine(t, 2, 12, "line3", sc.Line())
	assert.True(t, sc.Scan())
	assertLine(t, 3, 18, "line4", sc.Line())
	assert.False(t, sc.Scan())
}

func TestUnscan(t *testing.T) {
	sc := NewLineStringScanner(`line1
line2
line3
line4`)

	assert.True(t, sc.Scan())
	assertLine(t, 0, 0, "line1", sc.Line())
	assert.True(t, sc.Scan())
	assertLine(t, 1, 6, "line2", sc.Line())

	sc.Unscan()
	assert.True(t, sc.Scan())
	assertLine(t, 1, 6, "line2", sc.Line())

	assert.True(t, sc.Scan())
	assertLine(t, 2, 12, "line3", sc.Line())

	sc.Unscan()
	assert.True(t, sc.Scan())
	assertLine(t, 2, 12, "line3", sc.Line())

	assert.True(t, sc.Scan())
	assertLine(t, 3, 18, "line4", sc.Line())
	assert.False(t, sc.Scan())
}

func TestRepeatedUnscan(t *testing.T) {
	sc := NewLineStringScanner(`line1
line2
line3
line4`)

	sc.Scan()
	sc.Scan()
	sc.Scan()
	assertLine(t, 2, 12, "line3", sc.Line())

	// Unscan the same line repeatedly
	sc.Unscan()
	assert.True(t, sc.Scan())
	assertLine(t, 2, 12, "line3", sc.Line())
	sc.Unscan()
	assert.True(t, sc.Scan())
	assertLine(t, 2, 12, "line3", sc.Line())
	sc.Unscan()
	assert.True(t, sc.Scan())
	assertLine(t, 2, 12, "line3", sc.Line())

	assert.True(t, sc.Scan())
	assertLine(t, 3, 18, "line4", sc.Line())
	assert.False(t, sc.Scan())
}

func TestUnscanLast(t *testing.T) {
	sc := NewLineStringScanner(`line1
line2
line3
line4`)

	sc.Scan()
	sc.Scan()
	sc.Scan()
	sc.Scan()

	assertLine(t, 3, 18, "line4", sc.Line())
	sc.Unscan()
	assert.True(t, sc.Scan())
	assertLine(t, 3, 18, "line4", sc.Line())

	assert.False(t, sc.Scan())
}

func TestDetectNewLineLength(t *testing.T) {
	sc := NewLineStringScanner("line1\r\n" +
		"line2\r\n" +
		"line3\r\n" +
		"line4")

	assert.True(t, sc.Scan())
	assertLine(t, 0, 0, "line1", sc.Line())
	assert.True(t, sc.Scan())
	assertLine(t, 1, 7, "line2", sc.Line())
	assert.True(t, sc.Scan())
	assertLine(t, 2, 14, "line3", sc.Line())
	assert.True(t, sc.Scan())
	assertLine(t, 3, 21, "line4", sc.Line())
	assert.False(t, sc.Scan())
}

func assertLine(t *testing.T, expectedLineNum int, expectedOffset int,
	expectedContent string, line *Line) {

	assert.Equal(t, expectedContent, line.content)
	assert.Equal(t, expectedLineNum, line.lineNumber)
	assert.Equal(t, expectedOffset, line.offset)
}
