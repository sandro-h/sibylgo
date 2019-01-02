package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDaily(t *testing.T) {
	re := parseRe("[] bla (every day)")
	assert.NotNil(t, re)
	assert.Equal(t, RE_DAILY, re.recurrence)
	assert.NotNil(t, re.refDate)
	assert.Equal(t, 8, re.refDate.offset)
	assert.Equal(t, 9, re.refDate.length)
}

func TestWeekly(t *testing.T) {
	days := [...]interface{}{
		"sunday", time.Sunday,
		"monday", time.Monday,
		"tuesday", time.Tuesday,
		"wednesday", time.Wednesday,
		"thursday", time.Thursday,
		"friday", time.Friday,
		"saturday", time.Saturday}
	for i := 0; i < len(days); i += 2 {
		re := parseRe("[] bla (every " + days[i].(string) + ")")
		assert.NotNil(t, re)
		assert.Equal(t, RE_WEEKLY, re.recurrence)
		assert.Equal(t, days[i+1].(time.Weekday), re.refDate.time.Weekday())
	}
}

func TestMonthly(t *testing.T) {
	for i := 1; i <= 28; i++ {
		re := parseRe(fmt.Sprintf("[] bla (every %d.)", i))
		assert.NotNil(t, re)
		assert.Equal(t, RE_MONTHLY, re.recurrence)
		assert.Equal(t, i, re.refDate.time.Day())
	}
}

func TestYearly(t *testing.T) {
	re := parseRe("[] bla (every 2.5.)")
	assert.NotNil(t, re)
	assert.Equal(t, RE_YEARLY, re.recurrence)
	assert.Equal(t, 2, re.refDate.time.Day())
	assert.Equal(t, time.May, re.refDate.time.Month())
}

func TestInvalidRecurrence(t *testing.T) {
	re := parseRe("[] bla (every 2.5.2015)")
	assert.Nil(t, re)
}

func parseRe(content string) *Recurrence {
	line := &Line{content: content}
	re, _ := parseRecurrence(line, line.Content())
	return re
}
