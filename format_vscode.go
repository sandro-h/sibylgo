package main

import (
	"fmt"
)

func FormatVSCode(todos *Todos) string {

	res := ""
	for _, c := range todos.categories {
		appendFmt(&res, c.DocCoords, "cat")
	}
	for _, m := range todos.moments {
		formatMoment(&res, m, false)
	}
	return res
}

func formatMoment(res *string, m Moment, parentDone bool) {
	momFmt := "mom"
	done := parentDone || m.IsDone()
	if done {
		momFmt += ".done"
	} else {
		formatDueSoon(&momFmt, m)
		if m.GetPriority() > 0 {
			momFmt += ".priority"
		}
	}

	appendFmt(res, m.GetDocCoords(), momFmt)

	// Additional format lines:
	if done {
		for _, c := range m.GetComments() {
			appendFmt(res, c.DocCoords, "com.done")
		}
	} else {
		formatDates(res, m)
	}

	for _, s := range m.GetSubMoments() {
		formatMoment(res, s, done)
	}
}

func formatDates(res *string, m Moment) {
	switch v := m.(type) {
	case *SingleMoment:
		if v.start != nil {
			appendFmt(res, v.start.DocCoords, "mom.date")
		}
		if v.end != nil && (v.start == nil || v.end.DocCoords != v.start.DocCoords) {
			appendFmt(res, v.end.DocCoords, "mom.date")
		}
	case *RecurMoment:
		appendFmt(res, v.recurrence.refDate.DocCoords, "mom.date")
	}
}

func formatDueSoon(res *string, m Moment) {
	// TODO
}

func appendFmt(res *string, c DocCoords, format string) {
	*res += fmt.Sprintf("%d,%d,%s\n", c.offset, c.offset+c.length, format)
}
