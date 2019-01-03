package parse

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDaily(t *testing.T) {
	re := parseRe("[] bla (every day)")
	assert.NotNil(t, re)
	assert.Equal(t, moment.RE_DAILY, re.Recurrence)
	assert.NotNil(t, re.RefDate)
	assert.Equal(t, 8, re.RefDate.Offset)
	assert.Equal(t, 9, re.RefDate.Length)
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
		assert.Equal(t, moment.RE_WEEKLY, re.Recurrence)
		assert.Equal(t, days[i+1].(time.Weekday), re.RefDate.Time.Weekday())
	}
}

func TestMonthly(t *testing.T) {
	for i := 1; i <= 28; i++ {
		re := parseRe(fmt.Sprintf("[] bla (every %d.)", i))
		assert.NotNil(t, re)
		assert.Equal(t, moment.RE_MONTHLY, re.Recurrence)
		assert.Equal(t, i, re.RefDate.Time.Day())
	}
}

func TestYearly(t *testing.T) {
	re := parseRe("[] bla (every 2.5.)")
	assert.NotNil(t, re)
	assert.Equal(t, moment.RE_YEARLY, re.Recurrence)
	assert.Equal(t, 2, re.RefDate.Time.Day())
	assert.Equal(t, time.May, re.RefDate.Time.Month())
}

func TestInvalidRecurrence(t *testing.T) {
	re := parseRe("[] bla (every 2.5.2015)")
	assert.Nil(t, re)
}

func parseRe(content string) *moment.Recurrence {
	line := &Line{content: content}
	re, _ := parseRecurrence(line, line.Content())
	return re
}
