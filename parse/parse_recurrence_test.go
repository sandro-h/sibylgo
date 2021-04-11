package parse

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sandro-h/sibylgo/moment"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	getNow = func() time.Time { return time.Now() }
	retCode := m.Run()
	os.Exit(retCode)
}

func TestDaily(t *testing.T) {
	re := parseRe("[] bla (every day)")
	assert.NotNil(t, re)
	assert.Equal(t, moment.RecurDaily, re.Recurrence)
	assert.NotNil(t, re.RefDate)
	assert.Equal(t, 8, re.RefDate.Offset)
	assert.Equal(t, 9, re.RefDate.Length)
}

func TestDailyWithTodayAlias(t *testing.T) {
	re := parseRe("[] bla (today)")
	assert.NotNil(t, re)
	assert.Equal(t, moment.RecurDaily, re.Recurrence)
	assert.NotNil(t, re.RefDate)
	assert.Equal(t, 8, re.RefDate.Offset)
	assert.Equal(t, 5, re.RefDate.Length)
}

func TestDailyWithDifferentConfig(t *testing.T) {
	defer ResetConfig()
	ParseConfig.SetDailyPattern("jeden tag")
	re := parseRe("[] bla (jeden tag)")
	assert.NotNil(t, re)
	assert.Equal(t, moment.RecurDaily, re.Recurrence)
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
		assert.Equal(t, moment.RecurWeekly, re.Recurrence)
		assert.Equal(t, days[i+1].(time.Weekday), re.RefDate.Time.Weekday())
	}
}

func TestWeeklyWithDifferentConfig(t *testing.T) {
	defer ResetConfig()
	ParseConfig.SetWeeklyPattern("jeden (montag|dienstag|mittwoch|donnerstag|freitag|samstag|sonntag)")
	ParseConfig.SetWeekDaysFromList([]string{"sonntag", "montag", "dienstag", "mittwoch", "donnerstag", "freitag", "samstag"})

	days := [...]interface{}{
		"sonntag", time.Sunday,
		"montag", time.Monday,
		"dienstag", time.Tuesday,
		"mittwoch", time.Wednesday,
		"donnerstag", time.Thursday,
		"freitag", time.Friday,
		"samstag", time.Saturday}
	for i := 0; i < len(days); i += 2 {
		re := parseRe("[] bla (jeden " + days[i].(string) + ")")
		assert.NotNil(t, re)
		assert.Equal(t, moment.RecurWeekly, re.Recurrence)
		assert.Equal(t, days[i+1].(time.Weekday), re.RefDate.Time.Weekday())
	}
}

func TestMonthly(t *testing.T) {
	for i := 1; i <= 28; i++ {
		re := parseRe(fmt.Sprintf("[] bla (every %d.)", i))
		assert.NotNil(t, re)
		assert.Equal(t, moment.RecurMonthly, re.Recurrence)
		assert.Equal(t, i, re.RefDate.Time.Day())
	}
}

func TestMonthlyWithDifferentConfig(t *testing.T) {
	defer ResetConfig()
	ParseConfig.SetMonthlyPattern(`jeden (\d{1,2})\.?$`)

	for i := 1; i <= 28; i++ {
		re := parseRe(fmt.Sprintf("[] bla (jeden %d.)", i))
		assert.NotNil(t, re)
		assert.Equal(t, moment.RecurMonthly, re.Recurrence)
		assert.Equal(t, i, re.RefDate.Time.Day())
	}
}

func TestYearly(t *testing.T) {
	re := parseRe("[] bla (every 2.5.)")
	assert.NotNil(t, re)
	assert.Equal(t, moment.RecurYearly, re.Recurrence)
	assert.Equal(t, 2, re.RefDate.Time.Day())
	assert.Equal(t, time.May, re.RefDate.Time.Month())
}

func TestYearlyWithDifferentConfig(t *testing.T) {
	defer ResetConfig()
	ParseConfig.SetYearlyPattern(`jeden (\d{1,2})\.(\d{1,2})\.?$`)

	re := parseRe("[] bla (jeden 2.5.)")
	assert.NotNil(t, re)
	assert.Equal(t, moment.RecurYearly, re.Recurrence)
	assert.Equal(t, 2, re.RefDate.Time.Day())
	assert.Equal(t, time.May, re.RefDate.Time.Month())
}

func TestBiWeekly(t *testing.T) {
	doTestNWeekly(t, "[] bla (every 2nd thursday)", moment.RecurBiWeekly, 2, tu.DtUtc("18.10.2019"), tu.DtUtc("17.10.2019"))
}

func TestTriWeekly(t *testing.T) {
	doTestNWeekly(t, "[] bla (every 3rd thursday)", moment.RecurTriWeekly, 3, tu.DtUtc("08.11.2019"), tu.DtUtc("07.11.2019"))
}

func TestQuadriWeekly(t *testing.T) {
	doTestNWeekly(t, "[] bla (every 4th thursday)", moment.RecurQuadriWeekly, 4, tu.DtUtc("01.11.2019"), tu.DtUtc("31.10.2019"))
}

func TestNthWeeklyWithDifferentConfig(t *testing.T) {
	defer ResetConfig()
	ParseConfig.SetWeekDaysFromList([]string{"sonntag", "montag", "dienstag", "mittwoch", "donnerstag", "freitag", "samstag"})
	ParseConfig.SetNthWeeklyPattern("jeden (2\\.|3\\.|4\\.) (montag|dienstag|mittwoch|donnerstag|freitag|samstag|sonntag)")
	ParseConfig.SetNthsFromList([]string{"2.", "3.", "4."})

	doTestNWeekly(t, "[] bla (jeden 2. donnerstag)", moment.RecurBiWeekly, 2, tu.DtUtc("18.10.2019"), tu.DtUtc("17.10.2019"))
	doTestNWeekly(t, "[] bla (jeden 3. donnerstag)", moment.RecurTriWeekly, 3, tu.DtUtc("08.11.2019"), tu.DtUtc("07.11.2019"))
	doTestNWeekly(t, "[] bla (jeden 4. donnerstag)", moment.RecurQuadriWeekly, 4, tu.DtUtc("01.11.2019"), tu.DtUtc("31.10.2019"))
}

func doTestNWeekly(t *testing.T, mom string, exRe int, n int, firstNow time.Time, exFirstRef time.Time) {
	// The important part is that the ref date is fixed within the n-range,
	// i.e. when we're in next week, it doesn't just move the ref date by one week,
	// otherwise we end up with weekly recurrence.
	now := firstNow
	expectedRef := exFirstRef
	for i := 0; i < 5; i++ {
		getNow = func() time.Time {
			return now
		}
		re := parseRe(mom)
		assert.NotNil(t, re)
		assert.Equal(t, exRe, re.Recurrence)
		assert.Equal(t, expectedRef, re.RefDate.Time, "In week of %s, expected ref date %s", now, expectedRef)
		now = now.AddDate(0, 0, 7)
		if i%n == n-1 {
			expectedRef = expectedRef.AddDate(0, 0, n*7)
		}
	}
}

func TestInvalidRecurrence(t *testing.T) {
	re := parseRe("[] bla (every 2.5.2015)")
	assert.Nil(t, re)
}

func parseRe(content string) *moment.Recurrence {
	line := &Line{content: content}
	re, _, _ := parseRecurrence(line, line.Content())
	return re
}
