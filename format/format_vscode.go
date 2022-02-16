package format

import (
	"fmt"
	"time"
	"unicode"

	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
)

const catMarker = "cat"
const momMarker = "mom"
const commentMarker = "com"
const dateMarker = "date"
const timeMarker = "time"
const idMarker = "id"
const doneSuffix = ".done"
const prioritySuffix = ".priority"
const untilSuffix = ".until%d"

type format struct {
	start int
	end   int
	style string
}

var noFormat format = format{-1, -1, ""}

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
	formats := forVSCode(todos)
	return toFormatString(formats)
}

// ForVSCodeOptimized does the same as ForVSCode, but will merge format lines
// if they are the same type and consecutive in the todo file.
// This way the editor can apply the same format to a larger range instead
// of multiple formats to individual ranges.
func ForVSCodeOptimized(todos *moment.Todos, rawTodos string) string {
	unoptimized := forVSCode(todos)

	var optimized []format
	cur := noFormat
	for _, next := range unoptimized {
		merge := cur != noFormat && next.style == cur.style && onlyWhitespaceBetween(cur, next, rawTodos)

		if merge {
			cur.end = next.end
		} else {
			if cur != noFormat {
				optimized = append(optimized, cur)
			}
			cur = next
		}
	}
	if cur != noFormat {
		optimized = append(optimized, cur)
	}

	return toFormatString(optimized)
}

func forVSCode(todos *moment.Todos) []format {

	var formats []format

	for _, c := range todos.Categories {
		formats = appendFmt(formats, c.DocCoords, catMarker)
	}
	for _, m := range todos.Moments {
		formats = append(formats, formatMoment(m, false)...)
	}

	return formats
}

func formatMoment(m moment.Moment, parentDone bool) []format {
	var formats []format
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

	formats = appendFmt(formats, m.GetDocCoords(), momFmt)

	// Additional format lines:
	if done {
		for _, c := range m.GetComments() {
			formats = appendFmt(formats, c.DocCoords, commentMarker+doneSuffix)
		}
	} else {
		formats = append(formats, formatDates(m)...)
		if m.GetID() != nil {
			formats = appendFmt(formats, m.GetID().DocCoords, idMarker)
		}
	}

	for _, s := range m.GetSubMoments() {
		formats = append(formats, formatMoment(s, done)...)
	}

	return formats
}

func formatDates(m moment.Moment) []format {
	var formats []format

	switch v := m.(type) {
	case *moment.SingleMoment:
		if v.Start != nil {
			formats = appendFmt(formats, v.Start.DocCoords, dateMarker)
		}
		if v.End != nil && (v.Start == nil || v.End.DocCoords != v.Start.DocCoords) {
			formats = appendFmt(formats, v.End.DocCoords, dateMarker)
		}
	case *moment.RecurMoment:
		formats = appendFmt(formats, v.Recurrence.RefDate.DocCoords, dateMarker)
	}

	if m.GetTimeOfDay() != nil {
		formats = appendFmt(formats, m.GetTimeOfDay().DocCoords, timeMarker)
	}

	return formats
}

func formatDueSoon(momFmt *string, m moment.Moment) {
	// Due until 10 (n-1) days in the future
	n := 11
	today := util.SetToStartOfDay(getNow())
	nDaysFromToday := today.AddDate(0, 0, n)
	nRealHours := nDaysFromToday.Sub(today) / time.Hour

	insts := instances.GenerateWithoutSubs(m, today, nDaysFromToday)
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

func appendFmt(formats []format, c moment.DocCoords, style string) []format {
	return append(formats, format{c.Offset, c.Offset + c.Length, style})
}

func toFormatString(formats []format) string {
	res := ""
	for _, frmt := range formats {
		res += fmt.Sprintf("%d,%d,%s\n", frmt.start, frmt.end, frmt.style)
	}
	return res
}

func onlyWhitespaceBetween(a format, b format, content string) bool {
	for _, c := range content[a.end:b.start] {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}
