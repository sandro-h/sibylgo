package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleDate(t *testing.T) {
	line := &Line{content: "[] blabla (24.12.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "24.12.2015 23:59", dateStr(mom.end))
}

func TestRangeDate(t *testing.T) {
	line := &Line{content: "[] blabla (24.12.2015-25.12.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "25.12.2015 23:59", dateStr(mom.end))
}

func TestEndlessRangeDate(t *testing.T) {
	line := &Line{content: "[] blabla (24.12.2015-)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "nil", dateStr(mom.end))
}

func TestStartlessRangeDate(t *testing.T) {
	line := &Line{content: "[] blabla (-25.12.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "nil", dateStr(mom.start))
	assert.Equal(t, "25.12.2015 23:59", dateStr(mom.end))
}

func TestShortYearDate(t *testing.T) {
	line := &Line{content: "[] blabla (24.12.15)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "24.12.2015 23:59", dateStr(mom.end))
}

func TestZeroPaddedDate(t *testing.T) {
	line := &Line{content: "[] blabla (04.01.15)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "04.01.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "04.01.2015 23:59", dateStr(mom.end))
}

func TestNonZeroPaddedDate(t *testing.T) {
	line := &Line{content: "[] blabla (4.1.15)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "04.01.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "04.01.2015 23:59", dateStr(mom.end))
}

func TestFaultySingleDate(t *testing.T) {
	line := &Line{content: "[] blabla (4.1.asfasf)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "nil", dateStr(mom.start))
	assert.Equal(t, "nil", dateStr(mom.end))
}

func TestFaultyRangeDate(t *testing.T) {
	line := &Line{content: "[] blabla (4.1.2015-asgdgd)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, "nil", dateStr(mom.start))
	assert.Equal(t, "nil", dateStr(mom.end))
}

func TestCalculateSingleCoords(t *testing.T) {
	line := &Line{content: "[] blabla (4.1.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 11, mom.start.offset)
	assert.Equal(t, 8, mom.start.length)
	assert.Equal(t, 11, mom.end.offset)
	assert.Equal(t, 8, mom.end.length)

	line = &Line{content: "[] blabla (   4.1.2015  )"}
	mom, _ = parseSingleMoment(line, line.Content())
	assert.Equal(t, 14, mom.start.offset)
	assert.Equal(t, 8, mom.start.length)
	assert.Equal(t, 14, mom.end.offset)
	assert.Equal(t, 8, mom.end.length)
}

func TestCalculateRangeCoords(t *testing.T) {
	line := &Line{content: "[] blabla (4.1.2015-5.1.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 11, mom.start.offset)
	assert.Equal(t, 8, mom.start.length)
	assert.Equal(t, 20, mom.end.offset)
	assert.Equal(t, 8, mom.end.length)

	line = &Line{content: "[] blabla (  4.1.2015  -   5.1.2015  )"}
	mom, _ = parseSingleMoment(line, line.Content())

	assert.Equal(t, 13, mom.start.offset)
	assert.Equal(t, 8, mom.start.length)
	assert.Equal(t, 27, mom.end.offset)
	assert.Equal(t, 8, mom.end.length)
}

func TestCalculateEndlessRangeCoords(t *testing.T) {
	line := &Line{content: "[] blabla (  4.1.2015  -   )"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 13, mom.start.offset)
	assert.Equal(t, 8, mom.start.length)
}

func TestCalculateStartlessRangeCoords(t *testing.T) {
	line := &Line{content: "[] blabla (  -   5.1.2015  )"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 17, mom.end.offset)
	assert.Equal(t, 8, mom.end.length)
}

func dateStr(dt *Date) string {
	if dt == nil {
		return "nil"
	}
	return dt.time.Format("02.01.2006 15:04")
}
