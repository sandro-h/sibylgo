package parse

import (
	"testing"

	"github.com/sandro-h/sibylgo/moment"
	"github.com/stretchr/testify/assert"
)

func TestDone(t *testing.T) {
	mom, _ := parseMom("[] blabla")
	assert.False(t, mom.IsDone())

	mom, _ = parseMom("[x] blabla")
	assert.True(t, mom.IsDone())

	mom, _ = parseMom("[ x   ] blabla")
	assert.True(t, mom.IsDone())

	mom, _ = parseMom("[b] blabla")
	assert.False(t, mom.IsDone())
}

func TestBadDoneBrackets(t *testing.T) {
	_, err := parseMom("[x blabla")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Expected closing ] for moment line [x blabla")
}

func TestWorkState(t *testing.T) {
	mom, _ := parseMom("[] blabla")
	assert.Equal(t, moment.NewState, mom.GetWorkState())

	mom, _ = parseMom("[x] blabla")
	assert.Equal(t, moment.DoneState, mom.GetWorkState())

	mom, _ = parseMom("[w] blabla")
	assert.Equal(t, moment.WaitingState, mom.GetWorkState())

	mom, _ = parseMom("[p] blabla")
	assert.Equal(t, moment.InProgressState, mom.GetWorkState())
}

func TestPriority(t *testing.T) {
	mom, _ := parseMom("[] blabla!!!")
	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, 3, mom.GetPriority())

	smom, _ := parseSingleMom("[] blabla!! (1.2.2015)")
	assert.Equal(t, "blabla", smom.GetName())
	assert.Equal(t, 2, smom.GetPriority())
	assert.Equal(t, "01.02.2015 00:00", dateStr(smom.Start))
}

func TestNoDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "nil", dateStr(mom.Start))
	assert.Equal(t, "nil", dateStr(mom.End))
}

func TestSingleDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "24.12.2015 23:59", dateStr(mom.End))
}

func TestRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015-25.12.2015)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "25.12.2015 23:59", dateStr(mom.End))
}

func TestEndlessRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015-)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "nil", dateStr(mom.End))
}

func TestStartlessRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (-25.12.2015)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "nil", dateStr(mom.Start))
	assert.Equal(t, "25.12.2015 23:59", dateStr(mom.End))
}

func TestShortYearDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.15)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "24.12.2015 23:59", dateStr(mom.End))
}

func TestZeroPaddedDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (04.01.15)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "04.01.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "04.01.2015 23:59", dateStr(mom.End))
}

func TestNonZeroPaddedDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (4.1.15)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "04.01.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "04.01.2015 23:59", dateStr(mom.End))
}

func TestFaultySingleDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (4.1.asfasf)")

	assert.Equal(t, "blabla (4.1.asfasf)", mom.GetName())
	assert.Equal(t, "nil", dateStr(mom.Start))
	assert.Equal(t, "nil", dateStr(mom.End))
}

func TestFaultyRangeDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (4.1.2015-asgdgd)")

	assert.Equal(t, "blabla (4.1.2015-asgdgd)", mom.GetName())
	assert.Equal(t, "nil", dateStr(mom.Start))
	assert.Equal(t, "nil", dateStr(mom.End))
}

func TestSingleDateWithTime(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015 13:15)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "24.12.2015 23:59", dateStr(mom.End))
	assert.Equal(t, "13:15:00", timeStr(mom.TimeOfDay))
	assert.Equal(t, 22, mom.TimeOfDay.Offset)
	assert.Equal(t, 5, mom.TimeOfDay.Length)
}

func TestRangeDateWithTime(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (24.12.2015-25.12.2015 13:15)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, "24.12.2015 00:00", dateStr(mom.Start))
	assert.Equal(t, "25.12.2015 23:59", dateStr(mom.End))
	assert.Equal(t, "13:15:00", timeStr(mom.TimeOfDay))
	assert.Equal(t, 33, mom.TimeOfDay.Offset)
	assert.Equal(t, 5, mom.TimeOfDay.Length)
}

func TestCalculateSingleCoords(t *testing.T) {
	line := &Line{content: "[] blabla (4.1.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 11, mom.Start.Offset)
	assert.Equal(t, 8, mom.Start.Length)
	assert.Equal(t, 11, mom.End.Offset)
	assert.Equal(t, 8, mom.End.Length)

	line = &Line{content: "[] blabla (   4.1.2015  )"}
	mom, _ = parseSingleMoment(line, line.Content())
	assert.Equal(t, 14, mom.Start.Offset)
	assert.Equal(t, 8, mom.Start.Length)
	assert.Equal(t, 14, mom.End.Offset)
	assert.Equal(t, 8, mom.End.Length)
}

func TestCalculateSingleUnicodeCoords(t *testing.T) {
	line := &Line{content: "[] bläbla (4.1.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 11, mom.Start.Offset)
	assert.Equal(t, 8, mom.Start.Length)
	assert.Equal(t, 11, mom.End.Offset)
	assert.Equal(t, 8, mom.End.Length)
}

func TestCalculateRangeCoords(t *testing.T) {
	line := &Line{content: "[] blabla (4.1.2015-5.1.2015)"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 11, mom.Start.Offset)
	assert.Equal(t, 8, mom.Start.Length)
	assert.Equal(t, 20, mom.End.Offset)
	assert.Equal(t, 8, mom.End.Length)

	line = &Line{content: "[] blabla (  4.1.2015  -   5.1.2015  )"}
	mom, _ = parseSingleMoment(line, line.Content())

	assert.Equal(t, 13, mom.Start.Offset)
	assert.Equal(t, 8, mom.Start.Length)
	assert.Equal(t, 27, mom.End.Offset)
	assert.Equal(t, 8, mom.End.Length)
}

func TestCalculateEndlessRangeCoords(t *testing.T) {
	line := &Line{content: "[] blabla (  4.1.2015  -   )"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 13, mom.Start.Offset)
	assert.Equal(t, 8, mom.Start.Length)
}

func TestCalculateStartlessRangeCoords(t *testing.T) {
	line := &Line{content: "[] blabla (  -   5.1.2015  )"}
	mom, _ := parseSingleMoment(line, line.Content())

	assert.Equal(t, 17, mom.End.Offset)
	assert.Equal(t, 8, mom.End.Length)
}

func TestRecurringMoment(t *testing.T) {
	mom, _ := parseRecurMom("[] blabla (every 5.)")
	assert.NotNil(t, mom)
	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, 0, mom.Offset)
	assert.Equal(t, 20, mom.Length)
	assert.Equal(t, moment.RecurMonthly, mom.Recurrence.Recurrence)
	assert.Equal(t, 5, mom.Recurrence.RefDate.Time.Day())
	assert.Equal(t, 11, mom.Recurrence.RefDate.Offset)
	assert.Equal(t, 8, mom.Recurrence.RefDate.Length)
}

func TestUnicodeRecurringMoment(t *testing.T) {
	mom, _ := parseRecurMom("[] bläbla (every 5.)")
	assert.Equal(t, "bläbla", mom.GetName())
	assert.Equal(t, 0, mom.Offset)
	assert.Equal(t, 20, mom.Length)
	assert.Equal(t, 11, mom.Recurrence.RefDate.Offset)
	assert.Equal(t, 8, mom.Recurrence.RefDate.Length)
}

func TestDoneRecurringMoment(t *testing.T) {
	mom, _ := parseRecurMom("[x] blabla (every 5.)")
	assert.NotNil(t, mom)
	assert.Equal(t, "blabla", mom.GetName())
	assert.True(t, mom.IsDone())
	assert.Equal(t, 21, mom.Length)
	assert.Equal(t, moment.RecurMonthly, mom.Recurrence.Recurrence)
	assert.Equal(t, 5, mom.Recurrence.RefDate.Time.Day())
	assert.Equal(t, 12, mom.Recurrence.RefDate.Offset)
	assert.Equal(t, 8, mom.Recurrence.RefDate.Length)
}

func TestPriorityRecurringMoment(t *testing.T) {
	mom, _ := parseRecurMom("[] blabla! (every 5.)")
	assert.NotNil(t, mom)
	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, 1, mom.GetPriority())
	assert.Equal(t, 21, mom.Length)
	assert.Equal(t, moment.RecurMonthly, mom.Recurrence.Recurrence)
	assert.Equal(t, 5, mom.Recurrence.RefDate.Time.Day())
	assert.Equal(t, 12, mom.Recurrence.RefDate.Offset)
	assert.Equal(t, 8, mom.Recurrence.RefDate.Length)
}

func TestRecurringMomentWithTime(t *testing.T) {
	mom, _ := parseRecurMom("[] blabla (every 5. 13:15)")

	assert.Equal(t, "blabla", mom.GetName())
	assert.Equal(t, moment.RecurMonthly, mom.Recurrence.Recurrence)
	assert.Equal(t, 5, mom.Recurrence.RefDate.Time.Day())
	assert.Equal(t, "13:15:00", timeStr(mom.TimeOfDay))
	assert.Equal(t, 20, mom.TimeOfDay.Offset)
	assert.Equal(t, 5, mom.TimeOfDay.Length)
}

func TestEndingWithBracket(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla)")

	assert.Equal(t, "blabla)", mom.GetName())
}

func TestID(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla #my-id-123")

	assert.Equal(t, "blabla", mom.GetName())
	assert.NotNil(t, mom.GetID())
	assert.Equal(t, "my-id-123", mom.GetID().Value)
	assert.Equal(t, 10, mom.GetID().Offset)
	assert.Equal(t, 10, mom.GetID().Length)
}

func TestIDWithDate(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla (1.12.19) #my-id-123")

	assert.Equal(t, "blabla", mom.GetName())
	assert.NotNil(t, mom.GetID())
	assert.Equal(t, "my-id-123", mom.GetID().Value)
	assert.Equal(t, 20, mom.GetID().Offset)
	assert.Equal(t, 10, mom.GetID().Length)
	assert.Equal(t, "01.12.2019 00:00", dateStr(mom.Start))
}

func TestIDTrimming(t *testing.T) {
	mom, _ := parseSingleMom("[] blabla   #my-id-123  ")

	assert.Equal(t, "blabla", mom.GetName())
	assert.NotNil(t, mom.GetID())
	assert.Equal(t, "my-id-123", mom.GetID().Value)
	assert.Equal(t, 12, mom.GetID().Offset)
	assert.Equal(t, 10, mom.GetID().Length)
}

func dateStr(dt *moment.Date) string {
	if dt == nil {
		return "nil"
	}
	return dt.Time.Format("02.01.2006 15:04")
}

func timeStr(dt *moment.Date) string {
	if dt == nil {
		return "nil"
	}
	return dt.Time.Format("15:04:05")
}

func parseMom(content string) (moment.Moment, error) {
	line := &Line{content: content}
	return parseMoment(line, line.Content())
}

func parseSingleMom(content string) (*moment.SingleMoment, error) {
	mom, err := parseMom(content)
	return mom.(*moment.SingleMoment), err
}

func parseRecurMom(content string) (*moment.RecurMoment, error) {
	mom, err := parseMom(content)
	return mom.(*moment.RecurMoment), err
}
