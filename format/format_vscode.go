package format

import (
	"fmt"
	"github.com/sandro-h/sibylgo/generate"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

const catMarker = "cat"
const momMarker = "mom"
const commentMarker = "com"
const dateMarker = "date"
const timeMarker = "time"
const doneSuffix = ".done"
const prioritySuffix = ".priority"
const untilSuffix = ".until%d"

var getNow = func() time.Time {
	return time.Now()
}

// ForVSCode returns a string of lines containing formatting instructions that
// are understood by the sibyl Visual Studio Code extension to format a todo file.
//
// A formatting instruction has the syntax <start offset>,<end offset>,<formatting string>.
// The offsets are the absolute offsets of the string to be formatted, for the entire file
// (so not scoped by line number). The formatting strings are specific to the vscode extension
// and mark things like a done moment or a comment, or a moment that is due soon.
//
// Example output:
//     1,17,mom
//     10,16,date
//     18,34,mom.until10
//     27,33,date
func ForVSCode(todos *moment.Todos) string {

	res := ""
	for _, c := range todos.Categories {
		appendFmt(&res, c.DocCoords, catMarker)
	}
	for _, m := range todos.Moments {
		formatMoment(&res, m, false)
	}
	return res
}

func formatMoment(res *string, m moment.Moment, parentDone bool) {
	momFmt := momMarker
	done := parentDone || m.IsDone()
	if done {
		momFmt += doneSuffix
	} else {
		formatDueSoon(&momFmt, m)
		if m.GetPriority() > 0 {
			momFmt += prioritySuffix
		}
	}

	appendFmt(res, m.GetDocCoords(), momFmt)

	// Additional format lines:
	if done {
		for _, c := range m.GetComments() {
			appendFmt(res, c.DocCoords, commentMarker+doneSuffix)
		}
	} else {
		formatDates(res, m)
	}

	for _, s := range m.GetSubMoments() {
		formatMoment(res, s, done)
	}
}

func formatDates(res *string, m moment.Moment) {
	switch v := m.(type) {
	case *moment.SingleMoment:
		if v.Start != nil {
			appendFmt(res, v.Start.DocCoords, dateMarker)
		}
		if v.End != nil && (v.Start == nil || v.End.DocCoords != v.Start.DocCoords) {
			appendFmt(res, v.End.DocCoords, dateMarker)
		}
	case *moment.RecurMoment:
		appendFmt(res, v.Recurrence.RefDate.DocCoords, dateMarker)
	}

	if m.GetTimeOfDay() != nil {
		appendFmt(res, m.GetTimeOfDay().DocCoords, timeMarker)
	}
}

func formatDueSoon(momFmt *string, m moment.Moment) {
	// Due until 10 (n-1) days in the future
	n := 11
	today := util.SetToStartOfDay(getNow())
	nDaysFromToday := today.AddDate(0, 0, n)
	nRealHours := nDaysFromToday.Sub(today) / time.Hour

	insts := generate.InstancesWithoutSubs(m, today, nDaysFromToday)
	earliest := n
	for _, inst := range insts {
		// We need to compare hours here because of daylight saving time.
		// Instead of 264h (=11 days) it might only be 263h or 265h,
		// which would lead to the wrong number of days calculated.
		if inst.End.Sub(today)/time.Hour >= nRealHours {
			continue
		}

		d := int(inst.End.Sub(today) / util.Days)
		if d < earliest {
			earliest = d
		}
	}
	if earliest < n {
		*momFmt += fmt.Sprintf(untilSuffix, earliest)
	}
}

func appendFmt(res *string, c moment.DocCoords, format string) {
	*res += fmt.Sprintf("%d,%d,%s\n", c.Offset, c.Offset+c.Length, format)
}
