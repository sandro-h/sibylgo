package format

import (
	"fmt"
	"time"

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

	return forVSCode(todos, false)
}

// ForVSCodeOptimized does the same as ForVSCode, but will merge format lines
// if they are the same type and consecutive in the todo file.
// This way the editor can apply the same format to a larger range instead
// of multiple formats to individual ranges.
func ForVSCodeOptimized(todos *moment.Todos) string {

	return forVSCode(todos, true)
}

func forVSCode(todos *moment.Todos, optimize bool) string {

	res := ""
	fmtState := &formatterState{
		res:       &res,
		optimized: optimize,
	}

	for _, c := range todos.Categories {
		appendFmt(fmtState, c.DocCoords, catMarker)
	}
	for _, m := range todos.Moments {
		formatMoment(fmtState, m, false)
	}

	hasStoredOptimize := fmtState.optimizedMomFmt != ""
	if hasStoredOptimize {
		appendFmtRaw(fmtState, fmtState.optimizedMomCoords, fmtState.optimizedMomFmt)
	}

	return res
}

func formatMoment(fmtState *formatterState, m moment.Moment, parentDone bool) {
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

	if !fmtState.optimized {
		appendFmt(fmtState, m.GetDocCoords(), momFmt)
	} else {
		appendMomentFmtLater(fmtState, m, momFmt)
	}

	// Additional format lines:
	if done {
		for _, c := range m.GetComments() {
			appendFmt(fmtState, c.DocCoords, commentMarker+doneSuffix)
		}
	} else {
		formatDates(fmtState, m)
		formatID(fmtState, m)
	}

	for _, s := range m.GetSubMoments() {
		formatMoment(fmtState, s, done)
	}
}

func appendMomentFmtLater(fmtState *formatterState, m moment.Moment, momFmt string) {
	hasStoredOptimize := fmtState.optimizedMomFmt != ""
	if hasStoredOptimize && momFmt == fmtState.optimizedMomFmt {
		fmtState.optimizedMomCoords.Length = m.GetDocCoords().Offset - fmtState.optimizedMomCoords.Offset + m.GetDocCoords().Length
	} else {
		if hasStoredOptimize {
			appendFmtRaw(fmtState, fmtState.optimizedMomCoords, fmtState.optimizedMomFmt)
		}
		fmtState.optimizedMomFmt = momFmt
		fmtState.optimizedMomCoords = m.GetDocCoords()
	}

	// (Undone) comment lines aren't formatted, so we need to check explicitly.
	// For all other types, appendFmt will be triggered which will flush the stored fmt.
	canChainFurther := m.GetComments() == nil
	if !canChainFurther {
		appendFmtRaw(fmtState, fmtState.optimizedMomCoords, fmtState.optimizedMomFmt)
		fmtState.optimizedMomFmt = ""
	}
}

func formatDates(fmtState *formatterState, m moment.Moment) {
	switch v := m.(type) {
	case *moment.SingleMoment:
		if v.Start != nil {
			appendFmt(fmtState, v.Start.DocCoords, dateMarker)
		}
		if v.End != nil && (v.Start == nil || v.End.DocCoords != v.Start.DocCoords) {
			appendFmt(fmtState, v.End.DocCoords, dateMarker)
		}
	case *moment.RecurMoment:
		appendFmt(fmtState, v.Recurrence.RefDate.DocCoords, dateMarker)
	}

	if m.GetTimeOfDay() != nil {
		appendFmt(fmtState, m.GetTimeOfDay().DocCoords, timeMarker)
	}
}

func formatID(fmtState *formatterState, m moment.Moment) {
	if m.GetID() != nil {
		appendFmt(fmtState, m.GetID().DocCoords, idMarker)
	}
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

func appendFmt(fmtState *formatterState, c moment.DocCoords, format string) {
	if fmtState.optimizedMomFmt != "" {
		appendFmtRaw(fmtState, fmtState.optimizedMomCoords, fmtState.optimizedMomFmt)
		fmtState.optimizedMomFmt = ""
	}
	appendFmtRaw(fmtState, c, format)
}

func appendFmtRaw(fmtState *formatterState, c moment.DocCoords, format string) {
	*fmtState.res += fmt.Sprintf("%d,%d,%s\n", c.Offset, c.Offset+c.Length, format)
}

type formatterState struct {
	res *string

	optimized          bool
	optimizedMomFmt    string
	optimizedMomCoords moment.DocCoords
}
