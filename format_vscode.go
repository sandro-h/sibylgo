package main

import (
	"fmt"
	"time"
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

func formatDueSoon(momFmt *string, m Moment) {
	// Due until 10 (n-1) days in the future
	n := 11
	today := setToStartOfDay(time.Now())
	insts := GenerateInstancesWithoutSubs(m, today, today.AddDate(0, 0, n))
	earliest := n
	for _, inst := range insts {
		d := int(inst.end.Sub(today) / Days)
		if d < earliest {
			earliest = d
		}
	}
	if earliest < n {
		*momFmt += fmt.Sprintf(".until%d", earliest)
	}
}

func appendFmt(res *string, c DocCoords, format string) {
	*res += fmt.Sprintf("%d,%d,%s\n", c.offset, c.offset+c.length, format)
}
