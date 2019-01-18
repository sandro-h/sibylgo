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

func FormatVSCode(todos *moment.Todos) string {

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
	today := util.SetToStartOfDay(time.Now())
	insts := generate.GenerateInstancesWithoutSubs(m, today, today.AddDate(0, 0, n))
	earliest := n
	for _, inst := range insts {
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
