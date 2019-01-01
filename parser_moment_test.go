package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDone(t *testing.T) {
	mom, _ := parseMom("[] blabla")
	assert.False(t, mom.IsDone())

	mom, _ = parseMom("[x] blabla")
	assert.True(t, mom.IsDone())

	mom, _ = parseMom("[X] blabla")
	assert.True(t, mom.IsDone())

	mom, _ = parseMom("[ x   ] blabla")
	assert.True(t, mom.IsDone())

	mom, _ = parseMom("[b] blabla")
	assert.False(t, mom.IsDone())
}

func TestBadDoneBrackets(t *testing.T) {
	_, err := parseMom("[x blabla")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Expected closing x for moment line [x blabla")
}

func TestPriority(t *testing.T) {
	mom, _ := parseMom("[] blabla!!!")
	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, 3, mom.GetPriority())

	smom, _ := parseSingleMom("[] blabla!! (1.2.2015)")
	assert.Equal(t, "blabla", smom.GetName())
	assert.Equal(t, 2, smom.GetPriority())
	assert.Equal(t, "01.02.2015 00:00", dateStr(smom.start))
}

func TestNoDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "nil", dateStr(mom.start))
	assert.Equal(t, "nil", dateStr(mom.end))
}

func TestSingleDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015)")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "24.12.2015 23:59", dateStr(mom.end))
}

func TestRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015-25.12.2015)")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "25.12.2015 23:59", dateStr(mom.end))
}

func TestEndlessRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015-)")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "nil", dateStr(mom.end))
}

func TestStartlessRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (-25.12.2015)")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "nil", dateStr(mom.start))
	assert.Equal(t, "25.12.2015 23:59", dateStr(mom.end))
}

func TestShortYearDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.15)")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "24.12.2015 23:59", dateStr(mom.end))
}

func TestZeroPaddedDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (04.01.15)")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "04.01.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "04.01.2015 23:59", dateStr(mom.end))
}

func TestNonZeroPaddedDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (4.1.15)")

	assert.Equal(t, "blabla", mom.name)
	assert.Equal(t, "04.01.2015 00:00", dateStr(mom.start))
	assert.Equal(t, "04.01.2015 23:59", dateStr(mom.end))
}

func TestFaultySingleDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (4.1.asfasf)")

	assert.Equal(t, "blabla (4.1.asfasf)", mom.name)
	assert.Equal(t, "nil", dateStr(mom.start))
	assert.Equal(t, "nil", dateStr(mom.end))
}

func TestFaultyRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (4.1.2015-asgdgd)")

	assert.Equal(t, "blabla (4.1.2015-asgdgd)", mom.name)
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

func parseMom(content string) (Moment, error) {
	line := &Line{content: content}
	return parseMoment(line, line.Content(), "")
}

func parseSingleMom(content string) (*SingleMoment, error) {
	mom, err := parseMom(content)
	return mom.(*SingleMoment), err
}
